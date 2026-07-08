package services

import (
	"Lejematch/internal/citynorm"
	"Lejematch/internal/database/models"
	"Lejematch/internal/database/repo"
	"strings"
)

const minSeekerImages = 3

type CreateSeekerRequest struct {
	Title       string
	Description string
	City        string
	MaxBudget   int
	RoomType    string
	MoveInFrom  string
	Images      []string

	SeekingType         string
	NumPeople           *int
	NumRooms            *int
	FurnishedPreference string
	RentalPeriod        string
	RentalPeriodDetails string
}

type UpdateSeekerRequest struct {
	Title       *string
	Description *string
	City        *string
	MaxBudget   *int
	RoomType    *string
	Status      *string
	MoveInFrom  *string
	Images      []string

	SeekingType         *string
	NumPeople           *int
	NumRooms            *int
	FurnishedPreference *string
	RentalPeriod        *string
	RentalPeriodDetails *string
}

type SeekerService interface {
	Create(userID uint, req *CreateSeekerRequest) (*models.SeekerListing, error)
	Update(seekerID int, callerID uint, isAdmin bool, req *UpdateSeekerRequest) error
	Delete(seekerID int, callerID uint, isAdmin bool) error
}

type seekerService struct {
	seekerRepo *repo.SeekersRepo
}

func NewSeekerService(seekerRepo *repo.SeekersRepo) SeekerService {
	return &seekerService{seekerRepo: seekerRepo}
}

func (s *seekerService) Create(userID uint, req *CreateSeekerRequest) (*models.SeekerListing, error) {
	if len(req.Images) < minSeekerImages {
		return nil, ErrTooFewImages
	}

	seeker := &models.SeekerListing{
		UserID:      userID,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		City:        citynorm.Normalize(req.City),
		MaxBudget:   req.MaxBudget,
		RoomType:    models.RoomType(req.RoomType),
		Status:      models.ListingStatusActive,
		MoveInFrom:  req.MoveInFrom,
		Images:      req.Images,

		SeekingType:         strings.TrimSpace(req.SeekingType),
		NumPeople:           req.NumPeople,
		NumRooms:            req.NumRooms,
		FurnishedPreference: strings.TrimSpace(req.FurnishedPreference),
		RentalPeriod:        strings.TrimSpace(req.RentalPeriod),
		RentalPeriodDetails: strings.TrimSpace(req.RentalPeriodDetails),
	}

	if err := s.seekerRepo.Create(seeker); err != nil {
		return nil, err
	}
	return seeker, nil
}

func (s *seekerService) Update(seekerID int, callerID uint, isAdmin bool, req *UpdateSeekerRequest) error {
	seeker, err := s.seekerRepo.FindByID(seekerID)
	if err != nil {
		return err
	}

	if !isAdmin && seeker.UserID != callerID {
		return ErrNotOwner
	}

	if req.Images != nil && len(req.Images) < minSeekerImages {
		return ErrTooFewImages
	}

	fields := make(map[string]interface{})
	if req.Title != nil {
		fields["title"] = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		fields["description"] = strings.TrimSpace(*req.Description)
	}
	if req.City != nil {
		fields["city"] = citynorm.Normalize(*req.City)
	}
	if req.MaxBudget != nil {
		fields["max_budget"] = *req.MaxBudget
	}
	if req.RoomType != nil {
		fields["room_type"] = *req.RoomType
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	if req.MoveInFrom != nil {
		fields["move_in_from"] = *req.MoveInFrom
	}
	if req.Images != nil {
		fields["images"] = models.StringSlice(req.Images)
	}
	if req.SeekingType != nil {
		fields["seeking_type"] = strings.TrimSpace(*req.SeekingType)
	}
	if req.NumPeople != nil {
		fields["num_people"] = req.NumPeople
	}
	if req.NumRooms != nil {
		fields["num_rooms"] = req.NumRooms
	}
	if req.FurnishedPreference != nil {
		fields["furnished_preference"] = strings.TrimSpace(*req.FurnishedPreference)
	}
	if req.RentalPeriod != nil {
		fields["rental_period"] = strings.TrimSpace(*req.RentalPeriod)
	}
	if req.RentalPeriodDetails != nil {
		fields["rental_period_details"] = strings.TrimSpace(*req.RentalPeriodDetails)
	}

	if len(fields) == 0 {
		return nil
	}

	return s.seekerRepo.UpdateFields(seekerID, fields)
}

func (s *seekerService) Delete(seekerID int, callerID uint, isAdmin bool) error {
	seeker, err := s.seekerRepo.FindByID(seekerID)
	if err != nil {
		return err
	}

	if !isAdmin && seeker.UserID != callerID {
		return ErrNotOwner
	}

	return s.seekerRepo.Delete(seeker)
}
