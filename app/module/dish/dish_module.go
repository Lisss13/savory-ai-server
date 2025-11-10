package dish

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/dish/controller"
	dish_repo "savory-ai-server/app/module/dish/repository"
	"savory-ai-server/app/module/dish/service"
)

type DishRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewDishRouter(fiber *fiber.App, controller *controller.Controller) *DishRouter {
	return &DishRouter{
		App:        fiber,
		Controller: controller,
	}
}

var DishModule = fx.Options(
	fx.Provide(dish_repo.NewDishRepository),
	fx.Provide(service.NewDishService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewDishRouter),
)

func (r *DishRouter) RegisterDishRoutes(auth fiber.Handler) {
	dishController := r.Controller.Dish
	r.App.Route("/dishes", func(router fiber.Router) {
		router.Get("/", auth, dishController.GetAll)
		router.Get("/category", auth, dishController.GetDishCategory)
		router.Get("/dish-of-day", auth, dishController.GetDishOfDay)
		router.Get("/:id", auth, dishController.GetByID)
		router.Post("/", auth, dishController.Create)
		router.Put("/:id", auth, dishController.Update)
		router.Delete("/:id", auth, dishController.Delete)
		router.Post("/:id/dish-of-day", auth, dishController.SetDishOfDay)
	})
}
