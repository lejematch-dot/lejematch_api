package database

import (
	"Lejematch/internal/database/models"
	"Lejematch/internal/security"
	"time"
)

func Seed() {
	// Skip if data already exists
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	password, _ := security.HashPassword("password123")

	users := []models.User{
		{FirstName: "Anders", LastName: "Jensen", Email: "lejematch@gmail.com", Phone: "11111111", Password: password, IsActive: true},
		{FirstName: "Sofie", LastName: "Nielsen", Email: "sofie@example.com", Phone: "22222222", Password: password, IsActive: true},
		{FirstName: "Mikkel", LastName: "Hansen", Email: "mikkel@example.com", Phone: "33333333", Password: password, IsActive: true},
		{FirstName: "Laura", LastName: "Pedersen", Email: "laura@example.com", Phone: "44444444", Password: password, IsActive: true},
	}

	for i := range users {
		DB.Create(&users[i])

		DB.Create(&models.Profile{
			UserID:      users[i].ID,
			DisplayName: users[i].FirstName + " " + users[i].LastName,
			City:        []string{"København", "Aarhus", "Odense", "Aalborg"}[i],
			Bio:         "Hej! Jeg søger en god roommate.",
			Email:       users[i].Email,
			Phone:       users[i].Phone,
		})
	}

	promotedUntil := time.Now().Add(7 * 24 * time.Hour)

	listings := []models.Listing{
		{
			UserID:        users[0].ID,
			Title:         "Lyst værelse i hjertet af København",
			Description:   "Et møbleret værelse i en 3-værelses lejlighed. Delt køkken og badeværelse. Tæt på metro.",
			Price:         6500,
			City:          "København",
			Zip:           "2200",
			Area:          "Nørrebro",
			RoomType:      models.RoomTypePrivate,
			Status:        models.ListingStatusActive,
			AvailableFrom: "2024-06-01",
			Images:        models.StringSlice{"https://images.unsplash.com/photo-1522708323590-d24dbb6b0267?w=800&h=600&fit=crop"},
			PromotedUntil: &promotedUntil,
			ListingKind:   models.ListingTypeRoom,
			SizeSqm:       intPtr(14),
		},
		{
			UserID:        users[1].ID,
			Title:         "Hyggeligt værelse nær Aarhus C",
			Description:   "Roligt værelse i lejlighed med 2 andre studerende. Gode forbindelser til universitetet.",
			Price:         4800,
			City:          "Aarhus",
			Zip:           "8200",
			Area:          "Trøjborg",
			RoomType:      models.RoomTypePrivate,
			Status:        models.ListingStatusActive,
			AvailableFrom: "2024-07-01",
			Images:        models.StringSlice{"https://images.unsplash.com/photo-1505691938895-1758d7feb511?w=800&h=600&fit=crop"},
			ListingKind:   models.ListingTypeRoom,
			SizeSqm:       intPtr(18),
		},
		{
			UserID:        users[2].ID,
			Title:         "Hel lejlighed til leje — Odense",
			Description:   "Moderne 2-værelses lejlighed med altan. Perfekt for par eller to venner.",
			Price:         8200,
			City:          "Odense",
			Zip:           "5000",
			Area:          "Østerbro",
			RoomType:      models.RoomTypeApartment,
			Status:        models.ListingStatusActive,
			AvailableFrom: "2024-06-15",
			Images:        models.StringSlice{"https://images.unsplash.com/photo-1502672260266-1c1ef2d93688?w=800&h=600&fit=crop"},
			ListingKind:   models.ListingType2V,
			SizeSqm:       intPtr(58),
		},
		{
			UserID:        users[3].ID,
			Title:         "Delt værelse i Aalborg",
			Description:   "Vi er to venner der deler en stor lejlighed og søger en tredje roommate. Meget socialt.",
			Price:         3200,
			City:          "Aalborg",
			Zip:           "9000",
			Area:          "Midtbyen",
			RoomType:      models.RoomTypeShared,
			Status:        models.ListingStatusActive,
			AvailableFrom: "2024-05-01",
			Images:        models.StringSlice{"https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=800&h=600&fit=crop"},
			ListingKind:   models.ListingTypeRoom,
			SizeSqm:       intPtr(15),
		},
		{
			UserID:        users[0].ID,
			Title:         "Værelse udlejes — Frederiksberg",
			Description:   "Stort lyst værelse i rolig gade. Eget badeværelse. Ingen husdyr.",
			Price:         7000,
			City:          "København",
			Zip:           "2000",
			Area:          "Frederiksberg",
			RoomType:      models.RoomTypePrivate,
			Status:        models.ListingStatusActive,
			AvailableFrom: "2024-08-01",
			Images:        models.StringSlice{"https://images.unsplash.com/photo-1493809842364-78817add7ffb?w=800&h=600&fit=crop"},
			ListingKind:   models.ListingType1V,
			SizeSqm:       intPtr(35),
		},
	}

	for i := range listings {
		DB.Create(&listings[i])
	}
}

func intPtr(i int) *int {
	return &i
}
