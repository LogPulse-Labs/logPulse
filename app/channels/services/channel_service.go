package channel_service

import (
	"fmt"
	"log-pulse/pkg/repository"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteChannelsForProject(projectId string) bool {
	channelDeleted, err := repository.NewChannelRepository().
		Delete(bson.D{{Key: "project_id", Value: projectId}}, true)

	if err != nil {
		fmt.Print("Error deleting channels:", err)
	}

	return channelDeleted
}

func ChannelExistsForProject(projectId string) bool {
	return repository.NewChannelRepository().Exists(bson.D{{Key: "project_id", Value: projectId}})
}
