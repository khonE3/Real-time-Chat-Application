package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SimpleAuth is a basic middleware for nickname-based auth
// In production, replace with JWT authentication
func SimpleAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("X-User-ID")
		username := c.Get("X-Username")

		if userID == "" || username == "" {
			// Allow request to continue but without auth context
			return c.Next()
		}

		// Store in locals for handlers to access
		c.Locals("userID", userID)
		c.Locals("username", username)

		return c.Next()
	}
}
