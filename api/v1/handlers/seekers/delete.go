package seekers

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DeleteSeeker(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	seekerService := services.NewSeekerService(repo.NewSeekersRepo())
	if err := seekerService.Delete(id, caller.UserID, caller.IsAdmin); err != nil {
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
