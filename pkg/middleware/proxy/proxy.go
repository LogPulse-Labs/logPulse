package proxy

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func ProxyMiddleware() fiber.Handler {
	return proxy.Balancer(proxy.Config{
		Servers: strings.Split(strings.TrimSpace(os.Getenv("PROXY_SERVERS")), ","),
		ModifyRequest: func(c *fiber.Ctx) error {
			c.Request().Header.Add("X-Real-IP", c.IP())
			return nil
		},
		Timeout: 5 * time.Second,
		ModifyResponse: func(c *fiber.Ctx) error {
			c.Response().Header.Del(fiber.HeaderServer)
			return nil
		},
	})
}
