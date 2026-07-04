package listings

import (
	"Lejematch/internal/database/repo"
	"errors"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gorm.io/gorm"
)

// ContactListingRequest er strukturen for kontakt request
type ContactListingRequest struct {
	Message     string `json:"message" binding:"required"`
	SenderName  string `json:"senderName" binding:"required"`
	SenderEmail string `json:"senderEmail" binding:"required"`
	SenderPhone string `json:"senderPhone"`
}

// ContactListing handler kontakt til udlejer
func ContactListing(c *fiber.Ctx) error {
	// Parse listing ID
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Parse request body
	var req ContactListingRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate request
	if req.Message == "" || req.SenderName == "" || req.SenderEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "message, senderName, and senderEmail are required",
		})
	}

	// Get listing with user (udlejer)
	listingsRepo := repo.NewListingsRepo()
	listing, err := listingsRepo.FindByID(listingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}

	// Get udlejer's info
	usersRepo := repo.NewUsersRepo()
	udlejer, err := usersRepo.FindByID(int(listing.UserID))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	// Send email to udlejer
	if err := sendContactEmail(udlejer.Email, udlejer.FirstName, listing.Title, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Email sent to listing owner",
	})
}

// sendContactEmail sender email til udlejer
func sendContactEmail(recipientEmail, recipientName, listingTitle string, req ContactListingRequest) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return errors.New("SENDGRID_API_KEY not set in environment")
	}

	from := mail.NewEmail("LejeMatch", "noreply@lejematch.dk")
	subject := "Ny interesse i din annonce: " + listingTitle
	to := mail.NewEmail(recipientName, recipientEmail)

	// HTML email body
	htmlContent := `
	<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<h2>Ny interesse i din annonce</h2>
			<p>Hej ` + recipientName + `,</p>
			<p>Der er nogen interesseret i din annonce <strong>` + listingTitle + `</strong>.</p>
			
			<h3>Kontaktinfo fra interessent:</h3>
			<p><strong>Navn:</strong> ` + req.SenderName + `</p>
			<p><strong>Email:</strong> <a href="mailto:` + req.SenderEmail + `">` + req.SenderEmail + `</a></p>
	`

	if req.SenderPhone != "" {
		htmlContent += `<p><strong>Telefon:</strong> ` + req.SenderPhone + `</p>`
	}

	htmlContent += `
			<h3>Besked:</h3>
			<p>` + req.Message + `</p>
			
			<hr>
			<p style="color: #666; font-size: 12px;">
				Denne email er sendt via LejeMatch. 
				Kontakt interessenten direkte på ` + req.SenderEmail + ` eller ` + req.SenderPhone + `.
			</p>
		</body>
	</html>
	`

	message := mail.NewSingleEmail(from, subject, to, "Se email i HTML format", htmlContent)
	client := sendgrid.NewSendClient(apiKey)

	_, err := client.Send(message)
	return err
}
