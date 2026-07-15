package cmd

import (
	"Lejematch/api/v1/routes"
	"Lejematch/config"
	"Lejematch/internal/database"
	"Lejematch/internal/services"

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
		BodyLimit:    30 * 1024 * 1024, // 30MB — skal være over uploads.maxUploadSize (25MB) så vores egen fejlbesked kan nå at fyre først

		// Caddy er den eneste indgang til denne container (internt Docker-netværk,
		// ikke offentligt tilgængelig direkte). Uden dette ser c.IP() Caddys egen
		// IP for alle requests, hvilket får IP-baserede rate-limiters (login,
		// registrering, m.fl.) til fejlagtigt at dele én bucket mellem alle brugere.
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
		ProxyHeader:             fiber.HeaderXForwardedFor,
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

	services.StartDailyContactDigest()

	//Setup routes
	routes.SetupRoutes(app)

	var port = config.AppConfigInstance.APIPort
	return app.Listen(":" + port)
}
