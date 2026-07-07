package seekers

import (
	"Lejematch/internal/database/models"
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ContactSeekerRequest er strukturen for kontakt request
type ContactSeekerRequest struct {
	Message     string `json:"message" binding:"required"`
	SenderPhone string `json:"senderPhone"`
}

// ContactSeeker handler kontakt til en lejer der søger bolig. Kræver login,
// så beskeden kan knyttes til en rigtig bruger og vises i modtagerens
// "Beskeder"-oversigt.
func ContactSeeker(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	seekerID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	var req ContactSeekerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "message is required"})
	}

	seekersRepo := repo.NewSeekersRepo()
	seeker, err := seekersRepo.FindByID(seekerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}

	if seeker.UserID == caller.UserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Du kan ikke kontakte dig selv"})
	}

	usersRepo := repo.NewUsersRepo()
	lejer, err := usersRepo.FindByID(int(seeker.UserID))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	profilesRepo := repo.NewProfilesRepo()
	senderProfile, err := profilesRepo.FindByUserID(caller.UserID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	contactsRepo := repo.NewContactsRepo()
	contact := &models.Contact{
		SenderID:    caller.UserID,
		RecipientID: seeker.UserID,
		TargetType:  models.ContactTargetSeeker,
		TargetID:    seeker.ID,
		Message:     req.Message,
		SenderPhone: req.SenderPhone,
	}
	if err := contactsRepo.Create(contact); err != nil {
		return fiber.ErrInternalServerError
	}

	// Fejl i mailafsendelse må ikke forhindre at beskeden er gemt.
	_ = sendSeekerContactEmail(lejer.Email, lejer.FirstName, seeker.Title, senderProfile.DisplayName, caller.Email, req.SenderPhone, req.Message)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Besked sendt",
	})
}

// sendSeekerContactEmail sender email til lejeren via den delte mailer
func sendSeekerContactEmail(recipientEmail, recipientName, seekerTitle, senderName, senderEmail, senderPhone, message string) error {
	subject := "Ny henvendelse på dit opslag: " + seekerTitle

	htmlContent := `
	<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<h2>Ny henvendelse på dit opslag</h2>
			<p>Hej ` + recipientName + `,</p>
			<p>Der er nogen der har set dit opslag <strong>` + seekerTitle + `</strong> og vil kontakte dig.</p>

			<h3>Kontaktinfo fra afsender:</h3>
			<p><strong>Navn:</strong> ` + senderName + `</p>
			<p><strong>Email:</strong> <a href="mailto:` + senderEmail + `">` + senderEmail + `</a></p>
	`

	if senderPhone != "" {
		htmlContent += `<p><strong>Telefon:</strong> ` + senderPhone + `</p>`
	}

	htmlContent += `
			<h3>Besked:</h3>
			<p>` + message + `</p>

			<hr>
			<p style="color: #666; font-size: 12px;">
				Denne email er sendt via LejeMatch. Log ind for at se beskeden under "Beskeder".
			</p>
		</body>
	</html>
	`

	return services.SendEmail(recipientEmail, subject, htmlContent)
}
