package routes

import (
	"Lejematch/api/auth"
	handler "Lejematch/api/v1/handlers/seekers"

	"github.com/gofiber/fiber/v2"
)

func SetupSeekerRoutes(app fiber.Router) {
	public := app.Group("/seekers")
	public.Get("/", handler.ListSeekers)
	public.Get("/cities", handler.ListCities)
	public.Get("/:id", handler.GetSeeker)

	protected := app.Group("/seekers", auth.JWTmiddleware)
	protected.Post("/", handler.CreateSeeker)
	protected.Patch("/:id", handler.UpdateSeeker)
	protected.Delete("/:id", handler.DeleteSeeker)
	protected.Post("/:id/contact", handler.ContactSeeker)
}
