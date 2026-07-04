package listings

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UpdateListing(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	var req services.UpdateListingRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	listingService := services.NewListingService(repo.NewListingsRepo())
	if err := listingService.Update(id, caller.UserID, caller.IsAdmin, &req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.ErrNotFound
		}
		if errors.Is(err, services.ErrNotOwner) {
			return fiber.ErrForbidden
		}
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusNoContent)
}
