package reports

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"

	"github.com/gofiber/fiber/v2"
)

func ListReports(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)
	if !caller.IsAdmin {
		return fiber.ErrForbidden
	}

	reportService := services.NewReportService(repo.NewReportsRepo())
	list, err := reportService.List()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(list)
}
