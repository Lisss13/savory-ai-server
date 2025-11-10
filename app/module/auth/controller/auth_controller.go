package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/auth/payload"
	"savory-ai-server/app/module/auth/service"
	"savory-ai-server/utils/response"
)

type authController struct {
	authService service.AuthService
}

type AuthController interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
}

// NewAuthController
func NewAuthController(authService service.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

// Login
func (ac *authController) Login(c *fiber.Ctx) error {
	req := new(payload.LoginRequest)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	res, err := ac.authService.Login(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data:     res,
		Messages: response.Messages{"Login success"},
		Code:     fiber.StatusOK,
	})
}

// Register
func (ac *authController) Register(c *fiber.Ctx) error {
	req := new(payload.RegisterRequest)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	res, err := ac.authService.Register(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data:     res,
		Messages: response.Messages{"Register success"},
		Code:     fiber.StatusOK,
	})
}
