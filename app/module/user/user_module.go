package user

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/user/controller"
	user_repo "savory-ai-server/app/module/user/repository"
	"savory-ai-server/app/module/user/service"
)

type UserRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewUserRouter(fiber *fiber.App, controller *controller.Controller) *UserRouter {
	return &UserRouter{
		App:        fiber,
		Controller: controller,
	}
}

var UserModuler = fx.Options(
	fx.Provide(user_repo.NewUserRepository),
	fx.Provide(service.NewUserService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewUserRouter),
)

func (ur *UserRouter) RegisterUserRouters(auth fiber.Handler) {
	userController := ur.Controller.User
	ur.App.Route("/user", func(router fiber.Router) {
		router.Get("/:id", auth, userController.Get)
		router.Patch("/:id", auth, userController.Update)
		router.Post("/", auth, userController.Create)
	})
}
