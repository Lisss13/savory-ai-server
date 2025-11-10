package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/table/payload"
	"savory-ai-server/app/module/table/service"
	"savory-ai-server/utils/jwt"
	"savory-ai-server/utils/response"
	"strconv"
)

type tableController struct {
	tableService service.TableService
}

type TableController interface {
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetByRestaurantID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func NewTableController(service service.TableService) TableController {
	return &tableController{
		tableService: service,
	}
}

func (c *tableController) GetAll(ctx *fiber.Ctx) error {
	tables, err := c.tableService.GetAll()
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     tables,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *tableController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	table, err := c.tableService.GetByID(uint(id))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     table,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *tableController) GetByRestaurantID(ctx *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return err
	}

	tables, err := c.tableService.GetByRestaurantID(uint(restaurantID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     tables,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *tableController) Create(ctx *fiber.Ctx) error {
	req := new(payload.CreateTableReq)
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

	user := ctx.Locals("user").(jwt.JWTData)

	table, err := c.tableService.Create(req, user.CompanyID)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     table,
		Messages: response.Messages{"Table created successfully"},
		Code:     fiber.StatusCreated,
	})
}

func (c *tableController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	req := new(payload.UpdateTableReq)
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

	user := ctx.Locals("user").(jwt.JWTData)

	table, err := c.tableService.Update(uint(id), req, user.CompanyID)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     table,
		Messages: response.Messages{"Table updated successfully"},
		Code:     fiber.StatusOK,
	})
}

func (c *tableController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	if err := c.tableService.Delete(uint(id)); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"Table deleted successfully"},
		Code:     fiber.StatusOK,
	})
}
