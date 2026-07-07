package repo

import (
	"Lejematch/internal/database"
	"Lejematch/internal/database/models"
)

type ListingsRepo struct {
	*GenericRepo[models.Listing]
}

func NewListingsRepo() *ListingsRepo {
	return &ListingsRepo{NewGenericRepo[models.Listing](database.DB)}
}

type ListingFilters struct {
	City     string
	MinPrice int
	MaxPrice int
	RoomType string
	// LandlordType, FurnishedPreference, ListingKind og RentalPeriod er
	// lister, så flere værdier kan vælges samtidig (f.eks. både "2v" og "3v").
	LandlordType        []string
	FurnishedPreference []string
	ListingKind         []string
	RentalPeriod        []string
	// Category filtrerer på ListingKind: "vaerelse" = kun enkeltværelser,
	// "hele" (eller tom) = alt andet end enkeltværelser.
	Category string
	Page     int
	PageSize int
}

func (r *ListingsRepo) FindFiltered(f ListingFilters) ([]*models.Listing, int64, error) {
	query := r.db.Model(&models.Listing{}).Where("status = ?", models.ListingStatusActive)

	if f.City != "" {
		query = query.Where("city = ?", f.City)
	}
	if f.MinPrice > 0 {
		query = query.Where("price >= ?", f.MinPrice)
	}
	if f.MaxPrice > 0 {
		query = query.Where("price <= ?", f.MaxPrice)
	}
	if f.RoomType != "" {
		query = query.Where("room_type = ?", f.RoomType)
	}
	if len(f.LandlordType) > 0 {
		query = query.Where("landlord_type IN ?", f.LandlordType)
	}
	if len(f.FurnishedPreference) > 0 {
		query = query.Where("furnished_preference IN ?", f.FurnishedPreference)
	}
	if len(f.RentalPeriod) > 0 {
		query = query.Where("rental_period IN ?", f.RentalPeriod)
	}
	if f.Category == "vaerelse" {
		query = query.Where("listing_kind = ?", "room")
	} else if f.Category == "hele" {
		query = query.Where("listing_kind <> ?", "room")
	}
	if len(f.ListingKind) > 0 {
		query = query.Where("listing_kind IN ?", f.ListingKind)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.PageSize
	var listings []*models.Listing
	err := query.Order("(promoted_until > NOW()) DESC, promoted_until DESC, created_at DESC").Offset(offset).Limit(f.PageSize).Find(&listings).Error
	return listings, total, err
}

func (r *ListingsRepo) FindByUserID(userID uint) ([]*models.Listing, error) {
	var listings []*models.Listing
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&listings).Error
	return listings, err
}

// DistinctCities returner alle unikke byer der har mindst ét aktivt opslag
// inden for den angivne kategori ("hele"/"vaerelse", tom = alle), alfabetisk
// sorteret. Bruges til by-filteret på browse-siden, så en by kun vises når
// den rent faktisk har opslag i den fane brugeren kigger på.
func (r *ListingsRepo) DistinctCities(category string) ([]string, error) {
	query := r.db.Model(&models.Listing{}).Where("status = ?", models.ListingStatusActive)
	if category == "vaerelse" {
		query = query.Where("listing_kind = ?", "room")
	} else if category == "hele" {
		query = query.Where("listing_kind <> ?", "room")
	}

	var cities []string
	err := query.
		Distinct("city").
		Order("city ASC").
		Pluck("city", &cities).Error
	return cities, err
}
