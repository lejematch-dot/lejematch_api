package admin

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"

	"github.com/gofiber/fiber/v2"
)

// GetStats returnerer samlede statistikker til admin-dashboardet.
func GetStats(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)
	if !caller.IsAdmin {
		return fiber.ErrForbidden
	}

	stats, err := repo.NewStatsRepo().Get()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(stats)
}
