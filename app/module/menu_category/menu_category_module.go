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

func (r *MenuCategoryRouter) RegisterMenuCategoryRoutes(auth fiber.Handler) {
	menuCategoryController := r.Controller.MenuCategory
	r.App.Route("/categories", func(router fiber.Router) {
		router.Get("/", auth, menuCategoryController.GetAll)
		router.Get("/:id", auth, menuCategoryController.GetByID)
		router.Post("/", auth, menuCategoryController.Create)
		router.Delete("/:id", auth, menuCategoryController.Delete)
	})
}
