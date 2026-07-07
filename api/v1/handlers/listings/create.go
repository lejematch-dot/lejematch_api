package listings

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func CreateListing(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	var req services.CreateListingRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	listingService := services.NewListingService(repo.NewListingsRepo())
	listing, err := listingService.Create(caller.UserID, &req)
	if err != nil {
		if errors.Is(err, services.ErrTooFewImages) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tilføj mindst 5 billeder."})
		}
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        listing.ID,
		"createdAt": listing.CreatedAt,
	})
}
