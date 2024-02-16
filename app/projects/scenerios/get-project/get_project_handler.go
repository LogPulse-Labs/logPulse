package get_project

import (
	"errors"
	project_models "log-pulse/app/projects/models"
	project_service "log-pulse/app/projects/services"
	"log-pulse/pkg/repository"
	"log-pulse/pkg/utils"
	"log-pulse/platform/database"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProjectHandler(c *fiber.Ctx) error {
	user := utils.AuthUser(c)

	projectId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	project, err := repository.NewProjectRepository().FindOne(projectId, nil)

	if err != nil {
		if errors.Is(err, database.NOTFOUND) {
			return utils.ErrorResponse(c, fiber.Error{
				Message: "Project not found",
				Code:    fiber.StatusNotFound,
			})
		} else {
			return err
		}
	}

	return utils.SuccessResponse(c, "Project details", &project_models.ProjectResponse{
		ID:        project.ID,
		Name:      project.Name,
		CreatedAt: project.CreatedAt,
		Channels:  project_service.LoadChannelsForProject(user["id"].(string), project.ID),
	})
}
