package middleware

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/utils/jwt"
)

// AdminRequired проверяет, что пользователь является администратором
func AdminRequired() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(jwt.JWTData)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		if user.Role != "admin" {
			return fiber.NewError(fiber.StatusForbidden, "Access denied. Admin privileges required")
		}

		return ctx.Next()
	}
}
