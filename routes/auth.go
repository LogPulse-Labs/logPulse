package routes

import (
	forgot_password "log-pulse/app/auth/scenerios/forgot-password"
	login_scenerio "log-pulse/app/auth/scenerios/login"
	register_scenerio "log-pulse/app/auth/scenerios/register"
	rate_limit "log-pulse/pkg/middleware/rate-limit"

	"github.com/gofiber/fiber/v2"
)

func AuthV1Routes(api fiber.Router) {
	api.Route("/auth", func(api fiber.Router) {
		api.Post("/login", rate_limit.AuthLimiterMiddleware(), login_scenerio.LoginHandler).Name("login")
		api.Post("/register", register_scenerio.RegisterHandler).Name("register")
		api.Post("/forgot-password", forgot_password.ForgotPasswordHandler).Name("forgot-password")
	}, "auth.")
}
