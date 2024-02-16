package create_channel

import (
	channel_service "log-pulse/app/channels/services"
	"log-pulse/models"
	"log-pulse/pkg/constant"
	"log-pulse/pkg/repository"
	"log-pulse/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func CreateChannelHandler(c *fiber.Ctx) error {
	request := new(CreateChannelRequest)
	projectId := c.Params("id")

	if err := utils.ValidateRequest(c, request); err != nil {
		return utils.ErrorResponse(c, fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	user := utils.AuthUser(c)
	channelRepository := repository.NewChannelRepository()

	if channel_service.ChannelExistsForProject(projectId) == true {
		return utils.AlreadyExists(c, "Channel")
	}

	channel, err := channelRepository.CreateOne(&models.Channel{
		ProjectID: projectId,
		Name:      utils.ToSnakeCase(request.Name),
		UserID:    user["id"].(string),
	})

	if err != nil {
		return constant.ErrSomethingWentWrong
	}

	return utils.SuccessResponse(c, "Channel created successfully", channel)
}
