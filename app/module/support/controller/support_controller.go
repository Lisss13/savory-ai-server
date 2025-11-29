// Package controller содержит HTTP обработчики для модуля поддержки.
package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/support/payload"
	"savory-ai-server/app/module/support/service"
	"savory-ai-server/utils/jwt"
	"savory-ai-server/utils/response"
)

type supportController struct {
	service service.SupportService
}

// SupportController определяет интерфейс HTTP обработчиков для заявок в поддержку
type SupportController interface {
	// User endpoints
	Create(ctx *fiber.Ctx) error
	GetMyTickets(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error

	// Admin endpoints
	GetAll(ctx *fiber.Ctx) error
	UpdateStatus(ctx *fiber.Ctx) error
}

// NewSupportController создаёт новый экземпляр контроллера поддержки
func NewSupportController(service service.SupportService) SupportController {
	return &supportController{service: service}
}

// Create создаёт новую заявку в поддержку
// Метод: POST /support
// Требует: JWT авторизация
func (c *supportController) Create(ctx *fiber.Ctx) error {
	// Получаем текущего пользователя
	currentUser := ctx.Locals("user").(jwt.JWTData)

	// Парсим запрос
	req := new(payload.CreateSupportTicketReq)
	if err := ctx.BodyParser(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"неверный формат запроса"},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Валидация
	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Создаём заявку
	ticket, err := c.service.Create(currentUser.ID, req)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusInternalServerError,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     ticket,
		Messages: response.Messages{"заявка успешно создана"},
		Code:     fiber.StatusCreated,
	})
}

// GetMyTickets возвращает заявки текущего пользователя
// Метод: GET /support/my
// Query: page, page_size
// Требует: JWT авторизация
func (c *supportController) GetMyTickets(ctx *fiber.Ctx) error {
	// Получаем текущего пользователя
	currentUser := ctx.Locals("user").(jwt.JWTData)

	// Получаем параметры пагинации
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Получаем заявки
	tickets, err := c.service.GetMyTickets(currentUser.ID, page, pageSize)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusInternalServerError,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     tickets,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetByID возвращает заявку по ID
// Метод: GET /support/:id
// Требует: JWT авторизация
func (c *supportController) GetByID(ctx *fiber.Ctx) error {
	// Получаем ID из параметров
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"неверный ID заявки"},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Получаем заявку
	ticket, err := c.service.GetByID(uint(id))
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusNotFound,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     ticket,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetAll возвращает все заявки (для админов)
// Метод: GET /admin/support
// Query: page, page_size, status
// Требует: JWT авторизация + role: admin
func (c *supportController) GetAll(ctx *fiber.Ctx) error {
	// Получаем параметры пагинации
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))
	status := ctx.Query("status", "")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var tickets *payload.SupportTicketsListResp
	var err error

	// Если указан статус, фильтруем по нему
	if status != "" {
		tickets, err = c.service.GetByStatus(status, page, pageSize)
	} else {
		tickets, err = c.service.GetAll(page, pageSize)
	}

	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     tickets,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// UpdateStatus обновляет статус заявки (для админов)
// Метод: PATCH /admin/support/:id/status
// Требует: JWT авторизация + role: admin
func (c *supportController) UpdateStatus(ctx *fiber.Ctx) error {
	// Получаем ID из параметров
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"неверный ID заявки"},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Парсим запрос
	req := new(payload.UpdateSupportTicketStatusReq)
	if err := ctx.BodyParser(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"неверный формат запроса"},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Валидация
	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Обновляем статус
	ticket, err := c.service.UpdateStatus(uint(id), req)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     ticket,
		Messages: response.Messages{"статус успешно обновлён"},
		Code:     fiber.StatusOK,
	})
}
