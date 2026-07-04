package users

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DeleteUserAndProfile(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	if !caller.IsAdmin && caller.UserID != uint(userID) {
		return fiber.ErrForbidden
	}

	userRepo := repo.NewUsersRepo()

	user, err := userRepo.FindByID(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return fiber.ErrInternalServerError
	}

	err = userRepo.Delete(user)
	if err != nil {
		return err
	}

	return c.JSON(user.ID)
}
