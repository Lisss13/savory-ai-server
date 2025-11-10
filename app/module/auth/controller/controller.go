package controller

import (
	"savory-ai-server/app/module/auth/service"
)

type AuthControllerS struct {
	Auth AuthController
}

func NewController(authService service.AuthService) *AuthControllerS {
	return &AuthControllerS{
		Auth: NewAuthController(authService),
	}
}
