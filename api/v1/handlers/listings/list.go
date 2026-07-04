package listings

import (
	"Lejematch/internal/database/repo"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const defaultPageSize = 20

func ListListings(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	minPrice, _ := strconv.Atoi(c.Query("minPrice", "0"))
	maxPrice, _ := strconv.Atoi(c.Query("maxPrice", "0"))

	filters := repo.ListingFilters{
		City:     c.Query("city"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		RoomType: c.Query("roomType"),
		Page:     page,
		PageSize: defaultPageSize,
	}

	listingsRepo := repo.NewListingsRepo()
	results, total, err := listingsRepo.FindFiltered(filters)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"data":       results,
		"page":       page,
		"pageSize":   defaultPageSize,
		"total":      total,
		"totalPages": int(math.Ceil(float64(total) / float64(defaultPageSize))),
	})
}
