package auth

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"

	"github.com/gofiber/fiber/v2"
)

func ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Adgangskoden skal være mindst 8 tegn",
		})
	}

	usersRepo := repo.NewUsersRepo()
	if err := services.ResetPassword(usersRepo, req.Token, req.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Linket er ugyldigt eller udløbet",
		})
	}

	return c.JSON(fiber.Map{"success": true})
}
