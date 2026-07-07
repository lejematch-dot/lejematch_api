package contacts

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"

	"github.com/gofiber/fiber/v2"
)

func ListContacts(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	contactsRepo := repo.NewContactsRepo()
	contacts, err := contactsRepo.FindByRecipient(caller.UserID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(contacts)
}
