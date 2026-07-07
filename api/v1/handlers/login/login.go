package login

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/security"
	"Lejematch/internal/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	userRepo := repo.NewUsersRepo()
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	//Cast email to lowercase
	request.Email = strings.ToLower(request.Email)

	//Get user
	user, err := userRepo.GetByEmailWithPassword(request.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	//Compare hash and password
	match, err := security.VerifyPassword(request.Password, user.Password)
	if err != nil || !match {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bekræft din e-mail før du kan logge ind",
			"code":  "email_not_verified",
		})
	}

	//Generate JWT
	jwt, err := services.GenerateJWT(*user)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	//Return JWT
	return c.JSON(fiber.Map{"token": jwt, "userID": user.ID})

}
