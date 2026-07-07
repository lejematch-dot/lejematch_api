package models

import "gorm.io/gorm"

type FavoriteType string

const (
	FavoriteTypeListing FavoriteType = "listing"
	FavoriteTypeSeeker  FavoriteType = "seeker"
	FavoriteTypeProfile FavoriteType = "profile"
)

// Favorite lader en bruger gemme en reference til et Listing, SeekerListing
// eller en Profile, skelnet via FavoriteType.
type Favorite struct {
	gorm.Model

	UserID       uint         `gorm:"not null;uniqueIndex:idx_user_favorite"`
	User         User         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	FavoriteType FavoriteType `gorm:"not null;uniqueIndex:idx_user_favorite"`
	FavoriteID   uint         `gorm:"not null;uniqueIndex:idx_user_favorite"`
}
