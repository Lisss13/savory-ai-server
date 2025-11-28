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

// RegisterDishRoutes регистрирует маршруты для блюд.
// GET /dishes/restaurant/:restaurant_id - получить блюда ресторана
// GET /dishes/category/:restaurant_id - получить блюда по категориям
// GET /dishes/dish-of-day/:restaurant_id - получить блюдо дня
// POST /dishes/dish-of-day/:id - установить блюдо дня (требует авторизации)
// GET /dishes/:id - получить блюдо по ID
// POST /dishes - создать блюдо (требует авторизации)
// PUT /dishes/:id - обновить блюдо (требует авторизации)
// DELETE /dishes/:id - удалить блюдо (требует авторизации)
func (r *DishRouter) RegisterDishRoutes(auth fiber.Handler) {
	dishController := r.Controller.Dish
	r.App.Route("/dishes", func(router fiber.Router) {
		router.Get("/restaurant/:restaurant_id", dishController.GetByRestaurantID)
		router.Get("/category/:restaurant_id", dishController.GetDishCategory)
		router.Get("/dish-of-day/:restaurant_id", dishController.GetDishOfDay)
		router.Post("/dish-of-day/:id", auth, dishController.SetDishOfDay)
		router.Get("/:id", dishController.GetByID)
		router.Post("/", auth, dishController.Create)
		router.Put("/:id", auth, dishController.Update)
		router.Delete("/:id", auth, dishController.Delete)
	})
}
