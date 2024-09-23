package middleware

import (
	"church-app-backend/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware protects routes from unauthorized access
func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	claims, err := utils.VerifyJWTToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Store claims in the request context for later use
	c.Locals("user", claims)
	return c.Next()
}
