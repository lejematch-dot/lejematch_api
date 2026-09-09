package contacts

import (
	"Lejematch/internal/database/models"
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// mustBeParticipant henter tråden (Contact) og tjekker at kalderen enten er
// den oprindelige afsender eller modtager — begge kan se og svare i tråden.
func mustBeParticipant(c *fiber.Ctx, callerID uint) (*models.Contact, error) {
	contactID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return nil, fiber.ErrBadRequest
	}

	contactsRepo := repo.NewContactsRepo()
	contact, err := contactsRepo.FindByID(contactID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, fiber.ErrInternalServerError
	}

	if contact.SenderID != callerID && contact.RecipientID != callerID {
		return nil, fiber.ErrForbidden
	}

	return contact, nil
}

// ListReplies henter alle svar i en tråd. Kræver at kalderen er en af de to
// oprindelige parter i tråden.
func ListReplies(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	contact, err := mustBeParticipant(c, caller.UserID)
	if err != nil {
		return err
	}

	repliesRepo := repo.NewContactRepliesRepo()
	replies, err := repliesRepo.FindByContactID(contact.ID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(replies)
}

// CreateReplyRequest er strukturen for et nyt svar i en tråd.
type CreateReplyRequest struct {
	Message string `json:"message"`
}

// CreateReply tilføjer et svar til en tråd. Kalderen skal være en af de to
// oprindelige parter — modtageren af svaret er automatisk "den anden" part.
func CreateReply(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	contact, err := mustBeParticipant(c, caller.UserID)
	if err != nil {
		return err
	}

	var req CreateReplyRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	message := strings.TrimSpace(req.Message)
	if message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "besked er påkrævet"})
	}

	reply := &models.ContactReply{
		ContactID: contact.ID,
		SenderID:  caller.UserID,
		Message:   message,
	}

	repliesRepo := repo.NewContactRepliesRepo()
	if err := repliesRepo.Create(reply); err != nil {
		return fiber.ErrInternalServerError
	}

	// Den anden part i tråden — ikke nødvendigvis Contact.RecipientID, da
	// begge parter kan svare frem og tilbage.
	otherPartyID := contact.RecipientID
	if caller.UserID == contact.RecipientID {
		otherPartyID = contact.SenderID
	}

	// Fejl i mailafsendelse må ikke forhindre at svaret er gemt — det kan
	// stadig ses i modtagerens "Beskeder"-oversigt.
	usersRepo := repo.NewUsersRepo()
	profilesRepo := repo.NewProfilesRepo()
	if otherUser, err := usersRepo.FindByID(int(otherPartyID)); err == nil {
		senderName := "En bruger"
		if senderProfile, err := profilesRepo.FindByUserID(caller.UserID); err == nil {
			senderName = senderProfile.DisplayName
		}
		_ = sendReplyNotificationEmail(otherUser.Email, otherUser.FirstName, senderName, message)
	}

	return c.Status(fiber.StatusCreated).JSON(reply)
}

func sendReplyNotificationEmail(recipientEmail, recipientName, senderName, message string) error {
	subject := "Nyt svar på LejeMatch fra " + senderName

	html := `
	<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			` + services.EmailHeader() + `
			<h2>Nyt svar i din samtale</h2>
			<p>Hej ` + recipientName + `,</p>
			<p><strong>` + senderName + `</strong> har svaret dig på LejeMatch:</p>
			<p style="background: #f5f5f5; padding: 12px 16px; border-radius: 4px;">` + message + `</p>
			<hr>
			<p style="color: #666; font-size: 12px;">
				Log ind for at se og svare under "Beskeder".
			</p>
			` + services.EmailSignature() + `
		</body>
	</html>
	`

	return services.SendEmail(recipientEmail, subject, html)
}
