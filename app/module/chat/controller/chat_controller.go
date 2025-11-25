package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/chat/payload"
	"savory-ai-server/app/module/chat/service"
	"savory-ai-server/utils/response"
	"strconv"
)

// chatController реализует интерфейс ChatController.
type chatController struct {
	chatService service.ChatService
}

// ChatController определяет интерфейс HTTP-обработчиков для чата.
// Разделён на две группы: Table Chat и Restaurant Chat.
type ChatController interface {
	// Legacy endpoint
	GetRestaurantChats(ctx *fiber.Ctx) error // Получить чаты ресторана (устаревший)

	// =====================================================
	// Table Chat - чат для посетителей за столиком
	// =====================================================
	StartSessionFromTable(ctx *fiber.Ctx) error       // Начать сессию чата для столика
	CloseSessionFromTable(ctx *fiber.Ctx) error       // Закрыть сессию чата
	MessageFromTable(ctx *fiber.Ctx) error            // Отправить сообщение → получить ответ AI
	GetMessageFromTableSession(ctx *fiber.Ctx) error  // Получить историю сообщений сессии
	GetSessionsFromTable(ctx *fiber.Ctx) error        // Получить все сессии столика

	// =====================================================
	// Restaurant Chat - общий чат с рестораном
	// =====================================================
	StartRestaurantSession(ctx *fiber.Ctx) error          // Начать сессию чата
	CloseRestaurantSession(ctx *fiber.Ctx) error          // Закрыть сессию чата
	MessageFromRestaurant(ctx *fiber.Ctx) error           // Отправить сообщение → получить ответ AI
	GetRestaurantMessagesFromSession(ctx *fiber.Ctx) error // Получить историю сообщений
	GetRestaurantSessions(ctx *fiber.Ctx) error           // Получить все сессии ресторана
}

// NewChatController создаёт новый экземпляр контроллера чата.
func NewChatController(service service.ChatService) ChatController {
	return &chatController{
		chatService: service,
	}
}

// =====================================================
// Table Chat Methods - чат для посетителей за столиком
// =====================================================

// GetRestaurantChats возвращает все чат-сессии ресторана.
// Legacy endpoint, используйте GetRestaurantSessions вместо него.
//
// Метод: GET /chat/restaurant/:restaurant_id
func (c *chatController) GetRestaurantChats(ctx *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid restaurant ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	sessions, err := c.chatService.GetRestaurantChats(uint(restaurantID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     sessions,
		Messages: response.Messages{"Chat sessions retrieved successfully"},
		Code:     fiber.StatusOK,
	})
}

// StartSessionFromTable создаёт новую сессию чата для столика.
// Вызывается посетителем при сканировании QR-кода столика.
//
// Метод: POST /chat/table/session/start
// Тело запроса: { tableId, restaurantId }
// Ответ: Данные созданной сессии
func (c *chatController) StartSessionFromTable(ctx *fiber.Ctx) error {
	req := new(payload.StartTableSessionReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	// Validate request
	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	session, err := c.chatService.StartTableSession(req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     session,
		Messages: response.Messages{"Chat session started successfully"},
		Code:     fiber.StatusOK,
	})
}

// CloseSessionFromTable закрывает сессию чата для столика.
// После закрытия отправка сообщений в эту сессию невозможна.
//
// Метод: POST /chat/table/session/close/:session_id
func (c *chatController) CloseSessionFromTable(ctx *fiber.Ctx) error {
	sessionID, err := strconv.ParseUint(ctx.Params("session_id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid table ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	if err = c.chatService.CloseSessionFromTable(uint(sessionID)); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"Chat session closed successfully"},
		Code:     fiber.StatusOK,
	})
}

// MessageFromTable обрабатывает сообщение от посетителя и генерирует ответ AI.
// Сохраняет сообщение пользователя, генерирует ответ через Anthropic Claude,
// сохраняет ответ бота и возвращает его клиенту.
//
// Метод: POST /chat/table/message/send
// Тело запроса: { sessionId, content }
// Ответ: Сообщение от AI-бота
func (c *chatController) MessageFromTable(ctx *fiber.Ctx) error {
	req := new(payload.SendTableMessageReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	// Validate request
	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	message, err := c.chatService.MessageFromTable(req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     message,
		Messages: response.Messages{"Message sent successfully"},
		Code:     fiber.StatusOK,
	})
}

// GetMessageFromTableSession возвращает историю сообщений из сессии чата.
// Включает сообщения пользователя и ответы AI-бота.
//
// Метод: GET /chat/table/session/:session_id/messages
func (c *chatController) GetMessageFromTableSession(ctx *fiber.Ctx) error {
	sessionID, err := strconv.ParseUint(ctx.Params("session_id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid session ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	messages, err := c.chatService.GetTableMessagesFromSession(uint(sessionID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     messages,
		Messages: response.Messages{"Messages retrieved successfully"},
		Code:     fiber.StatusOK,
	})
}

// GetSessionsFromTable возвращает все чат-сессии для указанного столика.
// Используется персоналом для просмотра истории общения посетителей.
//
// Метод: GET /chat/table/session/:table_id
func (c *chatController) GetSessionsFromTable(ctx *fiber.Ctx) error {
	tableID, err := strconv.ParseUint(ctx.Params("table_id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid session ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	sessions, err := c.chatService.GetSessionsFromTable(uint(tableID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     sessions,
		Messages: response.Messages{"Sessions retrieved successfully"},
		Code:     fiber.StatusOK,
	})
}

// =====================================================
// Restaurant Chat Methods - общий чат с рестораном
// =====================================================

// StartRestaurantSession создаёт новую сессию чата для ресторана.
// Используется для общего чата с AI-ботом ресторана (бронирование, вопросы).
//
// Метод: POST /chat/restaurant/session/start
// Тело запроса: { restaurantId }
// Ответ: Данные созданной сессии
func (c *chatController) StartRestaurantSession(ctx *fiber.Ctx) error {
	req := new(payload.StartRestaurantSessionReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	// Validate request
	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	session, err := c.chatService.StartRestaurantSession(req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     session,
		Messages: response.Messages{"Restaurant chat session started successfully"},
		Code:     fiber.StatusOK,
	})
}

// CloseRestaurantSession закрывает сессию чата ресторана.
// После закрытия отправка сообщений в эту сессию невозможна.
//
// Метод: POST /chat/restaurant/session/close/:session_id
func (c *chatController) CloseRestaurantSession(ctx *fiber.Ctx) error {
	sessionID, err := strconv.ParseUint(ctx.Params("session_id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid session ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	if err = c.chatService.CloseRestaurantSession(uint(sessionID)); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"Restaurant chat session closed successfully"},
		Code:     fiber.StatusOK,
	})
}

// MessageFromRestaurant обрабатывает сообщение пользователя и генерирует ответ AI.
// Основной метод для взаимодействия с AI-ботом ресторана.
// Поддерживает tool calling для бронирования столиков через Anthropic Claude.
//
// Метод: POST /chat/restaurant/message/send
// Тело запроса: { sessionId, content }
// Ответ: Сообщение от AI-бота (может включать результаты бронирования)
func (c *chatController) MessageFromRestaurant(ctx *fiber.Ctx) error {
	req := new(payload.SendRestaurantMessageReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	// Validate request
	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	message, err := c.chatService.MessageFromRestaurant(req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     message,
		Messages: response.Messages{"Message sent successfully"},
		Code:     fiber.StatusOK,
	})
}

// GetRestaurantMessagesFromSession возвращает историю сообщений из сессии чата.
// Включает сообщения пользователя и ответы AI-бота с информацией об авторе.
//
// Метод: GET /chat/restaurant/session/:session_id/messages
func (c *chatController) GetRestaurantMessagesFromSession(ctx *fiber.Ctx) error {
	sessionID, err := strconv.ParseUint(ctx.Params("session_id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid session ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	messages, err := c.chatService.GetRestaurantMessagesFromSession(uint(sessionID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     messages,
		Messages: response.Messages{"Messages retrieved successfully"},
		Code:     fiber.StatusOK,
	})
}

// GetRestaurantSessions возвращает все чат-сессии для указанного ресторана.
// Используется персоналом для мониторинга общения с посетителями.
//
// Метод: GET /chat/restaurant/sessions/:restaurant_id
func (c *chatController) GetRestaurantSessions(ctx *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid restaurant ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	sessions, err := c.chatService.GetRestaurantSessions(uint(restaurantID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     sessions,
		Messages: response.Messages{"Restaurant chat sessions retrieved successfully"},
		Code:     fiber.StatusOK,
	})
}
