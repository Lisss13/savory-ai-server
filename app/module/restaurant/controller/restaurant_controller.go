package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/restaurant/payload"
	"savory-ai-server/app/module/restaurant/service"
	"savory-ai-server/utils/response"
	"strconv"
)

type restaurantController struct {
	restaurantService service.RestaurantService
}

type RestaurantController interface {
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetByOrganizationID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func NewRestaurantController(service service.RestaurantService) RestaurantController {
	return &restaurantController{
		restaurantService: service,
	}
}

func (c *restaurantController) GetAll(ctx *fiber.Ctx) error {
	restaurants, err := c.restaurantService.GetAll()
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     restaurants,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *restaurantController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	restaurant, err := c.restaurantService.GetByID(uint(id))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     restaurant,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *restaurantController) GetByOrganizationID(ctx *fiber.Ctx) error {
	organizationID, err := strconv.ParseUint(ctx.Params("organization_id"), 10, 32)
	if err != nil {
		return err
	}

	restaurants, err := c.restaurantService.GetByOrganizationID(uint(organizationID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     restaurants,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *restaurantController) Create(ctx *fiber.Ctx) error {
	req := new(payload.CreateRestaurantReq)
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

	restaurant, err := c.restaurantService.Create(req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     restaurant,
		Messages: response.Messages{"Restaurant created successfully"},
		Code:     fiber.StatusCreated,
	})
}

func (c *restaurantController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	req := new(payload.UpdateRestaurantReq)
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

	restaurant, err := c.restaurantService.Update(uint(id), req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     restaurant,
		Messages: response.Messages{"Restaurant updated successfully"},
		Code:     fiber.StatusOK,
	})
}

func (c *restaurantController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	if err := c.restaurantService.Delete(uint(id)); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"Restaurant deleted successfully"},
		Code:     fiber.StatusOK,
	})
}
