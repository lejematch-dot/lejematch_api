package users

import (
	"Lejematch/internal/database/repo"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetProfile(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	profileRepo := repo.NewProfilesRepo()
	profile, err := profileRepo.FindByUserID(uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"displayName": profile.DisplayName,
		"bio":         profile.Bio,
		"city":        profile.City,
		"imageURL":    profile.ImageURL,
	})
}
