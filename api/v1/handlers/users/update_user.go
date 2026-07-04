package users

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func UpdateUser(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	if !caller.IsAdmin && caller.UserID != uint(userID) {
		return fiber.ErrForbidden
	}

	var req services.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	userService := services.NewUserService(repo.NewUsersRepo(), repo.NewProfilesRepo())
	if err := userService.UpdateUser(userID, &req); err != nil {
		if errors.Is(err, services.ErrDuplicateEntry) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusNoContent)
}
