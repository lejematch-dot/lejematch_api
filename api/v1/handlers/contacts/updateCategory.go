package contacts

import (
	"Lejematch/internal/database/models"
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UpdateCategoryRequest er strukturen for opdatering af rød/gul/grøn-markering.
type UpdateCategoryRequest struct {
	Category models.ContactCategory `json:"category"`
}

// UpdateCategory sætter modtagerens private rød/gul/grøn-markering på en
// besked. Kun beskedens modtager kan ændre den — den er ikke synlig for
// afsenderen eller andre.
func UpdateCategory(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	contactID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	var req UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if !req.Category.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ugyldig kategori"})
	}

	contactsRepo := repo.NewContactsRepo()
	contact, err := contactsRepo.FindByID(contactID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}

	if contact.RecipientID != caller.UserID {
		return fiber.ErrForbidden
	}

	contact.Category = req.Category
	if err := contactsRepo.Update(contact); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{"success": true})
}
