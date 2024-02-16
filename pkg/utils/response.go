package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(c *fiber.Ctx, payload fiber.Error) error {
	return c.Status(payload.Code).JSON(fiber.Map{
		"success": false,
		"message": payload.Message,
	})
}

func AlreadyExists(c *fiber.Ctx, name string) error {
	return ErrorResponse(c, fiber.Error{
		Message: fmt.Sprintf("%v already exists", name),
		Code:    fiber.StatusBadRequest,
	})
}
