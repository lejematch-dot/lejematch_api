package services

import (
	"Lejematch/internal/database/models"
	"Lejematch/internal/database/repo"
	"errors"
	"strings"
	"time"
)

var ErrNotOwner = errors.New("not the owner of this listing")

type CreateListingRequest struct {
	Title         string
	Description   string
	Price         int
	City          string
	Zip           string
	Area          string
	RoomType      string
	AvailableFrom string
	Images        []string

	ListingKind         string
	SizeSqm             *int
	Deposit             *int
	RentalPeriod        string
	RentalPeriodDetails string
	LandlordType        string
	FurnishedPreference string
	Facilities          []string
	TargetAudience      string
}

type UpdateListingRequest struct {
	Title         *string
	Description   *string
	Price         *int
	City          *string
	Zip           *string
	Area          *string
	RoomType      *string
	Status        *string
	AvailableFrom *string
	Images        []string
	PromotedUntil *time.Time // admin only

	ListingKind         *string
	SizeSqm             *int
	Deposit             *int
	RentalPeriod        *string
	RentalPeriodDetails *string
	LandlordType        *string
	FurnishedPreference *string
	Facilities          []string
	TargetAudience      *string
}

type ListingService interface {
	Create(userID uint, req *CreateListingRequest) (*models.Listing, error)
	Update(listingID int, callerID uint, isAdmin bool, req *UpdateListingRequest) error
	Delete(listingID int, callerID uint, isAdmin bool) error
}

type listingService struct {
	listingRepo *repo.ListingsRepo
}

func NewListingService(listingRepo *repo.ListingsRepo) ListingService {
	return &listingService{listingRepo: listingRepo}
}

func (s *listingService) Create(userID uint, req *CreateListingRequest) (*models.Listing, error) {
	listing := &models.Listing{
		UserID:        userID,
		Title:         strings.TrimSpace(req.Title),
		Description:   strings.TrimSpace(req.Description),
		Price:         req.Price,
		City:          strings.TrimSpace(req.City),
		Zip:           strings.TrimSpace(req.Zip),
		Area:          strings.TrimSpace(req.Area),
		RoomType:      models.RoomType(req.RoomType),
		Status:        models.ListingStatusActive,
		AvailableFrom: req.AvailableFrom,
		Images:        req.Images,

		ListingKind:         models.ListingType(req.ListingKind),
		SizeSqm:             req.SizeSqm,
		Deposit:             req.Deposit,
		RentalPeriod:        strings.TrimSpace(req.RentalPeriod),
		RentalPeriodDetails: strings.TrimSpace(req.RentalPeriodDetails),
		LandlordType:        strings.TrimSpace(req.LandlordType),
		FurnishedPreference: strings.TrimSpace(req.FurnishedPreference),
		Facilities:          req.Facilities,
		TargetAudience:      strings.TrimSpace(req.TargetAudience),
	}

	if err := s.listingRepo.Create(listing); err != nil {
		return nil, err
	}
	return listing, nil
}

func (s *listingService) Update(listingID int, callerID uint, isAdmin bool, req *UpdateListingRequest) error {
	listing, err := s.listingRepo.FindByID(listingID)
	if err != nil {
		return err
	}

	if !isAdmin && listing.UserID != callerID {
		return ErrNotOwner
	}

	fields := make(map[string]interface{})
	if req.Title != nil {
		fields["title"] = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		fields["description"] = strings.TrimSpace(*req.Description)
	}
	if req.Price != nil {
		fields["price"] = *req.Price
	}
	if req.City != nil {
		fields["city"] = strings.TrimSpace(*req.City)
	}
	if req.Zip != nil {
		fields["zip"] = strings.TrimSpace(*req.Zip)
	}
	if req.Area != nil {
		fields["area"] = strings.TrimSpace(*req.Area)
	}
	if req.RoomType != nil {
		fields["room_type"] = *req.RoomType
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	if req.AvailableFrom != nil {
		fields["available_from"] = *req.AvailableFrom
	}
	if req.Images != nil {
		fields["images"] = models.StringSlice(req.Images)
	}
	if req.PromotedUntil != nil && isAdmin {
		fields["promoted_until"] = req.PromotedUntil
	}
	if req.ListingKind != nil {
		fields["listing_kind"] = *req.ListingKind
	}
	if req.SizeSqm != nil {
		fields["size_sqm"] = req.SizeSqm
	}
	if req.Deposit != nil {
		fields["deposit"] = req.Deposit
	}
	if req.RentalPeriod != nil {
		fields["rental_period"] = strings.TrimSpace(*req.RentalPeriod)
	}
	if req.RentalPeriodDetails != nil {
		fields["rental_period_details"] = strings.TrimSpace(*req.RentalPeriodDetails)
	}
	if req.LandlordType != nil {
		fields["landlord_type"] = strings.TrimSpace(*req.LandlordType)
	}
	if req.FurnishedPreference != nil {
		fields["furnished_preference"] = strings.TrimSpace(*req.FurnishedPreference)
	}
	if req.Facilities != nil {
		fields["facilities"] = models.StringSlice(req.Facilities)
	}
	if req.TargetAudience != nil {
		fields["target_audience"] = strings.TrimSpace(*req.TargetAudience)
	}
	if len(fields) == 0 {
		return nil
	}

	return s.listingRepo.UpdateFields(listingID, fields)
}

func (s *listingService) Delete(listingID int, callerID uint, isAdmin bool) error {
	listing, err := s.listingRepo.FindByID(listingID)
	if err != nil {
		return err
	}

	if !isAdmin && listing.UserID != callerID {
		return ErrNotOwner
	}

	return s.listingRepo.Delete(listing)
}
