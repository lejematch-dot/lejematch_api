package reports

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ResolveReport(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)
	if !caller.IsAdmin {
		return fiber.ErrForbidden
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	reportService := services.NewReportService(repo.NewReportsRepo())
	if err := reportService.Resolve(id); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusNoContent)
}
