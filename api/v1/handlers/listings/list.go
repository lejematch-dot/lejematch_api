package listings

import (
	"Lejematch/internal/database/repo"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const defaultPageSize = 20

// splitCSV splitter en kommasepareret query-param til en liste, så flere
// filterværdier kan vælges samtidig (f.eks. "landlordType=boligselskab,privat").
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	values := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			values = append(values, p)
		}
	}
	return values
}

func ListListings(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	minPrice, _ := strconv.Atoi(c.Query("minPrice", "0"))
	maxPrice, _ := strconv.Atoi(c.Query("maxPrice", "0"))

	filters := repo.ListingFilters{
		City:                c.Query("city"),
		MinPrice:            minPrice,
		MaxPrice:            maxPrice,
		RoomType:            c.Query("roomType"),
		LandlordType:        splitCSV(c.Query("landlordType")),
		FurnishedPreference: splitCSV(c.Query("furnishedPreference")),
		ListingKind:         splitCSV(c.Query("listingKind")),
		RentalPeriod:        splitCSV(c.Query("rentalPeriod")),
		Category:            c.Query("category"),
		Page:                page,
		PageSize:            defaultPageSize,
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
