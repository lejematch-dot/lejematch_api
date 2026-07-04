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
		BodyLimit:    10 * 1024 * 1024, // 10MB

	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // TODO: restrict to specific domain(s) before production (e.g. "https://lejematch.dk")
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	// Limit login attempts to 10 per minute per IP
	//app.Use("/api/v1/auth/login", limiter.New(limiter.Config{
	//	Max:        10,
	//	Expiration: 1 * time.Minute,
	//}))

	// Limit account creation to 5 per hour per IP
	//app.Use("/api/v1/users", limiter.New(limiter.Config{
	//	Max:        5,
	//	Expiration: 1 * time.Hour,
	//}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

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
