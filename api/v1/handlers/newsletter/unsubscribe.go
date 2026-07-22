package newsletter

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"

	"github.com/gofiber/fiber/v2"
)

// Unsubscribe slår nyhedsbrev-tilmelding fra for brugeren tokenet tilhører.
func Unsubscribe(c *fiber.Ctx) error {
	token := c.Params("token")

	usersRepo := repo.NewUsersRepo()
	if err := services.UnsubscribeNewsletter(usersRepo, token); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Linket er ugyldigt eller udløbet"})
	}

	return c.JSON(fiber.Map{"success": true})
}
