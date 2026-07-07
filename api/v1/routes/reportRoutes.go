package routes

import (
	"Lejematch/api/auth"
	"Lejematch/api/middleware"
	handler "Lejematch/api/v1/handlers/reports"

	"github.com/gofiber/fiber/v2"
)

func SetupReportRoutes(app fiber.Router) {
	protected := app.Group("/reports", auth.JWTmiddleware)
	protected.Post("/", middleware.ReportLimiter, handler.CreateReport)

	admin := app.Group("/admin/reports", auth.JWTmiddleware)
	admin.Get("/", handler.ListReports)
	admin.Patch("/:id/resolve", handler.ResolveReport)
}
