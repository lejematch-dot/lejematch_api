package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	FirstName       string
	LastName        string
	Email           string `gorm:"unique"`
	Phone           string `gorm:"unique"`
	Password        string
	IsAdmin         bool
	IsActive        bool
	NewsletterOptIn bool
}
