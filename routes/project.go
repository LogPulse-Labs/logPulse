package routes

import (
	create_channel "log-pulse/app/channels/scenerios/create-channel"
	create_project "log-pulse/app/projects/scenerios/create-project"
	delete_project "log-pulse/app/projects/scenerios/delete-project"
	get_project "log-pulse/app/projects/scenerios/get-project"
	get_projects "log-pulse/app/projects/scenerios/get-projects"
	"log-pulse/pkg/middleware/auth"

	"github.com/gofiber/fiber/v2"
)

func ProjectV1Routes(api fiber.Router) {
	project := api.Group("/projects", auth.AuthenticateMiddleware())

	project.Post("/", create_project.CreateProjectHandler).Name("projects.create")
	project.Get("/", get_projects.GetProjectsHandler).Name("projects.all")
	project.Get("/:id", get_project.GetProjectHandler).Name("projects.single")

	project.Delete("/:id", delete_project.DeleteProjectHandler).Name("projects.delete")
	project.Post("/:id/channels", create_channel.CreateChannelHandler).Name("project.channels.create")
}
