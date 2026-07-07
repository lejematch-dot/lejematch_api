package routes

import (
	"Lejematch/api/auth"
	"Lejematch/api/middleware"
	listingHandler "Lejematch/api/v1/handlers/listings"
	seekerHandler "Lejematch/api/v1/handlers/seekers"
	handler "Lejematch/api/v1/handlers/users"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app fiber.Router) {
	// This creates a subgroup for user creation that doesn't use the JWT middleware.
	create := app.Group("/users")
	create.Post("/", middleware.RegisterLimiter, handler.CreateUser)
	create.Get("/:id/profile", handler.GetProfile)
	create.Get("/:id/listings", listingHandler.ListByUser)
	create.Get("/:id/seekers", seekerHandler.ListByUser)

	api := app.Group("/users", auth.JWTmiddleware)
	api.Delete("/:id", handler.DeleteUserAndProfile)
	api.Get("/:id", handler.GetUser)
	api.Patch("/:id", handler.UpdateUser)
	api.Patch("/:id/profile", handler.UpdateProfile)
	api.Put("/:id/password", handler.UpdatePassword)
}
