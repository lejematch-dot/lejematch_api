package auth

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ForgotPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	usersRepo := repo.NewUsersRepo()
	// Afslører bevidst om kontoen findes (fravalgt anti-enumeration-praksis
	// efter aftale, til fordel for en tydelig brugeroplevelse) — men log
	// mail-fejl, så en fejlkonfigureret mailer ikke fejler stille.
	exists, err := services.RequestPasswordReset(usersRepo, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		log.Printf("failed to send password reset email: %v", err)
	}

	return c.JSON(fiber.Map{"success": true, "accountExists": exists})
}
