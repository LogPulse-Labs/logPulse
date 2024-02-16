package routes

import (
	"fmt"
	"log-pulse/pkg/utils"
	"net/http"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func AppRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		err := utils.Retry(3, func(attempts int) error {
			array := []int{1, 5, 6, 6, 7, 2, 9}
			sort.Ints(array)
			e := sort.IntsAreSorted(array)
			fmt.Println(e, array)
			return fmt.Errorf("failed", array)
		}, 1000, nil)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
				"success": false,
			})
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message": "Welcome to LogPulse API!",
			"success": true,
		})
	}).Name("home")
	app.Get("/health", monitor.New(monitor.Config{Title: "Health", Refresh: 5 * time.Second, APIOnly: true})).Name("health")
}
