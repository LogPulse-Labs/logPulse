package main

import (
	"fmt"
	"log"
	"log-pulse/pkg/config"
	"log-pulse/pkg/middleware"
	"log-pulse/platform/cache/redis"
	"log-pulse/platform/database"
	"log-pulse/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cancel, err := database.Connect()

	if err != nil {
		log.Fatal("Database Connection Error: ", err)
	}

	redis.Connect()

	app := fiber.New(config.FiberConfig())

	middleware.FiberMiddleware(app)

	routes.AppRoutes(app)
	routes.ApiRoutes(app)

	app.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  false,
				"message": fiber.ErrNotFound.Message,
			})
		},
	)

	app.Hooks().OnShutdown(func() error {
		fmt.Println("shutting down...")

		return nil
	})

	defer cancel()

	log.Fatal(app.Listen(":3100"))
}
