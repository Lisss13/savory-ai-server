package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/auth/payload"
	"savory-ai-server/app/module/auth/service"
	"savory-ai-server/utils/jwt"
	"savory-ai-server/utils/response"
)

type authController struct {
	authService service.AuthService
}

type AuthController interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
	RequestPasswordReset(c *fiber.Ctx) error
	VerifyPasswordReset(c *fiber.Ctx) error
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
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	res, err := ac.authService.Login(*req)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusUnauthorized,
		})
	}

	return response.Resp(c, response.Response{
		Data:     res,
		Messages: response.Messages{"Login success"},
		Code:     fiber.StatusOK,
	})
}

// Register
func (ac *authController) Register(ctx *fiber.Ctx) error {
	req := new(payload.RegisterRequest)
	if err := response.ParseAndValidate(ctx, req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	res, err := ac.authService.Register(*req)
	if err != nil {
		// Check for duplicate email error
		if err.Error() == "email already exists" {
			return response.Resp(ctx, response.Response{
				Messages: response.Messages{err.Error()},
				Code:     fiber.StatusConflict,
			})
		}
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     res,
		Messages: response.Messages{"Register success"},
		Code:     fiber.StatusCreated,
	})
}

// ChangePassword изменение пароля пользователя
func (ac *authController) ChangePassword(ctx *fiber.Ctx) error {
	// Get user ID from JWT token
	user := ctx.Locals("user").(jwt.JWTData)

	// Parse request
	req := new(payload.ChangePasswordRequest)
	if err := response.ParseAndValidate(ctx, req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Call service
	res, err := ac.authService.ChangePassword(user.ID, *req)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     res,
		Messages: response.Messages{"Password changed successfully"},
		Code:     fiber.StatusOK,
	})
}

// RequestPasswordReset handles a request to reset a password
func (ac *authController) RequestPasswordReset(ctx *fiber.Ctx) error {
	// Parse request
	req := new(payload.RequestPasswordResetRequest)
	if err := response.ParseAndValidate(ctx, req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Call service
	res, err := ac.authService.RequestPasswordReset(*req)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusInternalServerError,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     res,
		Messages: response.Messages{"Password reset code sent if email exists"},
		Code:     fiber.StatusOK,
	})
}

// VerifyPasswordReset verifies a password reset code and sets a new password
func (ac *authController) VerifyPasswordReset(ctx *fiber.Ctx) error {
	// Parse request
	req := new(payload.VerifyPasswordResetRequest)
	if err := response.ParseAndValidate(ctx, req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Call service
	res, err := ac.authService.VerifyPasswordReset(*req)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     res,
		Messages: response.Messages{"Password reset successful"},
		Code:     fiber.StatusOK,
	})
}
