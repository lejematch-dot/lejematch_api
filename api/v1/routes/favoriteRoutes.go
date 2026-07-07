package routes

import (
	"Lejematch/api/auth"
	handler "Lejematch/api/v1/handlers/favorites"

	"github.com/gofiber/fiber/v2"
)

func SetupFavoriteRoutes(app fiber.Router) {
	protected := app.Group("/favorites", auth.JWTmiddleware)
	protected.Get("/", handler.ListFavorites)
	protected.Post("/", handler.CreateFavorite)
	protected.Delete("/", handler.DeleteFavorite)
}
