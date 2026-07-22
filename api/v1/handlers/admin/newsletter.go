package admin

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

type SendNewsletterRequest struct {
	Subject string `json:"subject"`
	Html    string `json:"html"`
}

// SendNewsletterHandler sender et nyhedsbrev til alle tilmeldte brugere.
// Selve afsendelsen sker i en baggrunds-goroutine, så requesten ikke venter
// på at hver enkelt mail bliver sendt.
func SendNewsletterHandler(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)
	if !caller.IsAdmin {
		return fiber.ErrForbidden
	}

	var req SendNewsletterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if req.Subject == "" || req.Html == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "subject og html er påkrævet"})
	}

	usersRepo := repo.NewUsersRepo()
	recipients, err := usersRepo.FindNewsletterRecipients()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	go func() {
		sent, failed := services.SendNewsletter(usersRepo, req.Subject, req.Html)
		log.Printf("newsletter: send complete — %d sent, %d failed", sent, failed)
	}()

	return c.JSON(fiber.Map{
		"success":        true,
		"recipientCount": len(recipients),
	})
}

// SendNewsletterInviteHandler sender en engangs-invitation til alle
// eksisterende brugere der endnu ikke er tilmeldt nyhedsbreve, og spørger
// om de vil tilmelde sig. Sikker at køre flere gange — rammer kun dem der
// stadig mangler at tage stilling.
func SendNewsletterInviteHandler(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)
	if !caller.IsAdmin {
		return fiber.ErrForbidden
	}

	usersRepo := repo.NewUsersRepo()
	targets, err := usersRepo.FindNewsletterInviteTargets()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	go func() {
		sent, failed := services.SendNewsletterInvite(usersRepo)
		log.Printf("newsletter invite: send complete — %d sent, %d failed", sent, failed)
	}()

	return c.JSON(fiber.Map{
		"success":     true,
		"targetCount": len(targets),
	})
}
