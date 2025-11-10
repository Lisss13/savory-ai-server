package restaurant

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/restaurant/controller"
	restaurant_repo "savory-ai-server/app/module/restaurant/repository"
	"savory-ai-server/app/module/restaurant/service"
)

type RestaurantRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewRestaurantRouter(fiber *fiber.App, controller *controller.Controller) *RestaurantRouter {
	return &RestaurantRouter{
		App:        fiber,
		Controller: controller,
	}
}

var RestaurantModule = fx.Options(
	fx.Provide(restaurant_repo.NewRestaurantRepository),
	fx.Provide(service.NewRestaurantService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewRestaurantRouter),
)

func (r *RestaurantRouter) RegisterRestaurantRoutes(auth fiber.Handler) {
	restaurantController := r.Controller.Restaurant
	r.App.Route("/restaurants", func(router fiber.Router) {
		router.Get("/", auth, restaurantController.GetAll)
		router.Get("/:id", auth, restaurantController.GetByID)
		router.Get("/organization/:organization_id", auth, restaurantController.GetByOrganizationID)
		router.Post("/", auth, restaurantController.Create)
		router.Put("/:id", auth, restaurantController.Update)
		router.Delete("/:id", auth, restaurantController.Delete)
	})
}
