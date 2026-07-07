package services

import (
	"Lejematch/internal/database/models"
	"Lejematch/internal/database/repo"
	"errors"
	"strings"
)

var ErrInvalidTargetType = errors.New("invalid target type")
var ErrCannotReportSelf = errors.New("cannot report your own listing")

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
	reportRepo  *repo.ReportsRepo
	listingRepo *repo.ListingsRepo
	seekerRepo  *repo.SeekersRepo
}

func NewReportService(reportRepo *repo.ReportsRepo, listingRepo *repo.ListingsRepo, seekerRepo *repo.SeekersRepo) ReportService {
	return &reportService{reportRepo: reportRepo, listingRepo: listingRepo, seekerRepo: seekerRepo}
}

func (s *reportService) Create(reporterID uint, req *CreateReportRequest) (*models.Report, error) {
	switch models.ReportTargetType(req.TargetType) {
	case models.ReportTargetListing:
		if listing, err := s.listingRepo.FindByID(int(req.TargetID)); err == nil && listing.UserID == reporterID {
			return nil, ErrCannotReportSelf
		}
	case models.ReportTargetSeeker:
		if seeker, err := s.seekerRepo.FindByID(int(req.TargetID)); err == nil && seeker.UserID == reporterID {
			return nil, ErrCannotReportSelf
		}
	case models.ReportTargetProfile:
		if req.TargetID == reporterID {
			return nil, ErrCannotReportSelf
		}
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
