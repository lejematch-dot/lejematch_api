package listings

import (
	"Lejematch/internal/database/repo"

	"github.com/gofiber/fiber/v2"
)

func ListCities(c *fiber.Ctx) error {
	listingsRepo := repo.NewListingsRepo()
	cities, err := listingsRepo.DistinctCities(c.Query("category"))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{"cities": cities})
}
