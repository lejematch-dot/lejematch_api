package seekers

import (
	"Lejematch/internal/database/repo"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const defaultPageSize = 20

// splitCSV splitter en kommasepareret query-param til en liste, så flere
// filterværdier kan vælges samtidig (f.eks. "roomType=private,shared").
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

func ListSeekers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	maxBudget, _ := strconv.Atoi(c.Query("maxBudget", "0"))

	filters := repo.SeekerFilters{
		City:                c.Query("city"),
		MaxBudget:           maxBudget,
		RoomType:            splitCSV(c.Query("roomType")),
		FurnishedPreference: splitCSV(c.Query("furnishedPreference")),
		RentalPeriod:        splitCSV(c.Query("rentalPeriod")),
		Category:            c.Query("category"),
		Page:                page,
		PageSize:            defaultPageSize,
	}

	seekersRepo := repo.NewSeekersRepo()
	results, total, err := seekersRepo.FindFiltered(filters)
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
