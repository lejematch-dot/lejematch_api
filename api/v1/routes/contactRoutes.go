package routes

import (
	"Lejematch/api/auth"
	handler "Lejematch/api/v1/handlers/contacts"

	"github.com/gofiber/fiber/v2"
)

func SetupContactRoutes(app fiber.Router) {
	protected := app.Group("/contacts", auth.JWTmiddleware)
	protected.Get("/", handler.ListContacts)
}
