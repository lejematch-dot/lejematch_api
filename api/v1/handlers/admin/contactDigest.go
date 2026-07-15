package admin

import (
	"Lejematch/internal/services"

	"github.com/gofiber/fiber/v2"
)

// TriggerContactDigest sender den daglige kontakt-oversigt med det samme
// (i stedet for at vente til kl. 22). Kun for admins.
func TriggerContactDigest(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)
	if !caller.IsAdmin {
		return fiber.ErrForbidden
	}

	services.TriggerContactDigestNow()

	return c.JSON(fiber.Map{"success": true})
}
