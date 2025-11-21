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
}

func NewChatService(chatRepo repository.ChatRepository, tableService service.TableService) ChatService {
	return &chatService{
		chatRepo:     chatRepo,
		tableService: tableService,
	}
}

// ----------------------- Table Chat Methods ----------------------

func (s *chatService) GetRestaurantChats(restaurantID uint) (*payload.RestaurantChatSessionsResp, error) {
	// Get sessions from a repository
	sessions, err := s.chatRepo.FindTableSessionByRestaurantID(restaurantID)
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

	return &payload.RestaurantChatSessionsResp{
		Sessions: sessionResponse,
	}, nil
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
