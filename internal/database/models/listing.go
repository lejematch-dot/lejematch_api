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
	RoomTypePrivate RoomType = "private"
	RoomTypeShared  RoomType = "shared"
	RoomTypeApartment RoomType = "apartment"
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

	UserID      uint          `gorm:"not null;index"`
	User        User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Title       string        `gorm:"not null"`
	Description string
	Price       int           `gorm:"not null"` // DKK per month
	City        string        `gorm:"not null;index"`
	Zip         string        `gorm:"not null"`
	Area        string        // street/neighbourhood, no house number
	RoomType    RoomType      `gorm:"not null"`
	Status      ListingStatus `gorm:"not null;default:'active';index"`
	AvailableFrom string       // ISO date string e.g. "2024-06-01"
	Images        StringSlice  `gorm:"type:jsonb"`
	PromotedUntil *time.Time   `gorm:"index"`
}
