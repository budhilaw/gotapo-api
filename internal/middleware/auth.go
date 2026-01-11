package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// TapoCredentials extracts Tapo camera credentials from request headers
func TapoCredentials() fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Get("X-Tapo-Username")
		password := c.Get("X-Tapo-Password")

		if username == "" || password == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "Missing X-Tapo-Username or X-Tapo-Password headers",
			})
		}

		c.Locals("tapo_username", username)
		c.Locals("tapo_password", password)

		return c.Next()
	}
}

// GetTapoCredentials retrieves Tapo credentials from context
func GetTapoCredentials(c *fiber.Ctx) (username, password string) {
	username, _ = c.Locals("tapo_username").(string)
	password, _ = c.Locals("tapo_password").(string)
	return
}
