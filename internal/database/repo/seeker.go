package repo

import (
	"Lejematch/internal/database"
	"Lejematch/internal/database/models"
)

type SeekersRepo struct {
	*GenericRepo[models.SeekerListing]
}

func NewSeekersRepo() *SeekersRepo {
	return &SeekersRepo{NewGenericRepo[models.SeekerListing](database.DB)}
}

type SeekerFilters struct {
	City      string
	MaxBudget int
	// RoomType, FurnishedPreference og RentalPeriod er lister, så flere
	// værdier kan vælges samtidig (f.eks. både "private" og "shared").
	RoomType            []string
	FurnishedPreference []string
	RentalPeriod        []string
	// Category filtrerer på SeekingType: "vaerelse" = søger værelse i
	// bofællesskab ("roommate"), "hele" (eller tom) = søger hel bolig.
	Category string
	Page     int
	PageSize int
}

func (r *SeekersRepo) FindFiltered(f SeekerFilters) ([]*models.SeekerListing, int64, error) {
	query := r.db.Model(&models.SeekerListing{}).Where("status = ?", models.ListingStatusActive)

	if f.City != "" {
		query = query.Where("city = ?", f.City)
	}
	if f.MaxBudget > 0 {
		query = query.Where("max_budget <= ?", f.MaxBudget)
	}
	if len(f.RoomType) > 0 {
		query = query.Where("room_type IN ?", f.RoomType)
	}
	if len(f.FurnishedPreference) > 0 {
		query = query.Where("furnished_preference IN ?", f.FurnishedPreference)
	}
	if len(f.RentalPeriod) > 0 {
		query = query.Where("rental_period IN ?", f.RentalPeriod)
	}
	if f.Category == "vaerelse" {
		query = query.Where("seeking_type = ?", "roommate")
	} else if f.Category == "hele" {
		query = query.Where("seeking_type <> ?", "roommate")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.PageSize
	var seekers []*models.SeekerListing
	err := query.Order("created_at DESC").Offset(offset).Limit(f.PageSize).Find(&seekers).Error
	return seekers, total, err
}

func (r *SeekersRepo) FindByUserID(userID uint) ([]*models.SeekerListing, error) {
	var seekers []*models.SeekerListing
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&seekers).Error
	return seekers, err
}

// DistinctCities returner alle unikke byer der har mindst ét aktivt
// lejeropslag inden for den angivne kategori ("hele"/"vaerelse", tom = alle),
// alfabetisk sorteret. Bruges til by-filteret på lejere-siden.
func (r *SeekersRepo) DistinctCities(category string) ([]string, error) {
	query := r.db.Model(&models.SeekerListing{}).Where("status = ?", models.ListingStatusActive)
	if category == "vaerelse" {
		query = query.Where("seeking_type = ?", "roommate")
	} else if category == "hele" {
		query = query.Where("seeking_type <> ?", "roommate")
	}

	var cities []string
	err := query.
		Distinct("city").
		Order("city ASC").
		Pluck("city", &cities).Error
	return cities, err
}
