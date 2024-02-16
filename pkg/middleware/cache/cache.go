package cache

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/utils"
)

func CacheControlMiddleware() fiber.Handler {
	return cache.New(cache.Config{
		ExpirationGenerator: func(c *fiber.Ctx, cfg *cache.Config) time.Duration {
			newCacheTime, _ := strconv.Atoi(c.GetRespHeader("Cache-Time", "600"))
			return time.Second * time.Duration(newCacheTime)
		},
		Next: func(c *fiber.Ctx) bool {
			return c.Get("No-Cache") == "true"
		},
		CacheControl: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			if c.Params("id") != "" {
				return utils.CopyString(c.Params("id"))
			}
			return utils.CopyString(c.Path())
		},
		Methods: []string{fiber.MethodGet, fiber.MethodHead},
	})
}
