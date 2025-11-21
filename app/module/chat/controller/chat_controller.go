package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/chat/payload"
	"savory-ai-server/app/module/chat/service"
	"savory-ai-server/utils/response"
	"strconv"
)

type chatController struct {
	chatService service.ChatService
}

type ChatController interface {
	GetRestaurantChats(ctx *fiber.Ctx) error
	// Чаты для столиков в ресторане
	StartSessionFromTable(ctx *fiber.Ctx) error
	CloseSessionFromTable(ctx *fiber.Ctx) error
	MessageFromTable(ctx *fiber.Ctx) error
	GetMessageFromTableSession(ctx *fiber.Ctx) error
	GetSessionsFromTable(ctx *fiber.Ctx) error
}

func NewChatController(service service.ChatService) ChatController {
	return &chatController{
		chatService: service,
	}
}

// ----------------------- Table Chat Methods ----------------------

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

// StartSessionFromTable создание чата посетителем для столика в ресторане
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

// MessageFromTable сообщения, которые отправляют клиенты
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
