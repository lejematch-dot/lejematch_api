package users

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func UpdateProfile(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	if !caller.IsAdmin && caller.UserID != uint(userID) {
		return fiber.ErrForbidden
	}

	var req services.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	userService := services.NewUserService(repo.NewUsersRepo(), repo.NewProfilesRepo())
	if err := userService.UpdateProfile(userID, &req); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusNoContent)
}
