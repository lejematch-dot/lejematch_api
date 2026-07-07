package cmd

import (
	"Lejematch/api/v1/routes"
	"Lejematch/config"
	"Lejematch/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var (
	Version = "testing"
	Build   = "none"
)

func Run() error {
	config.Load()

	//Initialize fiber
	app := fiber.New(fiber.Config{
		AppName:      "Backend",
		ServerHeader: "Lejematch/api/" + Version + " (Build " + Build + ")",
		BodyLimit:    60 * 1024 * 1024, // 60MB — skal være over uploads.maxUploadSize (50MB) så vores egen fejlbesked kan nå at fyre først

	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: config.AppConfigInstance.FrontendURL,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Static("/uploads", "./uploads")

	// Initialize Database
	database.InitDB()
	defer database.CloseDB()

	//Migrate models
	database.Migrate()

	// Seed dummy data — comment out when no longer needed
	database.Seed()

	//Setup routes
	routes.SetupRoutes(app)

	var port = config.AppConfigInstance.APIPort
	return app.Listen(":" + port)
}
