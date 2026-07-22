package newsletter

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"

	"github.com/gofiber/fiber/v2"
)

// Subscribe slår nyhedsbrev-tilmelding til for brugeren tokenet tilhører.
func Subscribe(c *fiber.Ctx) error {
	token := c.Params("token")

	usersRepo := repo.NewUsersRepo()
	if err := services.SubscribeNewsletter(usersRepo, token); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Linket er ugyldigt eller udløbet"})
	}

	return c.JSON(fiber.Map{"success": true})
}
