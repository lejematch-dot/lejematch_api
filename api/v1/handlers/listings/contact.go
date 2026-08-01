package listings

import (
	"Lejematch/internal/database/models"
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ContactListingRequest er strukturen for kontakt request
type ContactListingRequest struct {
	Message          string                         `json:"message" binding:"required"`
	SenderPhone      string                         `json:"senderPhone"`
	NumPeople        int                            `json:"numPeople"`
	RelationshipType models.ContactRelationshipType `json:"relationshipType"`
	Ages             []int                          `json:"ages"`
	Employment       models.ContactEmployment       `json:"employment"`
	HasPets          bool                           `json:"hasPets"`
}

// ContactListing handler kontakt til udlejer. Kræver login, så beskeden kan
// knyttes til en rigtig bruger og vises i modtagerens "Beskeder"-oversigt.
func ContactListing(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	var req ContactListingRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "message is required"})
	}
	if req.NumPeople < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "numPeople skal være mindst 1"})
	}
	// Der skal være én alder pr. person, op til dropdown-loftet på 5 ("5+").
	wantAges := req.NumPeople
	if wantAges > 5 {
		wantAges = 5
	}
	if len(req.Ages) != wantAges {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "angiv en alder pr. person"})
	}
	for _, age := range req.Ages {
		if age < 1 || age > 120 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ugyldig alder"})
		}
	}
	if !req.Employment.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ugyldig beskæftigelse"})
	}
	if !req.RelationshipType.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ugyldig par/venner-værdi"})
	}

	listingsRepo := repo.NewListingsRepo()
	listing, err := listingsRepo.FindByID(listingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}

	if listing.UserID == caller.UserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Du kan ikke kontakte dig selv"})
	}

	usersRepo := repo.NewUsersRepo()
	udlejer, err := usersRepo.FindByID(int(listing.UserID))
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
		SenderID:         caller.UserID,
		RecipientID:      listing.UserID,
		TargetType:       models.ContactTargetListing,
		TargetID:         listing.ID,
		Message:          req.Message,
		SenderPhone:      req.SenderPhone,
		NumPeople:        req.NumPeople,
		RelationshipType: req.RelationshipType,
		Ages:             models.IntSlice(req.Ages),
		Employment:       req.Employment,
		HasPets:          req.HasPets,
	}
	if err := contactsRepo.Create(contact); err != nil {
		return fiber.ErrInternalServerError
	}

	// Fejl i mailafsendelse må ikke forhindre at beskeden er gemt — den kan
	// stadig ses i modtagerens "Beskeder"-oversigt.
	_ = sendContactEmail(udlejer.Email, udlejer.FirstName, listing.Title, senderProfile.DisplayName, caller.Email, req.SenderPhone, req.Message, contact)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Besked sendt",
	})
}

// sendContactEmail sender email til udlejer via den delte mailer
func sendContactEmail(recipientEmail, recipientName, listingTitle, senderName, senderEmail, senderPhone, message string, contact *models.Contact) error {
	subject := "Ny interesse i din annonce: " + listingTitle

	htmlContent := `
	<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			` + services.EmailHeader() + `
			<h2>Ny interesse i din annonce</h2>
			<p>Hej ` + recipientName + `,</p>
			<p>Der er nogen interesseret i din annonce <strong>` + listingTitle + `</strong>.</p>

			<h3>Kontaktinfo fra interessent:</h3>
			<p><strong>Navn:</strong> ` + senderName + `</p>
			<p><strong>Email:</strong> <a href="mailto:` + senderEmail + `">` + senderEmail + `</a></p>
	`

	if senderPhone != "" {
		htmlContent += `<p><strong>Telefon:</strong> ` + senderPhone + `</p>`
	}

	htmlContent += `
			<h3>Om interessenten:</h3>
			<p><strong>Antal personer:</strong> ` + contact.NumPeopleSummary() + `</p>
			<p><strong>Alder:</strong> ` + contact.AgesSummary() + `</p>
			<p><strong>Beskæftigelse:</strong> ` + contact.Employment.Label() + `</p>
			<p><strong>Kæledyr:</strong> ` + contact.PetsLabel() + `</p>
	`

	htmlContent += `
			<h3>Besked:</h3>
			<p>` + message + `</p>

			<hr>
			<p style="color: #666; font-size: 12px;">
				Denne email er sendt via LejeMatch. Log ind for at se beskeden under "Beskeder".
			</p>
			` + services.EmailSignature() + `
		</body>
	</html>
	`

	return services.SendEmail(recipientEmail, subject, htmlContent)
}
