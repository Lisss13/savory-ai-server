package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/menu_category/payload"
	"savory-ai-server/app/module/menu_category/service"
	"savory-ai-server/utils/response"
	"strconv"
)

type menuCategoryController struct {
	menuCategoryService service.MenuCategoryService
}

// MenuCategoryController определяет интерфейс контроллера категорий меню.
type MenuCategoryController interface {
	GetByRestaurantID(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	UpdateSortOrder(c *fiber.Ctx) error
}

func NewMenuCategoryController(service service.MenuCategoryService) MenuCategoryController {
	return &menuCategoryController{
		menuCategoryService: service,
	}
}

// GetByRestaurantID возвращает все категории меню для указанного ресторана.
// Метод: GET /categories/restaurant/:restaurant_id
func (c *menuCategoryController) GetByRestaurantID(ctx *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid restaurant ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	categories, err := c.menuCategoryService.GetByRestaurantID(uint(restaurantID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     categories,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetByID возвращает категорию меню по ID.
// Метод: GET /categories/:id
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

// Create создаёт новую категорию меню.
// Метод: POST /categories
// Требует: JWT авторизация, restaurant_id в теле запроса
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

	category, err := c.menuCategoryService.Create(req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     category,
		Messages: response.Messages{"Category created successfully"},
		Code:     fiber.StatusCreated,
	})
}

// Update обновляет категорию меню.
// Метод: PATCH /categories/:id
// Требует: JWT авторизация
func (c *menuCategoryController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid category ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	req := new(payload.UpdateMenuCategoryReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	category, err := c.menuCategoryService.Update(uint(id), req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     category,
		Messages: response.Messages{"Category updated successfully"},
		Code:     fiber.StatusOK,
	})
}

// Delete удаляет категорию меню по ID.
// Метод: DELETE /categories/:id
func (c *menuCategoryController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	if err := c.menuCategoryService.Delete(uint(id)); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data: struct {
			ID uint `json:"id"`
		}{ID: uint(id)},
		Messages: response.Messages{"Category deleted successfully"},
		Code:     fiber.StatusOK,
	})
}

// UpdateSortOrder массово обновляет порядок сортировки категорий.
// Метод: PUT /categories/sort-order
// Требует: JWT авторизация
// Принимает массив категорий с их новыми позициями.
func (c *menuCategoryController) UpdateSortOrder(ctx *fiber.Ctx) error {
	req := new(payload.UpdateCategoriesSortOrderReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	// Валидация запроса
	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	if err := c.menuCategoryService.UpdateCategoriesSortOrder(req); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"Categories sort order updated successfully"},
		Code:     fiber.StatusOK,
	})
}
