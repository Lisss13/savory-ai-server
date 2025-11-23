package organization

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/organization/controller"
	org_repo "savory-ai-server/app/module/organization/repository"
	"savory-ai-server/app/module/organization/service"
)

type OrganizationRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewOrganizationRouter(fiber *fiber.App, controller *controller.Controller) *OrganizationRouter {
	return &OrganizationRouter{
		App:        fiber,
		Controller: controller,
	}
}

var OrganizationModuler = fx.Options(
	// Repositories
	fx.Provide(org_repo.NewOrganizationRepository),
	fx.Provide(org_repo.NewLanguageRepository),

	// Services
	fx.Provide(service.NewOrganizationService),
	fx.Provide(service.NewLanguageService),

	// Controllers and Router
	fx.Provide(controller.NewControllers),
	fx.Provide(NewOrganizationRouter),
)

func (or *OrganizationRouter) RegisterOrganizationRoutes(auth fiber.Handler) {
	organizationController := or.Controller.Organization
	languageController := or.Controller.Language

	// Organization routes
	or.App.Route("/organization", func(router fiber.Router) {
		router.Get("/", auth, organizationController.GetAll)
		router.Get("/:id", auth, organizationController.GetByID)
		router.Patch("/:id", auth, organizationController.Update)
		router.Post("/:id/users", auth, organizationController.AddUser)
		router.Delete("/:id/users", auth, organizationController.RemoveUser)

		// Organization language routes
		router.Get("/:id/languages", auth, languageController.GetLanguagesByOrganizationID)
		router.Post("/:id/languages", auth, languageController.AddLanguageToOrganization)
		router.Delete("/:id/languages", auth, languageController.RemoveLanguageFromOrganization)
	})

	// Language routes
	or.App.Route("/languages", func(router fiber.Router) {
		router.Get("/", auth, languageController.GetAll)
		router.Get("/:id", auth, languageController.GetByID)
		router.Post("/", auth, languageController.Create)
		router.Patch("/:id", auth, languageController.Update)
		router.Delete("/:id", auth, languageController.Delete)
	})
}
