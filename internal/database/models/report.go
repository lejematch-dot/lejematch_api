package models

import "gorm.io/gorm"

type ReportTargetType string

const (
	ReportTargetListing ReportTargetType = "listing"
	ReportTargetSeeker  ReportTargetType = "seeker"
	ReportTargetProfile ReportTargetType = "profile"
)

type ReportStatus string

const (
	ReportStatusPending  ReportStatus = "pending"
	ReportStatusResolved ReportStatus = "resolved"
)

// Report er en anmeldelse af et opslag eller en profil, indsendt af en
// logget ind bruger. Gennemgås manuelt af en admin via /admin/reports.
type Report struct {
	gorm.Model

	ReporterID uint             `gorm:"not null;index"`
	TargetType ReportTargetType `gorm:"not null;index"` // "listing" | "seeker" | "profile"
	TargetID   uint             `gorm:"not null;index"`
	Reason     string           `gorm:"not null"` // "spam" | "svindel" | "andet"
	Message    string
	Status     ReportStatus `gorm:"not null;default:'pending';index"`
}
