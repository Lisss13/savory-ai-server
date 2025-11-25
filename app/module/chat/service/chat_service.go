// Package service содержит бизнес-логику для модуля чата.
// Управляет сессиями чата, обрабатывает сообщения и интегрируется с AI.
package service

import (
	"context"
	"errors"
	aiReservation "savory-ai-server/app/module/ai_reservation/service"
	"savory-ai-server/app/module/chat/payload"
	"savory-ai-server/app/module/chat/repository"
	"savory-ai-server/app/module/table/service"
	"savory-ai-server/app/storage"
	"time"
)

// chatService реализует интерфейс ChatService.
// Содержит зависимости от репозитория, сервиса столиков и AI-сервиса.
type chatService struct {
	chatRepo            repository.ChatRepository            // Репозиторий для работы с БД
	tableService        service.TableService                 // Сервис для проверки существования столиков
	aiReservationService aiReservation.AIReservationService  // Сервис AI-бронирования (может быть nil)
}

// ChatService определяет интерфейс бизнес-логики для чата.
// Разделён на две группы методов: Table Chat и Restaurant Chat.
type ChatService interface {
	// =====================================================
	// Table Chat - чат для посетителей за столиком
	// =====================================================
	StartTableSession(req *payload.StartTableSessionReq) (*payload.TableChatSessionsResp, error)     // Создать сессию
	CloseSessionFromTable(sessionID uint) error                                                       // Закрыть сессию
	MessageFromTable(req *payload.SendTableMessageReq) (*payload.MessagesRespFormBot, error)         // Отправить сообщение
	GetTableMessagesFromSession(sessionID uint) (*payload.TableChatMessagesResp, error)              // Получить историю
	GetRestaurantChats(restaurantID uint) (*payload.RestaurantChatSessionsResp, error)               // Legacy: получить чаты
	GetSessionsFromTable(tableID uint) (*payload.TableChatSessionsByTableIDResp, error)              // Получить сессии столика

	// =====================================================
	// Restaurant Chat - общий чат с рестораном
	// =====================================================
	StartRestaurantSession(req *payload.StartRestaurantSessionReq) (*payload.RestaurantChatSessionResp, error)   // Создать сессию
	CloseRestaurantSession(sessionID uint) error                                                                  // Закрыть сессию
	MessageFromRestaurant(req *payload.SendRestaurantMessageReq) (*payload.RestaurantMessagesRespFormBot, error) // Отправить сообщение
	GetRestaurantMessagesFromSession(sessionID uint) (*payload.RestaurantChatMessagesResp, error)                // Получить историю
	GetRestaurantSessions(restaurantID uint) (*payload.RestaurantChatSessionsResp, error)                        // Получить сессии
}

// NewChatService создаёт новый экземпляр сервиса чата.
// aiReservationService может быть nil - тогда используются простые ответы бота.
func NewChatService(chatRepo repository.ChatRepository, tableService service.TableService, aiReservationSvc aiReservation.AIReservationService) ChatService {
	return &chatService{
		chatRepo:            chatRepo,
		tableService:        tableService,
		aiReservationService: aiReservationSvc,
	}
}

// =====================================================
// Table Chat Methods - чат для посетителей за столиком
// =====================================================

// GetRestaurantChats возвращает чат-сессии ресторана.
// Legacy метод - используйте GetRestaurantSessions вместо него.
func (s *chatService) GetRestaurantChats(restaurantID uint) (*payload.RestaurantChatSessionsResp, error) {
	// This method is now a wrapper around GetRestaurantSessions for backward compatibility
	return s.GetRestaurantSessions(restaurantID)
}

// StartTableSession создаёт новую сессию чата для столика.
// Проверяет существование столика перед созданием сессии.
func (s *chatService) StartTableSession(req *payload.StartTableSessionReq) (*payload.TableChatSessionsResp, error) {
	// Check if table exists
	_, err := s.tableService.GetByID(req.TableID)
	if err != nil {
		return nil, errors.New("table not found")
	}

	// Create new session
	session := &storage.TableChatSessions{
		TableID:      req.TableID,
		RestaurantID: req.RestaurantID,
		Active:       true,
		LastActive:   time.Now(),
	}

	createdSession, err := s.chatRepo.CreateTableSession(session)
	if err != nil {
		return nil, err
	}

	resp := &payload.TableChatSessionsResp{
		Session: payload.TableChatSessions{
			ID:         createdSession.ID,
			Active:     createdSession.Active,
			LastActive: createdSession.LastActive,
			Table: payload.TableResp{
				ID:   createdSession.Table.ID,
				Name: createdSession.Table.Name,
			},
		},
	}

	return resp, nil
}

// CloseSessionFromTable закрывает сессию чата для столика.
// Устанавливает флаг active = false.
func (s *chatService) CloseSessionFromTable(sessionID uint) error {
	return s.chatRepo.CloseTableSession(sessionID)
}

// MessageFromTable обрабатывает сообщение от посетителя столика.
// Процесс:
//  1. Проверяет что сессия существует и активна
//  2. Сохраняет сообщение пользователя
//  3. Получает историю чата для контекста AI
//  4. Генерирует ответ через Anthropic Claude (или fallback на простые ответы)
//  5. Сохраняет ответ бота
//  6. Обновляет время активности сессии
func (s *chatService) MessageFromTable(req *payload.SendTableMessageReq) (*payload.MessagesRespFormBot, error) {
	// Check if the session exists and is active
	session, err := s.chatRepo.FindTableSessionByID(req.SessionID)
	if err != nil {
		return nil, err
	}
	if !session.Active {
		return nil, errors.New("session is not active")
	}

	// Create a table chat message
	message := &storage.TableChatMessage{
		TableID:       session.TableID,
		RestaurantID:  session.RestaurantID,
		ChatSessionID: session.ID,
		Content:       req.Content,
		AuthorType:    storage.UserAuthor,
		SentAt:        time.Now(),
	}

	// Save message
	createdUserMessage, err := s.chatRepo.CreateTableMessage(message)
	if err != nil {
		return nil, err
	}

	// Get chat history for context
	chatHistory, err := s.chatRepo.GetMessagesBySessionID(req.SessionID)
	if err != nil {
		return nil, err
	}

	// Convert to AI message format
	var aiMessages []aiReservation.ChatMessage
	for _, msg := range chatHistory {
		role := "user"
		if msg.AuthorType == storage.BotAuthor {
			role = "assistant"
		}
		aiMessages = append(aiMessages, aiReservation.ChatMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Generate AI response
	var messageFromBot string
	if s.aiReservationService != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		messageFromBot, err = s.aiReservationService.GenerateResponse(ctx, session.RestaurantID, aiMessages)
		if err != nil {
			// Fallback to simple response on error
			messageFromBot, _ = generateBotResponse(createdUserMessage.Content)
		}
	} else {
		messageFromBot, _ = generateBotResponse(createdUserMessage.Content)
	}

	botMessage := &storage.TableChatMessage{
		TableID:       session.TableID,
		RestaurantID:  session.RestaurantID,
		ChatSessionID: session.ID,
		Content:       messageFromBot,
		AuthorType:    storage.BotAuthor,
		SentAt:        time.Now(),
	}

	createdBotMessage, err := s.chatRepo.CreateTableMessage(botMessage)
	if err != nil {
		return nil, err
	}

	// Update session activity
	if err = s.chatRepo.UpdateTableSessionActivity(req.SessionID); err != nil {
		return nil, err
	}

	// Map to response
	return &payload.MessagesRespFormBot{
		Message: payload.Message{
			ID:      createdBotMessage.ID,
			Content: createdBotMessage.Content,
			SentAt:  createdBotMessage.SentAt,
		},
	}, nil
}

// GetTableMessagesFromSession возвращает историю сообщений из сессии чата столика.
// Включает сообщения пользователя (UserAuthor) и ответы бота (BotAuthor).
func (s *chatService) GetTableMessagesFromSession(sessionID uint) (*payload.TableChatMessagesResp, error) {
	// Get messages from a repository
	messages, err := s.chatRepo.GetMessagesBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Map messages to response
	var messageResps []payload.TableChatMessageResp
	for _, message := range messages {
		messageResps = append(messageResps, payload.TableChatMessageResp{
			ID:         message.ID,
			Content:    message.Content,
			SentAt:     message.SentAt,
			AuthorType: message.AuthorType,
		})
	}

	return &payload.TableChatMessagesResp{
		Messages: messageResps,
	}, nil
}

// GetSessionsFromTable возвращает все чат-сессии для указанного столика.
// Включает историю сообщений для каждой сессии.
// Пропускает сессии без сообщений.
func (s *chatService) GetSessionsFromTable(tableID uint) (*payload.TableChatSessionsByTableIDResp, error) {
	// Get sessions from a repository
	sessions, err := s.chatRepo.GetSessionsByTableID(tableID)
	if err != nil {
		return nil, err
	}

	// Map sessions to response format
	var sessionResponse []payload.TableChatSessionsResp
	for _, session := range sessions {
		if len(session.Messages) == 0 {
			continue
		}

		sessionResponse = append(sessionResponse, payload.TableChatSessionsResp{
			Session: payload.TableChatSessions{
				ID:         session.ID,
				Active:     session.Active,
				LastActive: session.LastActive,
				Table: payload.TableResp{
					ID:   session.Table.ID,
					Name: session.Table.Name,
				},
				Messages: []payload.TableMessageResp{},
			},
		})
		for _, message := range session.Messages {
			sessionResponse[len(sessionResponse)-1].Session.Messages = append(sessionResponse[len(sessionResponse)-1].Session.Messages, payload.TableMessageResp{
				ID:      message.ID,
				Content: message.Content,
				SentAt:  message.SentAt,
			})
		}
	}

	return &payload.TableChatSessionsByTableIDResp{
		Sessions: sessionResponse,
	}, nil
}

// =====================================================
// Restaurant Chat Methods - общий чат с рестораном
// =====================================================

// StartRestaurantSession создаёт новую сессию общего чата для ресторана.
// Используется для бронирования столиков и вопросов через AI-бота.
func (s *chatService) StartRestaurantSession(req *payload.StartRestaurantSessionReq) (*payload.RestaurantChatSessionResp, error) {
	// Create new session
	session := &storage.RestaurantChatSessions{
		RestaurantID: req.RestaurantID,
		Active:       true,
		LastActive:   time.Now(),
	}

	createdSession, err := s.chatRepo.CreateRestaurantSession(session)
	if err != nil {
		return nil, err
	}

	resp := &payload.RestaurantChatSessionResp{
		Session: payload.RestaurantChatSession{
			ID:         createdSession.ID,
			Active:     createdSession.Active,
			LastActive: createdSession.LastActive,
			Restaurant: payload.RestaurantResp{
				ID:   createdSession.Restaurant.ID,
				Name: createdSession.Restaurant.Name,
			},
		},
	}

	return resp, nil
}

// CloseRestaurantSession закрывает сессию чата ресторана.
// Устанавливает флаг active = false.
func (s *chatService) CloseRestaurantSession(sessionID uint) error {
	return s.chatRepo.CloseRestaurantSession(sessionID)
}

// MessageFromRestaurant обрабатывает сообщение пользователя в чате ресторана.
// Основной метод для взаимодействия с AI-ботом.
//
// Процесс:
//  1. Проверяет что сессия существует и активна
//  2. Сохраняет сообщение пользователя
//  3. Получает историю чата для контекста AI
//  4. Генерирует ответ через Anthropic Claude с поддержкой tool calling:
//     - get_available_slots: проверка доступных слотов
//     - create_reservation: создание бронирования
//     - get_my_reservations: получение броней по телефону
//     - cancel_reservation: отмена брони
//     - get_restaurant_info: информация о ресторане
//  5. Сохраняет ответ бота
//  6. Обновляет время активности сессии
func (s *chatService) MessageFromRestaurant(req *payload.SendRestaurantMessageReq) (*payload.RestaurantMessagesRespFormBot, error) {
	// Check if the session exists and is active
	session, err := s.chatRepo.FindRestaurantSessionByID(req.SessionID)
	if err != nil {
		return nil, err
	}
	if !session.Active {
		return nil, errors.New("session is not active")
	}

	// Create a restaurant chat message
	message := &storage.RestaurantChatMessage{
		RestaurantID:  session.RestaurantID,
		ChatSessionID: session.ID,
		Content:       req.Content,
		AuthorType:    storage.UserAuthor,
		SentAt:        time.Now(),
	}

	// Save message
	createdUserMessage, err := s.chatRepo.CreateRestaurantMessage(message)
	if err != nil {
		return nil, err
	}

	// Get chat history for context
	chatHistory, err := s.chatRepo.GetRestaurantMessagesBySessionID(req.SessionID)
	if err != nil {
		return nil, err
	}

	// Convert to AI message format
	var aiMessages []aiReservation.ChatMessage
	for _, msg := range chatHistory {
		role := "user"
		if msg.AuthorType == storage.BotAuthor {
			role = "assistant"
		}
		aiMessages = append(aiMessages, aiReservation.ChatMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Generate AI response
	var messageFromBot string
	if s.aiReservationService != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		messageFromBot, err = s.aiReservationService.GenerateResponse(ctx, session.RestaurantID, aiMessages)
		if err != nil {
			// Fallback to simple response on error
			messageFromBot, _ = generateBotResponse(createdUserMessage.Content)
		}
	} else {
		messageFromBot, _ = generateBotResponse(createdUserMessage.Content)
	}

	botMessage := &storage.RestaurantChatMessage{
		RestaurantID:  session.RestaurantID,
		ChatSessionID: session.ID,
		Content:       messageFromBot,
		AuthorType:    storage.BotAuthor,
		SentAt:        time.Now(),
	}

	createdBotMessage, err := s.chatRepo.CreateRestaurantMessage(botMessage)
	if err != nil {
		return nil, err
	}

	// Update session activity
	if err = s.chatRepo.UpdateRestaurantSessionActivity(req.SessionID); err != nil {
		return nil, err
	}

	// Map to response
	return &payload.RestaurantMessagesRespFormBot{
		Message: payload.Message{
			ID:      createdBotMessage.ID,
			Content: createdBotMessage.Content,
			SentAt:  createdBotMessage.SentAt,
		},
	}, nil
}

// GetRestaurantMessagesFromSession возвращает историю сообщений из сессии чата.
// Включает сообщения пользователя (UserAuthor) и ответы бота (BotAuthor).
func (s *chatService) GetRestaurantMessagesFromSession(sessionID uint) (*payload.RestaurantChatMessagesResp, error) {
	// Get messages from a repository
	messages, err := s.chatRepo.GetRestaurantMessagesBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	// Map messages to response
	var messageResps []payload.RestaurantChatMessageResp
	for _, message := range messages {
		messageResps = append(messageResps, payload.RestaurantChatMessageResp{
			ID:         message.ID,
			Content:    message.Content,
			SentAt:     message.SentAt,
			AuthorType: message.AuthorType,
		})
	}

	return &payload.RestaurantChatMessagesResp{
		Messages: messageResps,
	}, nil
}

// GetRestaurantSessions возвращает все чат-сессии для указанного ресторана.
// Включает историю сообщений для каждой сессии.
// Пропускает сессии без сообщений.
func (s *chatService) GetRestaurantSessions(restaurantID uint) (*payload.RestaurantChatSessionsResp, error) {
	// Get sessions from a repository
	sessions, err := s.chatRepo.FindRestaurantSessionByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	// Map sessions to response format
	var sessionResponse []payload.RestaurantChatSessionResp
	for _, session := range sessions {
		if len(session.Messages) == 0 {
			continue
		}

		sessionResponse = append(sessionResponse, payload.RestaurantChatSessionResp{
			Session: payload.RestaurantChatSession{
				ID:         session.ID,
				Active:     session.Active,
				LastActive: session.LastActive,
				Restaurant: payload.RestaurantResp{
					ID:   session.Restaurant.ID,
					Name: session.Restaurant.Name,
				},
				Messages: []payload.RestaurantMessageResp{},
			},
		})
		for _, message := range session.Messages {
			sessionResponse[len(sessionResponse)-1].Session.Messages = append(sessionResponse[len(sessionResponse)-1].Session.Messages, payload.RestaurantMessageResp{
				ID:      message.ID,
				Content: message.Content,
				SentAt:  message.SentAt,
			})
		}
	}

	return &payload.RestaurantChatSessionsResp{
		Sessions: sessionResponse,
	}, nil
}
