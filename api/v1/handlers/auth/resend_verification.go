package auth

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ResendVerification(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	usersRepo := repo.NewUsersRepo()
	err := services.ResendVerification(usersRepo, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		if errors.Is(err, services.ErrAlreadyVerified) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "E-mailen er allerede bekræftet"})
		}
		log.Printf("failed to resend verification email: %v", err)
	}

	// Svar altid succes ellers, uanset om e-mailen findes.
	return c.JSON(fiber.Map{"success": true})
}
