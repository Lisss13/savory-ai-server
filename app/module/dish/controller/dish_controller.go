package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/dish/payload"
	"savory-ai-server/app/module/dish/service"
	fileUploadService "savory-ai-server/app/module/file_upload/service"
	"savory-ai-server/utils/jwt"
	"savory-ai-server/utils/response"
	"strconv"
)

type dishController struct {
	dishService       service.DishService
	fileUploadService fileUploadService.FileUploadService
}

type DishController interface {
	GetAll(c *fiber.Ctx) error
	GetDishCategory(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetDishOfDay(c *fiber.Ctx) error
	SetDishOfDay(c *fiber.Ctx) error
}

func NewDishController(service service.DishService, fileUploadSvc fileUploadService.FileUploadService) DishController {
	return &dishController{
		dishService:       service,
		fileUploadService: fileUploadSvc,
	}
}

func (c *dishController) GetAll(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(jwt.JWTData)
	dishes, err := c.dishService.GetByOrganizationID(user.CompanyID)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dishes,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *dishController) GetDishCategory(ctx *fiber.Ctx) error {
	dishes, err := c.dishService.GetDishByMenuCategory()
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dishes,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *dishController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	dish, err := c.dishService.GetByID(uint(id))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dish,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *dishController) Create(ctx *fiber.Ctx) error {

	// Parse the request body
	req := new(payload.CreateDishReq)
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

	dish, err := c.dishService.Create(req, user.CompanyID)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dish,
		Messages: response.Messages{"Dish created successfully"},
		Code:     fiber.StatusCreated,
	})
}

func (c *dishController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	req := new(payload.UpdateDishReq)
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

	dish, err := c.dishService.Update(uint(id), req, user.CompanyID)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dish,
		Messages: response.Messages{"Dish updated successfully"},
		Code:     fiber.StatusOK,
	})
}

func (c *dishController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	if err := c.dishService.Delete(uint(id)); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"Dish deleted successfully"},
		Code:     fiber.StatusOK,
	})
}

func (c *dishController) GetDishOfDay(ctx *fiber.Ctx) error {
	dish, err := c.dishService.GetDishOfDay()
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dish,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *dishController) SetDishOfDay(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	dish, err := c.dishService.SetDishOfDay(uint(id))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dish,
		Messages: response.Messages{"Dish set as dish of the day successfully"},
		Code:     fiber.StatusOK,
	})
}
