package routes

import (
	"Lejematch/api/auth"
	"Lejematch/api/middleware"
	handler "Lejematch/api/v1/handlers/uploads"

	"github.com/gofiber/fiber/v2"
)

func SetupUploadRoutes(app fiber.Router) {
	protected := app.Group("/uploads", auth.JWTmiddleware)
	protected.Post("/", handler.CreateUpload)

	// Uautentificeret, men rate-begrænset — bruges kun til profilbilledet på
	// registreringssiden, hvor der endnu ikke findes en konto/JWT.
	// NB: stien må IKKE starte med "uploads" — Fibers Group-middleware
	// matcher præfikset som ren streng (ikke sti-segment-bevidst), så selv
	// "/uploads-registration" ramte stadig JWTmiddleware ovenfor.
	app.Post("/registration-upload", middleware.PublicUploadLimiter, handler.CreateUpload)
}
