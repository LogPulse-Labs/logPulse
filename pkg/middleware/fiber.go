package middleware

import (
	"log-pulse/pkg/middleware/cache"
	"log-pulse/pkg/middleware/cors"
	rate_limit "log-pulse/pkg/middleware/rate-limit"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func FiberMiddleware(app *fiber.App) {
	app.Use(
		cors.CorsMiddleware(),
		AcceptJsonMiddleware(),
		logger.New(logger.Config{
			DisableColors: false,
			TimeFormat:    time.RFC3339Nano,
			TimeZone:      "Africa/Lagos",
			Format:        "[${time}] - [${ip}]:${port} - ${status} - ${latency} ${method} ${path}\n",
			Done: func(c *fiber.Ctx, logString []byte) {
				if c.Response().StatusCode() != fiber.StatusOK {
					return
				}
			},
		}),
		cache.CacheControlMiddleware(),
		favicon.New(),
		helmet.New(),
		rate_limit.DefaultLimiterMiddleware(),
		requestid.New(),
		recover.New(),
	)
}
