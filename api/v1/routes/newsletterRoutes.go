package routes

import (
	handler "Lejematch/api/v1/handlers/newsletter"

	"github.com/gofiber/fiber/v2"
)

func SetupNewsletterRoutes(app fiber.Router) {
	app.Get("/newsletter/unsubscribe/:token", handler.Unsubscribe)
	app.Get("/newsletter/subscribe/:token", handler.Subscribe)
}
