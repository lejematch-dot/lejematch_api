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