package auth

import (
	"Lejematch/internal/services"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func JWTmiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.ErrUnauthorized
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return fiber.ErrUnauthorized
	}
	tokenString := strings.TrimSpace(parts[1])

	payload, err := services.GetPayloadFromJWT(tokenString)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	if time.Now().Unix() > payload.ExpiresAt {
		return fiber.ErrUnauthorized
	}

	c.Locals("user", payload)
	return c.Next()
}
