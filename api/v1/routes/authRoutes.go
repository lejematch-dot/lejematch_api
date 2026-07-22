package routes

import (
	"Lejematch/api/middleware"
	handler "Lejematch/api/v1/handlers/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthActionRoutes(app fiber.Router) {
	app.Post("/forgot-password", middleware.ForgotPasswordLimiter, handler.ForgotPassword)
	app.Post("/reset-password", handler.ResetPassword)
}
