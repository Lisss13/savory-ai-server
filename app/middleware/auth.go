package middleware

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/utils/config"
	"savory-ai-server/utils/jwt"
	"strings"
)

// AuthRequired user auth middleware
func AuthRequired(cfg *config.Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")

		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		isValid, data := jwt.
			NewJWT(cfg.Middleware.Jwt.Secret, cfg.Middleware.Jwt.Expiration).
			ParseToken(token)

		if !isValid {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		ctx.Locals("user", *data)

		return ctx.Next()
	}
}
