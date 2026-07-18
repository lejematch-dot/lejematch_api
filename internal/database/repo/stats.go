package repo

import (
	"Lejematch/internal/database"
	"Lejematch/internal/database/models"
	"time"
)

type CityCount struct {
	City  string
	Count int64
}

type Stats struct {
	TotalUsers  int64
	ActiveUsers int64
	NewUsers7d  int64
	NewUsers30d int64

	TotalListings    int64
	ActiveListings   int64
	RentedListings   int64
	ArchivedListings int64
	NewListings7d    int64

	TotalSeekers    int64
	ActiveSeekers   int64
	ArchivedSeekers int64
	NewSeekers7d    int64

	TotalContacts int64
	Contacts7d    int64
	Contacts30d   int64

	TopCities []CityCount
}

type StatsRepo struct{}

func NewStatsRepo() *StatsRepo {
	return &StatsRepo{}
}

func (r *StatsRepo) Get() (*Stats, error) {
	db := database.DB
	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, -7)
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	var s Stats

	if err := db.Model(&models.User{}).Count(&s.TotalUsers).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.User{}).Where("is_active = ?", true).Count(&s.ActiveUsers).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.User{}).Where("created_at >= ?", sevenDaysAgo).Count(&s.NewUsers7d).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.User{}).Where("created_at >= ?", thirtyDaysAgo).Count(&s.NewUsers30d).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.Listing{}).Count(&s.TotalListings).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.Listing{}).Where("status = ?", models.ListingStatusActive).Count(&s.ActiveListings).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.Listing{}).Where("status = ?", models.ListingStatusRented).Count(&s.RentedListings).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.Listing{}).Where("status = ?", models.ListingStatusArchived).Count(&s.ArchivedListings).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.Listing{}).Where("created_at >= ?", sevenDaysAgo).Count(&s.NewListings7d).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.SeekerListing{}).Count(&s.TotalSeekers).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.SeekerListing{}).Where("status = ?", models.ListingStatusActive).Count(&s.ActiveSeekers).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.SeekerListing{}).Where("status = ?", models.ListingStatusArchived).Count(&s.ArchivedSeekers).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.SeekerListing{}).Where("created_at >= ?", sevenDaysAgo).Count(&s.NewSeekers7d).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.Contact{}).Count(&s.TotalContacts).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.Contact{}).Where("created_at >= ?", sevenDaysAgo).Count(&s.Contacts7d).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.Contact{}).Where("created_at >= ?", thirtyDaysAgo).Count(&s.Contacts30d).Error; err != nil {
		return nil, err
	}

	var topCities []CityCount
	if err := db.Model(&models.Listing{}).
		Select("city, count(*) as count").
		Group("city").
		Order("count DESC").
		Limit(5).
		Scan(&topCities).Error; err != nil {
		return nil, err
	}
	s.TopCities = topCities

	return &s, nil
}
