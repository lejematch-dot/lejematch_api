package routes

import (
	"Lejematch/api/middleware"
	handler "Lejematch/api/v1/handlers/login"

	"github.com/gofiber/fiber/v2"
)

func setupLoginRoutes(app fiber.Router) {
	api := app.Group("/login")
	api.Post("/", middleware.LoginLimiter, handler.Login)
}
