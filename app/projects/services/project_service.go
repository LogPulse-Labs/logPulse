package project_service

import (
	"fmt"
	project_models "log-pulse/app/projects/models"
	"log-pulse/pkg/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AuthUser(c *fiber.Ctx) jwt.MapClaims {
	return c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
}

func LoadChannelsForProject(userId string, projectId string) []*project_models.ChannelResponse {
	channels, err := repository.
		NewChannelRepository().
		All(
			bson.D{{Key: "user_id", Value: userId}, {Key: "project_id", Value: projectId}},
			options.Find().SetProjection(bson.D{{Key: "_id", Value: 1}, {Key: "name", Value: 1}}),
		)

	if err != nil {
		fmt.Println("error getting channels:", err)
		return []*project_models.ChannelResponse{}
	}

	mapChannels := make([]*project_models.ChannelResponse, 0, len(*channels))

	for _, channel := range *channels {
		mapChannels = append(mapChannels, &project_models.ChannelResponse{
			ID:   channel.ID,
			Name: channel.Name,
		})
	}

	return mapChannels
}
