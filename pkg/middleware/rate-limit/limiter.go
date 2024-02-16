package rate_limit

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
)

var (
	TooManyAttemptsMessage = fiber.Map{
		"status":  false,
		"message": strings.Join([]string{fiber.ErrTooManyRequests.Message, ".Please try again after 1 minutes"}, ""),
	}
)

var redisPort, _ = strconv.Atoi(os.Getenv("REDIS_PORT"))
var redisDB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))

var redisStorage = redis.New(redis.Config{
	Host:     os.Getenv("REDIS_HOST"),
	Port:     redisPort,
	Database: redisDB,
	// redis cluster setup
	// Addrs:    strings.Split(strings.TrimSpace(os.Getenv("REDIS_CLUSTER")), ","),
})

func DefaultLimiterMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        500,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("X-Forwarded-For", c.IP())
		},
		Storage: redisStorage,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.ErrTooManyRequests.Code).JSON(TooManyAttemptsMessage)
		},
	})
}

func AuthLimiterMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Params("email")
		},
		Storage: redisStorage,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.ErrTooManyRequests.Code).JSON(TooManyAttemptsMessage)
		},
		SkipSuccessfulRequests: true,
	})
}
