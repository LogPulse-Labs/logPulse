package get_organization

import (
	"errors"
	"log-pulse/pkg/repository"
	"log-pulse/pkg/utils"
	"log-pulse/platform/database"

	"github.com/gofiber/fiber/v2"
)

func GetOrganizationHandler(c *fiber.Ctx) error {
	user := utils.AuthUser(c)

	organization, err := repository.NewOrganizationRepository().FindByUser(user["id"].(string))

	if err != nil {
		if errors.Is(err, database.NOTFOUND) {
			return utils.ErrorResponse(c, fiber.Error{
				Message: "Organization not found",
				Code:    fiber.StatusNotFound,
			})
		} else {
			return err
		}
	}

	return utils.SuccessResponse(c, "Project Details", organization)
}
