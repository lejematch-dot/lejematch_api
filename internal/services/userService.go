package services

import (
	"Lejematch/internal/database/models"
	"Lejematch/internal/database/repo"
	"Lejematch/internal/security"
	"errors"
	"strings"
	"unicode"
)

// isDuplicateKeyError detects PostgreSQL unique constraint violations (code 23505).
func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "23505")
}

var (
	ErrInvalidCredentials = errors.New("invalid current password")
	ErrDuplicateEntry     = errors.New("email or phone already in use")
)

type CreateUserRequest struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Password  string
	City      string
	ImageURL  string
}

type UpdatePasswordRequest struct {
	CurrentPassword string
	NewPassword     string
}

type UpdateUserRequest struct {
	FirstName *string
	LastName  *string
	Email     *string
	Phone     *string
}

type UpdateProfileRequest struct {
	DisplayName *string
	Bio         *string
	City        *string
	ImageURL    *string
}

type UserService interface {
	CreateUserWithProfile(req *CreateUserRequest) (*models.User, error)
	UpdatePassword(userID int, req *UpdatePasswordRequest) error
	UpdateUser(userID int, req *UpdateUserRequest) error
	UpdateProfile(userID int, req *UpdateProfileRequest) error
}

type userService struct {
	userRepo    *repo.UsersRepo
	profileRepo *repo.ProfilesRepo
}

func NewUserService(userRepo *repo.UsersRepo, profileRepo *repo.ProfilesRepo) UserService {
	return &userService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
	}
}

// titleCase capitalizes the first letter of each word and lowercases the rest
func titleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			for j := 1; j < len(runes); j++ {
				runes[j] = unicode.ToLower(runes[j])
			}
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

func (s *userService) CreateUserWithProfile(req *CreateUserRequest) (*models.User, error) {
	// Clean input
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.City = strings.TrimSpace(req.City)
	req.ImageURL = strings.TrimSpace(req.ImageURL)

	// Hash password
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashedPassword,
		IsAdmin:   false,
		IsActive:  true,
	}

	// Save user
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	displayName := titleCase(req.FirstName + " " + req.LastName)

	// Create profile linked to user
	profile := &models.Profile{
		UserID:      user.ID,
		DisplayName: displayName,
		City:        req.City,
		ImageURL:    req.ImageURL,
		Phone:       req.Phone,
		Email:       req.Email,
	}

	// Save profile
	err = s.profileRepo.Create(profile)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(userID int, req *UpdateUserRequest) error {
	fields := make(map[string]interface{})
	if req.FirstName != nil {
		fields["first_name"] = strings.TrimSpace(*req.FirstName)
	}
	if req.LastName != nil {
		fields["last_name"] = strings.TrimSpace(*req.LastName)
	}
	if req.Email != nil {
		fields["email"] = strings.TrimSpace(*req.Email)
	}
	if req.Phone != nil {
		fields["phone"] = strings.TrimSpace(*req.Phone)
	}
	if len(fields) == 0 {
		return nil
	}

	if err := s.userRepo.UpdateFields(userID, fields); err != nil {
		if isDuplicateKeyError(err) {
			return ErrDuplicateEntry
		}
		return err
	}

	// Keep profile email/phone in sync.
	profileFields := make(map[string]interface{})
	if req.Email != nil {
		profileFields["email"] = fields["email"]
	}
	if req.Phone != nil {
		profileFields["phone"] = fields["phone"]
	}
	if len(profileFields) > 0 {
		if err := s.profileRepo.UpdateByUserID(uint(userID), profileFields); err != nil {
			return err
		}
	}

	return nil
}

func (s *userService) UpdateProfile(userID int, req *UpdateProfileRequest) error {
	fields := make(map[string]interface{})
	if req.DisplayName != nil {
		fields["display_name"] = strings.TrimSpace(*req.DisplayName)
	}
	if req.Bio != nil {
		fields["bio"] = strings.TrimSpace(*req.Bio)
	}
	if req.City != nil {
		fields["city"] = strings.TrimSpace(*req.City)
	}
	if req.ImageURL != nil {
		fields["image_url"] = strings.TrimSpace(*req.ImageURL)
	}
	if len(fields) == 0 {
		return nil
	}

	return s.profileRepo.UpdateByUserID(uint(userID), fields)
}

func (s *userService) UpdatePassword(userID int, req *UpdatePasswordRequest) error {
	user, err := s.userRepo.FindByIDWithPassword(userID)
	if err != nil {
		return err
	}

	ok, err := security.VerifyPassword(req.CurrentPassword, user.Password)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidCredentials
	}

	hashed, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashed
	return s.userRepo.Update(user)
}
