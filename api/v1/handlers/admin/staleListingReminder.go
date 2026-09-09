package admin

import (
	"Lejematch/internal/services"

	"github.com/gofiber/fiber/v2"
)

// TriggerStaleListingReminders sender "er dit opslag stadig aktuelt?"-mails
// med det samme, i stedet for at vente på kl. 09. Kun for admins.
func TriggerStaleListingReminders(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)
	if !caller.IsAdmin {
		return fiber.ErrForbidden
	}

	services.TriggerStaleListingRemindersNow()

	return c.JSON(fiber.Map{"success": true})
}
