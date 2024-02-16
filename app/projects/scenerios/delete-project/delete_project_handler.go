package delete_project

import (
	"errors"
	"fmt"
	channel_service "log-pulse/app/channels/services"
	"log-pulse/pkg/constant"
	"log-pulse/pkg/repository"
	"log-pulse/pkg/utils"
	"log-pulse/platform/database"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DeleteProjectHandler(c *fiber.Ctx) error {
	projectId, _ := primitive.ObjectIDFromHex(strings.TrimSpace(c.Params("id")))
	projectRepo := repository.NewProjectRepository()

	_, err := projectRepo.FindOne(
		projectId,
		options.FindOne().SetProjection(bson.D{{Key: "id", Value: 1}}),
	)

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

	deleted, delErr := projectRepo.Delete(projectId)
	if delErr != nil {
		fmt.Println("err deleting project:", err)
		return constant.ErrSomethingWentWrong
	}

	channel_service.DeleteChannelsForProject(projectId.String())

	return utils.SuccessResponse(c, "Project deleted successfully", deleted)
}
