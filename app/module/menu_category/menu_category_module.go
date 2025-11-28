package menu_category

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/menu_category/controller"
	menu_category_repo "savory-ai-server/app/module/menu_category/repository"
	"savory-ai-server/app/module/menu_category/service"
)

type MenuCategoryRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewMenuCategoryRouter(fiber *fiber.App, controller *controller.Controller) *MenuCategoryRouter {
	return &MenuCategoryRouter{
		App:        fiber,
		Controller: controller,
	}
}

var MenuCategoryModule = fx.Options(
	fx.Provide(menu_category_repo.NewMenuCategoryRepository),
	fx.Provide(service.NewMenuCategoryService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewMenuCategoryRouter),
)

// RegisterMenuCategoryRoutes регистрирует маршруты для категорий меню.
// GET /categories/restaurant/:restaurant_id - получить категории ресторана
// GET /categories/:id - получить категорию по ID
// POST /categories - создать категорию (требует авторизации)
// DELETE /categories/:id - удалить категорию (требует авторизации)
func (r *MenuCategoryRouter) RegisterMenuCategoryRoutes(auth fiber.Handler) {
	menuCategoryController := r.Controller.MenuCategory
	r.App.Route("/categories", func(router fiber.Router) {
		router.Get("/restaurant/:restaurant_id", menuCategoryController.GetByRestaurantID)
		router.Get("/:id", menuCategoryController.GetByID)
		router.Post("/", auth, menuCategoryController.Create)
		router.Delete("/:id", auth, menuCategoryController.Delete)
	})
}
