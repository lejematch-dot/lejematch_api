package users

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/security"
	"Lejematch/internal/services"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UpdatePassword(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	if !caller.IsAdmin && caller.UserID != uint(userID) {
		return fiber.ErrForbidden
	}

	var req services.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	userService := services.NewUserService(repo.NewUsersRepo(), repo.NewProfilesRepo())
	if err := userService.UpdatePassword(userID, &req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.ErrNotFound
		}
		if errors.Is(err, services.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, security.ErrPasswordTooShort) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusNoContent)
}
