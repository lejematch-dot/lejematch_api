package users

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	var req services.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userService := services.NewUserService(repo.NewUsersRepo(), repo.NewProfilesRepo())
	user, err := userService.CreateUserWithProfile(&req)
	if err != nil {
		if errors.Is(err, services.ErrDuplicateEntry) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, services.ErrImageRequired) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tilføj et profilbillede."})
		}
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        user.ID,
		"createdAt": user.CreatedAt,
	})
}
