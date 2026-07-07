package favorites

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func DeleteFavorite(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	favoriteType := c.Query("favoriteType")
	favoriteID, err := strconv.Atoi(c.Query("favoriteId"))
	if favoriteType == "" || err != nil {
		return fiber.ErrBadRequest
	}

	favoritesRepo := repo.NewFavoritesRepo()
	if err := favoritesRepo.DeleteByUserAndTarget(caller.UserID, favoriteType, uint(favoriteID)); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusNoContent)
}
