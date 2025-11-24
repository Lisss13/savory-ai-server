package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/admin/payload"
	"savory-ai-server/app/module/admin/service"
	"savory-ai-server/utils/jwt"
	"savory-ai-server/utils/response"
	"strconv"
)

type adminController struct {
	adminService service.AdminService
}

type AdminController interface {
	// Статистика
	GetStats(c *fiber.Ctx) error

	// Пользователи
	GetAllUsers(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	UpdateUserStatus(c *fiber.Ctx) error
	UpdateUserRole(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error

	// Организации
	GetAllOrganizations(c *fiber.Ctx) error
	GetOrganizationByID(c *fiber.Ctx) error
	DeleteOrganization(c *fiber.Ctx) error

	// Модерация контента
	GetAllDishes(c *fiber.Ctx) error
	DeleteDish(c *fiber.Ctx) error

	// Логи
	GetAllLogs(c *fiber.Ctx) error
	GetMyLogs(c *fiber.Ctx) error
}

func NewAdminController(service service.AdminService) AdminController {
	return &adminController{adminService: service}
}

// ==================== Статистика ====================

// GetStats - GET /admin/stats
// Возвращает общую статистику системы
func (c *adminController) GetStats(ctx *fiber.Ctx) error {
	stats, err := c.adminService.GetStats()
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     stats,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// ==================== Пользователи ====================

// GetAllUsers - GET /admin/users
// Возвращает список всех пользователей с пагинацией
func (c *adminController) GetAllUsers(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize", "20"))

	users, err := c.adminService.GetAllUsers(page, pageSize)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     users,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetUserByID - GET /admin/users/:id
// Возвращает информацию о конкретном пользователе
func (c *adminController) GetUserByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	user, err := c.adminService.GetUserByID(uint(id))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     user,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// UpdateUserStatus - PATCH /admin/users/:id/status
// Блокирует или разблокирует пользователя
func (c *adminController) UpdateUserStatus(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	req := new(payload.UpdateUserStatusReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	admin := ctx.Locals("user").(jwt.JWTData)
	ipAddress := ctx.IP()

	if err := c.adminService.UpdateUserStatus(admin.ID, uint(id), req.IsActive, ipAddress); err != nil {
		return err
	}

	action := "blocked"
	if req.IsActive {
		action = "unblocked"
	}

	return response.Resp(ctx, response.Response{
		Data:     struct{ ID uint `json:"id"` }{ID: uint(id)},
		Messages: response.Messages{"User " + action + " successfully"},
		Code:     fiber.StatusOK,
	})
}

// UpdateUserRole - PATCH /admin/users/:id/role
// Изменяет роль пользователя
func (c *adminController) UpdateUserRole(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	req := new(payload.UpdateUserRoleReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	admin := ctx.Locals("user").(jwt.JWTData)
	ipAddress := ctx.IP()

	if err := c.adminService.UpdateUserRole(admin.ID, uint(id), req.Role, ipAddress); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     struct{ ID uint `json:"id"` }{ID: uint(id)},
		Messages: response.Messages{"User role updated successfully"},
		Code:     fiber.StatusOK,
	})
}

// DeleteUser - DELETE /admin/users/:id
// Удаляет пользователя
func (c *adminController) DeleteUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	admin := ctx.Locals("user").(jwt.JWTData)
	ipAddress := ctx.IP()

	if err := c.adminService.DeleteUser(admin.ID, uint(id), ipAddress); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     struct{ ID uint `json:"id"` }{ID: uint(id)},
		Messages: response.Messages{"User deleted successfully"},
		Code:     fiber.StatusOK,
	})
}

// ==================== Организации ====================

// GetAllOrganizations - GET /admin/organizations
// Возвращает список всех организаций
func (c *adminController) GetAllOrganizations(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize", "20"))

	orgs, err := c.adminService.GetAllOrganizations(page, pageSize)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     orgs,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetOrganizationByID - GET /admin/organizations/:id
// Возвращает информацию о конкретной организации
func (c *adminController) GetOrganizationByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	org, err := c.adminService.GetOrganizationByID(uint(id))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     org,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// DeleteOrganization - DELETE /admin/organizations/:id
// Удаляет организацию
func (c *adminController) DeleteOrganization(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	admin := ctx.Locals("user").(jwt.JWTData)
	ipAddress := ctx.IP()

	if err := c.adminService.DeleteOrganization(admin.ID, uint(id), ipAddress); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     struct{ ID uint `json:"id"` }{ID: uint(id)},
		Messages: response.Messages{"Organization deleted successfully"},
		Code:     fiber.StatusOK,
	})
}

// ==================== Модерация контента ====================

// GetAllDishes - GET /admin/dishes
// Возвращает список всех блюд для модерации
func (c *adminController) GetAllDishes(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize", "20"))

	dishes, err := c.adminService.GetAllDishes(page, pageSize)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dishes,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// DeleteDish - DELETE /admin/dishes/:id
// Удаляет блюдо (модерация)
func (c *adminController) DeleteDish(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	admin := ctx.Locals("user").(jwt.JWTData)
	ipAddress := ctx.IP()

	if err := c.adminService.DeleteDish(admin.ID, uint(id), ipAddress); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     struct{ ID uint `json:"id"` }{ID: uint(id)},
		Messages: response.Messages{"Dish deleted successfully"},
		Code:     fiber.StatusOK,
	})
}

// ==================== Логи ====================

// GetAllLogs - GET /admin/logs
// Возвращает логи всех действий администраторов
func (c *adminController) GetAllLogs(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize", "50"))

	logs, err := c.adminService.GetAllLogs(page, pageSize)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     logs,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetMyLogs - GET /admin/logs/me
// Возвращает логи текущего администратора
func (c *adminController) GetMyLogs(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize", "50"))

	admin := ctx.Locals("user").(jwt.JWTData)

	logs, err := c.adminService.GetLogsByAdminID(admin.ID, page, pageSize)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     logs,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}
