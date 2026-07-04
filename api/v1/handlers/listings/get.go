package listings

import (
	"Lejematch/internal/database/repo"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetListing(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	listingsRepo := repo.NewListingsRepo()
	listing, err := listingsRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}

	return c.JSON(listing)
}
