package seekers

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"

	"github.com/gofiber/fiber/v2"
)

func CreateSeeker(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	var req services.CreateSeekerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	seekerService := services.NewSeekerService(repo.NewSeekersRepo())
	seeker, err := seekerService.Create(caller.UserID, &req)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        seeker.ID,
		"createdAt": seeker.CreatedAt,
	})
}
