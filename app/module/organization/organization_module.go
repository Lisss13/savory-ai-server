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
	fx.Provide(org_repo.NewOrganizationRepository),
	fx.Provide(service.NewOrganizationService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewOrganizationRouter),
)

func (or *OrganizationRouter) RegisterOrganizationRoutes(auth fiber.Handler) {
	organizationController := or.Controller.Organization
	or.App.Route("/organization", func(router fiber.Router) {
		router.Get("/", auth, organizationController.GetAll)
		router.Get("/:id", auth, organizationController.GetByID)
		router.Patch("/:id", auth, organizationController.Update)
		router.Post("/:id/users", auth, organizationController.AddUser)
		router.Delete("/:id/users", auth, organizationController.RemoveUser)
	})
}
