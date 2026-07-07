package repo

import (
	"Lejematch/internal/database"
	"Lejematch/internal/database/models"
)

type ReportsRepo struct {
	*GenericRepo[models.Report]
}

func NewReportsRepo() *ReportsRepo {
	return &ReportsRepo{NewGenericRepo[models.Report](database.DB)}
}

// FindAllOrdered returner alle rapporter, nyeste først — pending før resolved.
func (r *ReportsRepo) FindAllOrdered() ([]*models.Report, error) {
	var reports []*models.Report
	err := r.db.Order("status ASC, created_at DESC").Find(&reports).Error
	return reports, err
}
