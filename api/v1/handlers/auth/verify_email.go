package auth

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"

	"github.com/gofiber/fiber/v2"
)

func VerifyEmail(c *fiber.Ctx) error {
	token := c.Params("token")

	usersRepo := repo.NewUsersRepo()
	if err := services.VerifyEmail(usersRepo, token); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Linket er ugyldigt eller udløbet",
		})
	}

	return c.JSON(fiber.Map{"success": true})
}
