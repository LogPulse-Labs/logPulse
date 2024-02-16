package cors

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CorsMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: strings.Join([]string{fiber.HeaderOrigin, fiber.HeaderContentType, fiber.HeaderAccept}, ","),
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowCredentials: false,
		AllowOriginsFunc: func(origin string) bool {
			return os.Getenv("APP_ENV") == "development"
		},
	})
}
