package favorites

import (
	"Lejematch/internal/database/models"
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CreateFavoriteRequest struct {
	FavoriteType string `json:"FavoriteType"`
	FavoriteID   uint   `json:"FavoriteID"`
}

func isDuplicateFavorite(err error) bool {
	return strings.Contains(err.Error(), "23505")
}

func CreateFavorite(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	var req CreateFavoriteRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	switch req.FavoriteType {
	case string(models.FavoriteTypeListing), string(models.FavoriteTypeSeeker), string(models.FavoriteTypeProfile):
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid FavoriteType"})
	}
	if req.FavoriteID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "FavoriteID is required"})
	}

	favoritesRepo := repo.NewFavoritesRepo()
	favorite := &models.Favorite{
		UserID:       caller.UserID,
		FavoriteType: models.FavoriteType(req.FavoriteType),
		FavoriteID:   req.FavoriteID,
	}

	if err := favoritesRepo.Create(favorite); err != nil {
		if isDuplicateFavorite(err) {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "alreadyFavorited": true})
		}
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        favorite.ID,
		"createdAt": favorite.CreatedAt,
	})
}
