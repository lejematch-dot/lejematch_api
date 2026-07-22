package routes

import (
	"Lejematch/api/auth"
	handler "Lejematch/api/v1/handlers/admin"

	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(app fiber.Router) {
	admin := app.Group("/admin", auth.JWTmiddleware)
	admin.Post("/contact-digest/trigger", handler.TriggerContactDigest)
	admin.Get("/stats", handler.GetStats)
	admin.Post("/newsletter/send", handler.SendNewsletterHandler)
	admin.Post("/newsletter/invite-existing", handler.SendNewsletterInviteHandler)
}
