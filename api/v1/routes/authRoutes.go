package routes

import (
	"Lejematch/api/middleware"
	handler "Lejematch/api/v1/handlers/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthActionRoutes(app fiber.Router) {
	app.Get("/verify-email/:token", handler.VerifyEmail)
	app.Post("/resend-verification", middleware.ResendVerificationLimiter, handler.ResendVerification)
	app.Post("/forgot-password", middleware.ForgotPasswordLimiter, handler.ForgotPassword)
	app.Post("/reset-password", handler.ResetPassword)
}
