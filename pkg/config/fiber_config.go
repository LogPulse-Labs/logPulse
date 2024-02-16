package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func FiberConfig() fiber.Config {
	return fiber.Config{
		AppName:            "LogPulse v1",
		BodyLimit:          5 * 1024 * 1024,
		CaseSensitive:      true,
		DisableDefaultDate: false,
		EnablePrintRoutes:  true,
		ProxyHeader:        fiber.HeaderXForwardedFor,
		TrustedProxies:     []string{"*__*"},
		ReadTimeout:        5 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Set Content-Type: text/plain; charset=utf-8
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			if os.Getenv("GO_ENV") == "production" {
				log.Printf(err.Error())

				return c.Status(code).JSON(fiber.Map{
					"status":  false,
					"message": "Internal Server Error",
				})
			}
			// Return status code with error message
			return c.Status(code).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		},
	}
}
