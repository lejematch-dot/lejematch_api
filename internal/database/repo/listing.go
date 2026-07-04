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
