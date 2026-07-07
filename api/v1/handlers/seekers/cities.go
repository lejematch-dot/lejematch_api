package seekers

import (
	"Lejematch/internal/database/repo"

	"github.com/gofiber/fiber/v2"
)

func ListCities(c *fiber.Ctx) error {
	seekersRepo := repo.NewSeekersRepo()
	cities, err := seekersRepo.DistinctCities(c.Query("category"))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{"cities": cities})
}
