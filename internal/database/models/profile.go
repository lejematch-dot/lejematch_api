package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model

	UserID uint `gorm:"not null;unique"`
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	DisplayName string
	Bio         string
	City        string
	ImageURL    string
	Phone       string
	Email       string
}
