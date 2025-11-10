package table

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/table/controller"
	table_repo "savory-ai-server/app/module/table/repository"
	"savory-ai-server/app/module/table/service"
)

type TableRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewTableRouter(fiber *fiber.App, controller *controller.Controller) *TableRouter {
	return &TableRouter{
		App:        fiber,
		Controller: controller,
	}
}

var TableModule = fx.Options(
	fx.Provide(table_repo.NewTableRepository),
	fx.Provide(service.NewTableService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewTableRouter),
)

func (r *TableRouter) RegisterTableRoutes(auth fiber.Handler) {
	tableController := r.Controller.Table
	r.App.Route("/tables", func(router fiber.Router) {
		router.Get("/", auth, tableController.GetAll)
		router.Get("/:id", auth, tableController.GetByID)
		router.Get("/restaurant/:restaurant_id", auth, tableController.GetByRestaurantID)
		router.Post("/", auth, tableController.Create)
		router.Put("/:id", auth, tableController.Update)
		router.Delete("/:id", auth, tableController.Delete)
	})
}
