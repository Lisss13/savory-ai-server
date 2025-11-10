package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/menu_category/payload"
	"savory-ai-server/app/module/menu_category/service"
	"savory-ai-server/utils/jwt"
	"savory-ai-server/utils/response"
	"strconv"
)

type menuCategoryController struct {
	menuCategoryService service.MenuCategoryService
}

type MenuCategoryController interface {
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func NewMenuCategoryController(service service.MenuCategoryService) MenuCategoryController {
	return &menuCategoryController{
		menuCategoryService: service,
	}
}

func (c *menuCategoryController) GetAll(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(jwt.JWTData)
	categories, err := c.menuCategoryService.GetByOrganizationID(user.CompanyID)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     categories,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *menuCategoryController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	category, err := c.menuCategoryService.GetByID(uint(id))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     category,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *menuCategoryController) Create(ctx *fiber.Ctx) error {
	req := new(payload.CreateMenuCategoryReq)
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

	category, err := c.menuCategoryService.Create(req, user.CompanyID)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     category,
		Messages: response.Messages{"Category created successfully"},
		Code:     fiber.StatusCreated,
	})
}

func (c *menuCategoryController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	if err := c.menuCategoryService.Delete(uint(id)); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"Category deleted successfully"},
		Code:     fiber.StatusOK,
	})
}
