package subscription

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/subscription/controller"
	subscription_repo "savory-ai-server/app/module/subscription/repository"
	"savory-ai-server/app/module/subscription/service"
)

type SubscriptionRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewSubscriptionRouter(fiber *fiber.App, controller *controller.Controller) *SubscriptionRouter {
	return &SubscriptionRouter{
		App:        fiber,
		Controller: controller,
	}
}

var SubscriptionModule = fx.Options(
	fx.Provide(subscription_repo.NewSubscriptionRepository),
	fx.Provide(service.NewSubscriptionService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewSubscriptionRouter),
)

func (r *SubscriptionRouter) RegisterSubscriptionRoutes(auth fiber.Handler) {
	subscriptionController := r.Controller.Subscription

	r.App.Route("/subscriptions", func(router fiber.Router) {
		// Get all subscriptions
		router.Get("/", auth, subscriptionController.GetAll)

		// Get subscription by ID
		router.Get("/:id", auth, subscriptionController.GetByID)

		// Get subscriptions by organization ID
		router.Get("/organization/:organizationId", auth, subscriptionController.GetByOrganizationID)

		// Get active subscription by organization ID
		router.Get("/organization/:organizationId/active", auth, subscriptionController.GetActiveByOrganizationID)

		// Create a new subscription
		router.Post("/", auth, subscriptionController.Create)

		// Update subscription
		router.Put("/:id", auth, subscriptionController.Update)

		// Extend subscription (add more months)
		router.Post("/:id/extend", auth, subscriptionController.Extend)

		// Deactivate subscription
		router.Post("/:id/deactivate", auth, subscriptionController.Deactivate)

		// Delete subscription
		router.Delete("/:id", auth, subscriptionController.Delete)
	})
}
