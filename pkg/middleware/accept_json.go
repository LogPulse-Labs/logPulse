package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func AcceptJsonMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Finished-At", time.Now().String())
		c.Set("Accept", fiber.MIMEApplicationJSON)
		c.Set("Content-Type", fiber.MIMEApplicationJSON)

		return c.Next()
	}
}
