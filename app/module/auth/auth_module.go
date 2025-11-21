package auth

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/auth/controller"
	"savory-ai-server/app/module/auth/service"
)

// struct of AuthRouter
type AuthRouter struct {
	App        fiber.Router
	Controller *controller.AuthControllerS
}

// register bulky of auth module
var NewAuthModule = fx.Options(
	// register service of auth module
	fx.Provide(service.NewAuthService),

	// register controller of auth module
	fx.Provide(controller.NewController),

	// register router of auth module
	fx.Provide(NewAuthRouter),
)

// NewAuthRouter init AuthRouter
func NewAuthRouter(fiber *fiber.App, controller *controller.AuthControllerS) *AuthRouter {
	return &AuthRouter{
		App:        fiber,
		Controller: controller,
	}
}

// RegisterAuthRoutes register routes of auth module
func (ar *AuthRouter) RegisterAuthRoutes(auth fiber.Handler) {
	authController := ar.Controller.Auth
	ar.App.Route("/auth", func(router fiber.Router) {
		router.Post("/login", authController.Login)
		router.Post("/register", authController.Register)
		router.Post("/change-password", auth, authController.ChangePassword)

		// Password reset routes
		router.Post("/request-password-reset", authController.RequestPasswordReset)
		router.Post("/verify-password-reset", authController.VerifyPasswordReset)
	})
}
