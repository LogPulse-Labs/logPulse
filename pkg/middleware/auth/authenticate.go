package auth

import (
	"fmt"
	"log-pulse/pkg/utils"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthenticateMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(fmt.Sprintf("%v", os.Getenv("JWT_SECRET")))},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	return utils.ErrorResponse(c, fiber.Error{
		Message: "Unauthorized.",
		Code:    fiber.StatusUnauthorized,
	})
}
