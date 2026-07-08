package services

import (
	"Lejematch/config"
	"Lejematch/internal/database/models"
	"Lejematch/internal/database/repo"
	"Lejematch/internal/security"
	"errors"
	"log"
	"strings"
	"time"
	"unicode"
)

const emailVerificationTTL = 1 * time.Hour

// isDuplicateKeyError detects PostgreSQL unique constraint violations (code 23505).
func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "23505")
}

var (
	ErrInvalidCredentials = errors.New("invalid current password")
	ErrDuplicateEntry     = errors.New("email or phone already in use")
	ErrImageRequired      = errors.New("profile image is required")
)

type CreateUserRequest struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Password  string
	City      string
	ImageURL  string

	Age      *int
	UserType string
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
	Age         *int
	UserType    *string
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
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Phone = strings.TrimSpace(req.Phone)
	req.City = strings.TrimSpace(req.City)
	req.ImageURL = strings.TrimSpace(req.ImageURL)

	if req.ImageURL == "" {
		return nil, ErrImageRequired
	}

	// Hash password
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user — inaktiv indtil e-mailen er bekræftet
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashedPassword,
		IsAdmin:   false,
		IsActive:  false,
	}

	// Save user
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	displayName := titleCase(req.FirstName + " " + req.LastName)

	userType := strings.TrimSpace(req.UserType)
	if userType == "" {
		userType = "tenant"
	}

	// Create profile linked to user
	profile := &models.Profile{
		UserID:      user.ID,
		DisplayName: displayName,
		City:        req.City,
		ImageURL:    req.ImageURL,
		Phone:       req.Phone,
		Email:       req.Email,
		Age:         req.Age,
		UserType:    userType,
	}

	// Save profile
	err = s.profileRepo.Create(profile)
	if err != nil {
		return nil, err
	}

	// Send bekræftelses-mail. Fejl her må ikke fejle selve oprettelsen —
	// brugeren kan altid bede om en ny via "resend-verification" — men skal
	// stadig logges, ellers er en fejlkonfigureret mailer usynlig.
	if err := sendVerificationEmail(user.ID, user.Email, user.FirstName); err != nil {
		log.Printf("failed to send verification email to %s: %v", user.Email, err)
	}

	return user, nil
}

// sendVerificationEmail sender et bekræftelseslink til den nyoprettede bruger.
func sendVerificationEmail(userID uint, email, firstName string) error {
	token, err := GenerateActionToken(userID, "verify_email", emailVerificationTTL)
	if err != nil {
		return err
	}

	link := config.AppConfigInstance.FrontendURL + "/bekraeft-email/" + token
	subject := "Bekræft din e-mail hos LejeMatch"
	html := `
	<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<h2>Velkommen til LejeMatch, ` + firstName + `!</h2>
			<p>Klik på linket for at bekræfte din e-mail og aktivere din konto:</p>
			<p><a href="` + link + `">Bekræft din e-mail</a></p>
			<p style="color: #666; font-size: 12px;">Linket udløber om 1 time.</p>
		</body>
	</html>
	`
	return SendEmail(email, subject, html)
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
		fields["email"] = strings.ToLower(strings.TrimSpace(*req.Email))
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
	if req.Age != nil {
		fields["age"] = req.Age
	}
	if req.UserType != nil {
		fields["user_type"] = strings.TrimSpace(*req.UserType)
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
