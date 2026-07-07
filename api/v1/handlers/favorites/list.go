package favorites

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"

	"github.com/gofiber/fiber/v2"
)

func ListFavorites(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	favoritesRepo := repo.NewFavoritesRepo()
	favorites, err := favoritesRepo.FindByUser(caller.UserID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(favorites)
}
