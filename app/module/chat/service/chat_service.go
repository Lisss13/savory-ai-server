package service

import (
	"errors"
	"savory-ai-server/app/module/chat/payload"
	"savory-ai-server/app/module/chat/repository"
	"savory-ai-server/app/module/table/service"
	"savory-ai-server/app/storage"
	"time"
)

type chatService struct {
	chatRepo     repository.ChatRepository
	tableService service.TableService
}

type ChatService interface {
	// Message operations for tables
	StartTableSession(req *payload.StartTableSessionReq) (*payload.TableChatSessionsResp, error)
	CloseSessionFromTable(sessionID uint) error
	MessageFromTable(req *payload.SendTableMessageReq) (*payload.MessagesRespFormBot, error)
	GetTableMessagesFromSession(sessionID uint) (*payload.TableChatMessagesResp, error)
	GetRestaurantChats(restaurantID uint) (*payload.RestaurantChatSessionsResp, error)
	GetSessionsFromTable(tableID uint) (*payload.TableChatSessionsByTableIDResp, error)

	// Message operations for restaurants
	StartRestaurantSession(req *payload.StartRestaurantSessionReq) (*payload.RestaurantChatSessionResp, error)
	CloseRestaurantSession(sessionID uint) error
	MessageFromRestaurant(req *payload.SendRestaurantMessageReq) (*payload.RestaurantMessagesRespFormBot, error)
	GetRestaurantMessagesFromSession(sessionID uint) (*payload.RestaurantChatMessagesResp, error)
	GetRestaurantSessions(restaurantID uint) (*payload.RestaurantChatSessionsResp, error)
}

func NewChatService(chatRepo repository.ChatRepository, tableService service.TableService) ChatService {
	return &chatService{
		chatRepo:     chatRepo,
		tableService: tableService,
	}
}

// ----------------------- Table Chat Methods ----------------------

// GetRestaurantChats gets table chat sessions for a restaurant (legacy method)
func (s *chatService) GetRestaurantChats(restaurantID uint) (*payload.RestaurantChatSessionsResp, error) {
	// This method is now a wrapper around GetRestaurantSessions for backward compatibility
	return s.GetRestaurantSessions(restaurantID)
}

// StartTableSession создание чата посетителем для столика в ресторане
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

// CloseSessionFromTable закрытие сессии чата для столика
func (s *chatService) CloseSessionFromTable(sessionID uint) error {
	return s.chatRepo.CloseTableSession(sessionID)
}

// MessageFromTable сообщения, которые отправляют клиенты и получают от бота
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

	messageFromBot, _ := generateBotResponse(createdUserMessage.Content)

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

// GetTableMessagesFromSession получение сообщений из сессии чата для столика
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

// GetSessionsFromTable получение сессий чата для столика
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

// ----------------------- Restaurant Chat Methods ----------------------

// StartRestaurantSession создание чата для ресторана
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

// CloseRestaurantSession закрытие сессии чата для ресторана
func (s *chatService) CloseRestaurantSession(sessionID uint) error {
	return s.chatRepo.CloseRestaurantSession(sessionID)
}

// MessageFromRestaurant сообщения, которые отправляют клиенты и получают от бота
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

	messageFromBot, _ := generateBotResponse(createdUserMessage.Content)

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

// GetRestaurantMessagesFromSession получение сообщений из сессии чата для ресторана
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

// GetRestaurantSessions получение сессий чата для ресторана
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
