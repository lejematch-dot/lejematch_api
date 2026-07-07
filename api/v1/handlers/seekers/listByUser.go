package seekers

import (
	"Lejematch/internal/database/repo"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ListByUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	seekersRepo := repo.NewSeekersRepo()
	seekers, err := seekersRepo.FindByUserID(uint(id))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(seekers)
}
