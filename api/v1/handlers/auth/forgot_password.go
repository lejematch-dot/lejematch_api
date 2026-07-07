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
	// Svar altid succes uanset om e-mailen findes, for ikke at afsløre
	// hvilke e-mails er registreret — men log fejl, så en fejlkonfigureret
	// mailer ikke fejler stille.
	if err := services.RequestPasswordReset(usersRepo, strings.ToLower(strings.TrimSpace(req.Email))); err != nil {
		log.Printf("failed to send password reset email: %v", err)
	}

	return c.JSON(fiber.Map{"success": true})
}
