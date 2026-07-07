package routes

import (
	"Lejematch/api/auth"
	handler "Lejematch/api/v1/handlers/uploads"

	"github.com/gofiber/fiber/v2"
)

func SetupUploadRoutes(app fiber.Router) {
	protected := app.Group("/uploads", auth.JWTmiddleware)
	protected.Post("/", handler.CreateUpload)
}
