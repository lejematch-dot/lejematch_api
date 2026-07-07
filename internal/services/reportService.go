package services

import (
	"Lejematch/internal/database/models"
	"Lejematch/internal/database/repo"
	"errors"
	"strings"
)

var ErrInvalidTargetType = errors.New("invalid target type")

type CreateReportRequest struct {
	TargetType string
	TargetID   uint
	Reason     string
	Message    string
}

type ReportService interface {
	Create(reporterID uint, req *CreateReportRequest) (*models.Report, error)
	List() ([]*models.Report, error)
	Resolve(reportID int) error
}

type reportService struct {
	reportRepo *repo.ReportsRepo
}

func NewReportService(reportRepo *repo.ReportsRepo) ReportService {
	return &reportService{reportRepo: reportRepo}
}

func (s *reportService) Create(reporterID uint, req *CreateReportRequest) (*models.Report, error) {
	switch models.ReportTargetType(req.TargetType) {
	case models.ReportTargetListing, models.ReportTargetSeeker, models.ReportTargetProfile:
	default:
		return nil, ErrInvalidTargetType
	}

	report := &models.Report{
		ReporterID: reporterID,
		TargetType: models.ReportTargetType(req.TargetType),
		TargetID:   req.TargetID,
		Reason:     strings.TrimSpace(req.Reason),
		Message:    strings.TrimSpace(req.Message),
		Status:     models.ReportStatusPending,
	}
	if err := s.reportRepo.Create(report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *reportService) List() ([]*models.Report, error) {
	return s.reportRepo.FindAllOrdered()
}

func (s *reportService) Resolve(reportID int) error {
	return s.reportRepo.UpdateFields(reportID, map[string]interface{}{"status": string(models.ReportStatusResolved)})
}
