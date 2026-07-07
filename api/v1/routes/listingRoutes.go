package routes
import (
	"Lejematch/api/auth"
	handler "Lejematch/api/v1/handlers/listings"
	"github.com/gofiber/fiber/v2"
)
func SetupListingRoutes(app fiber.Router) {
	public := app.Group("/listings")
	public.Get("/", handler.ListListings)
	public.Get("/cities", handler.ListCities)
	public.Get("/:id", handler.GetListing)
	protected := app.Group("/listings", auth.JWTmiddleware)
	protected.Post("/", handler.CreateListing)
	protected.Patch("/:id", handler.UpdateListing)
	protected.Delete("/:id", handler.DeleteListing)
	protected.Post("/:id/contact", handler.ContactListing)
}
