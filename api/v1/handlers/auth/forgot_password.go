package auth

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
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
	// Fejl ignoreres bevidst — svar altid succes, uanset om e-mailen findes,
	// for ikke at afsløre hvilke e-mails er registreret.
	_ = services.RequestPasswordReset(usersRepo, strings.ToLower(strings.TrimSpace(req.Email)))

	return c.JSON(fiber.Map{"success": true})
}
