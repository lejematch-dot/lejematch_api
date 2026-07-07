package services

import (
	"Lejematch/internal/database/models"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtExpiration = 24 * time.Hour

type JWTPayload struct {
	UserID    uint
	Email     string
	IsAdmin   bool
	IsActive  bool
	IssuedAt  int64
	ExpiresAt int64
}

func GenerateJWT(user models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub":       user.ID,
		"email":     user.Email,
		"is_admin":  user.IsAdmin,
		"is_active": user.IsActive,
		"iat":       now.Unix(),
		"exp":       now.Add(jwtExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func GetPayloadFromJWT(tokenString string) (*JWTPayload, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	sub, ok1 := claims["sub"].(float64)
	email, ok2 := claims["email"].(string)
	isAdmin, ok3 := claims["is_admin"].(bool)
	isActive, ok4 := claims["is_active"].(bool)
	iat, ok5 := claims["iat"].(float64)
	exp, ok6 := claims["exp"].(float64)
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return nil, errors.New("invalid token claims")
	}

	return &JWTPayload{
		UserID:    uint(sub),
		Email:     email,
		IsAdmin:   isAdmin,
		IsActive:  isActive,
		IssuedAt:  int64(iat),
		ExpiresAt: int64(exp),
	}, nil
}

// GenerateActionToken laver et kortlivet, formålsbundet token (e-mail-bekræftelse,
// nulstil adgangskode) ved at genbruge den eksisterende JWT-infrastruktur —
// ingen separat token-tabel nødvendig.
func GenerateActionToken(userID uint, purpose string, ttl time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub":     userID,
		"purpose": purpose,
		"iat":     now.Unix(),
		"exp":     now.Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseActionToken verificerer et token lavet af GenerateActionToken og
// tjekker at det har det forventede formål.
func ParseActionToken(tokenString, expectedPurpose string) (uint, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return 0, errors.New("JWT_SECRET is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	purpose, ok1 := claims["purpose"].(string)
	sub, ok2 := claims["sub"].(float64)
	if !ok1 || !ok2 || purpose != expectedPurpose {
		return 0, errors.New("invalid token claims")
	}

	return uint(sub), nil
}
