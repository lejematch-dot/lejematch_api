package models

import (
	"gorm.io/gorm"
)

// SeekerListing er et opslag fra en LEJER (person der søger bolig),
// modsat Listing som er et opslag fra en UDLEJER.
type SeekerListing struct {
	gorm.Model

	UserID      uint          `gorm:"not null;index"`
	User        User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Title       string        `gorm:"not null"` // f.eks. "Studerende søger værelse i København"
	Description string        // om personen/personerne + hvad de søger
	City        string        `gorm:"not null;index"` // ønsket by
	MaxBudget   int           `gorm:"not null"`       // DKK per måned, max de vil betale
	RoomType    RoomType      `gorm:"not null"`       // ønsket boligtype: private/shared/apartment
	Status      ListingStatus `gorm:"not null;default:'active';index"`
	MoveInFrom  string        // ISO date string, f.eks. "2024-08-01"
	Images      StringSlice   `gorm:"type:jsonb"` // profilbilleder

	SeekingType         string // "bolig" | "roommate" | "begge"
	NumPeople           *int
	NumRooms            *int   // antal værelser søgt, kun relevant når SeekingType er "roommate" eller "begge"
	FurnishedPreference string `gorm:"index"` // "furnished" | "unfurnished" | "any"
	RentalPeriod        string `gorm:"index"` // ønsket lejeperiode: "unlimited" | "limited"
	RentalPeriodDetails string // fri tekst, kun relevant når RentalPeriod = "limited"
}
