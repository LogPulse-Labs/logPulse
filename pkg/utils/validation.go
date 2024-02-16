package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateRequest(c *fiber.Ctx, request interface{}) error {
	if err := c.BodyParser(request); err != nil {
		return err
	}

	if err := validator.New().Struct(request); err != nil {
		return errors.New(CustomValidationMessages(err))
	}

	return nil
}
