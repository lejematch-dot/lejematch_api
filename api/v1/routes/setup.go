package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1")
	SetupUserRoutes(v1)
	SetupListingRoutes(v1)

	auth := v1.Group("/auth")
	setupLoginRoutes(auth)

}
