package routes

import "github.com/gofiber/fiber/v2"

func ApiRoutes(app *fiber.App) {
	api := app.Group("/api/v1", func(c *fiber.Ctx) error {
		c.Accepts(fiber.MIMETextPlain, fiber.MIMEApplicationJSON)
		c.Accepts(fiber.MIMEApplicationJSON, fiber.MIMETextHTML)

		return c.Next()
	})

	AuthV1Routes(api)
	ProjectV1Routes(api)
}
