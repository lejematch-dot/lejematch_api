package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ListingStatus string

const (
	ListingStatusActive   ListingStatus = "active"
	ListingStatusRented   ListingStatus = "rented"
	ListingStatusArchived ListingStatus = "archived"
)

type RoomType string

const (
	RoomTypePrivate   RoomType = "private"
	RoomTypeShared    RoomType = "shared"
	RoomTypeApartment RoomType = "apartment"
)

// ListingType beskriver boligens størrelse/type (antal værelser), som supplement til RoomType.
type ListingType string

const (
	ListingTypeRoom  ListingType = "room"
	ListingType1V    ListingType = "1v"
	ListingType2V    ListingType = "2v"
	ListingType3V    ListingType = "3v"
	ListingType4V    ListingType = "4v"
	ListingType5V    ListingType = "5v"
	ListingTypeHouse ListingType = "house"
)

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringSlice")
	}
	return json.Unmarshal(bytes, s)
}

type Listing struct {
	gorm.Model

	UserID        uint   `gorm:"not null;index"`
	User          User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Title         string `gorm:"not null"`
	Description   string
	Price         int           `gorm:"not null"` // DKK per month
	City          string        `gorm:"not null;index"`
	Zip           string        `gorm:"not null"`
	Area          string        // street/neighbourhood, no house number
	RoomType      RoomType      `gorm:"not null"`
	Status        ListingStatus `gorm:"not null;default:'active';index"`
	AvailableFrom string        // ISO date string e.g. "2024-06-01"
	Images        StringSlice   `gorm:"type:jsonb"`
	PromotedUntil *time.Time    `gorm:"index"`

	ListingKind         ListingType `gorm:"index"` // room/1v/2v/3v/4v/5v/house — supplerer RoomType
	SizeSqm             *int
	Deposit             *int
	RentalPeriod        string      `gorm:"index"` // "unlimited" | "limited"
	RentalPeriodDetails string      // fri tekst, kun relevant når RentalPeriod = "limited"
	LandlordType        string      `gorm:"index"` // "boligselskab" | "privat"
	FurnishedPreference string      `gorm:"index"` // "furnished" | "unfurnished" | "any"
	Facilities          StringSlice `gorm:"type:jsonb"`
	TargetAudience      string
	RoommatesWanted     *int // antal nye roomies søgt, kun relevant når ListingKind = "room"
}
