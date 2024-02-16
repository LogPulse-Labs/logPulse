package get_projects_scenerio

import (
	project_models "log-pulse/app/projects/models"
	project_service "log-pulse/app/projects/services"
	"log-pulse/pkg/constant"
	"log-pulse/pkg/repository"
	"log-pulse/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetProjectsHandler(c *fiber.Ctx) error {
	user := project_service.AuthUser(c)

	projects, err := repository.NewProjectRepository().All(
		bson.D{{Key: "user_id", Value: user["id"]}},
		nil,
	)

	if err != nil {
		return utils.ErrorResponse(c, fiber.Error{
			Message: constant.ErrSomethingWentWrong.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	projectsWithChannels := make([]*project_models.ProjectResponse, 0, len(*projects))

	for _, project := range *projects {
		projectsWithChannels = append(projectsWithChannels, &project_models.ProjectResponse{
			ID:        project.ID,
			Name:      project.Name,
			CreatedAt: project.CreatedAt,
			Channels:  project_service.LoadChannelsForProject(user["id"].(string), project.ID),
		})
	}

	return utils.SuccessResponse(c, "Projects Fetched", projectsWithChannels)
}
