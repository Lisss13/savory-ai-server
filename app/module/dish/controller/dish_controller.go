package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/dish/payload"
	"savory-ai-server/app/module/dish/service"
	fileUploadService "savory-ai-server/app/module/file_upload/service"
	"savory-ai-server/utils/response"
	"strconv"
)

type dishController struct {
	dishService       service.DishService
	fileUploadService fileUploadService.FileUploadService
}

// DishController определяет интерфейс контроллера блюд.
type DishController interface {
	GetByRestaurantID(c *fiber.Ctx) error
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

// GetByRestaurantID возвращает все блюда для указанного ресторана.
// Метод: GET /dishes/restaurant/:restaurant_id
func (c *dishController) GetByRestaurantID(ctx *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid restaurant ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	dishes, err := c.dishService.GetByRestaurantID(uint(restaurantID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dishes,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetDishCategory возвращает блюда, сгруппированные по категориям для ресторана.
// Метод: GET /dishes/category/:restaurant_id
func (c *dishController) GetDishCategory(ctx *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid restaurant ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	dishes, err := c.dishService.GetDishByMenuCategory(uint(restaurantID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dishes,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetByID возвращает блюдо по ID.
// Метод: GET /dishes/:id
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

// Create создаёт новое блюдо.
// Метод: POST /dishes
// Требует: JWT авторизация, restaurant_id в теле запроса
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

	dish, err := c.dishService.Create(req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dish,
		Messages: response.Messages{"Dish created successfully"},
		Code:     fiber.StatusCreated,
	})
}

// Update обновляет блюдо.
// Метод: PUT /dishes/:id
// Требует: JWT авторизация, restaurant_id в теле запроса
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

	dish, err := c.dishService.Update(uint(id), req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dish,
		Messages: response.Messages{"Dish updated successfully"},
		Code:     fiber.StatusOK,
	})
}

// Delete удаляет блюдо.
// Метод: DELETE /dishes/:id
func (c *dishController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	if err := c.dishService.Delete(uint(id)); err != nil {
		return err
	}

	return response.RespDelete(ctx, id, response.Response{
		Messages: response.Messages{"Dish deleted successfully"},
		Code:     fiber.StatusOK,
	})
}

// GetDishOfDay возвращает блюдо дня для ресторана.
// Метод: GET /dishes/dish-of-day/:restaurant_id
func (c *dishController) GetDishOfDay(ctx *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid restaurant ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	dish, err := c.dishService.GetDishOfDay(uint(restaurantID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     dish,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// SetDishOfDay устанавливает блюдо дня.
// Метод: POST /dishes/dish-of-day/:id
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
