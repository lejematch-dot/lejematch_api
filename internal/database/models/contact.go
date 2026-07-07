package models

import "gorm.io/gorm"

// ContactTargetType angiver om beskeden vedrører et Listing eller en SeekerListing.
type ContactTargetType string

const (
	ContactTargetListing ContactTargetType = "listing"
	ContactTargetSeeker  ContactTargetType = "seeker"
)

// Contact er en besked sendt fra én bruger (Sender) til en anden (Recipient)
// om et konkret opslag. Gemmes i databasen ud over at blive sendt som e-mail,
// så modtageren kan se historikken i sit dashboard.
type Contact struct {
	gorm.Model

	SenderID    uint              `gorm:"not null;index"`
	Sender      User              `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	RecipientID uint              `gorm:"not null;index"`
	Recipient   User              `gorm:"foreignKey:RecipientID;constraint:OnDelete:CASCADE"`
	TargetType  ContactTargetType `gorm:"not null"`
	TargetID    uint              `gorm:"not null"`
	Message     string            `gorm:"not null"`
	SenderPhone string
}
