package create_project

import (
	"log-pulse/models"
	"log-pulse/pkg/constant"
	"log-pulse/pkg/repository"
	"log-pulse/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateProjectHandler(c *fiber.Ctx) error {
	payload := new(CreateProjectRequest)

	if err := utils.ValidateRequest(c, payload); err != nil {
		return utils.ErrorResponse(c, fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	projectRepository := repository.NewProjectRepository()

	projectName := strings.TrimSpace(payload.Name)
	user := utils.AuthUser(c)

	exist := projectRepository.
		Exists(bson.D{{Key: "user_id", Value: user["id"]}, {Key: "name", Value: projectName}})

	if exist == true {
		return utils.AlreadyExists(c, "Project")
	}

	project, err := projectRepository.CreateOne(&models.Project{
		Name:   projectName,
		UserID: user["id"].(string),
	})

	if err != nil {
		return constant.ErrSomethingWentWrong
	}

	return utils.SuccessResponse(c, "Project created successfully", project)
}
