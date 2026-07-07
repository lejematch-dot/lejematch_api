package reports

import (
	"Lejematch/internal/database/repo"
	"Lejematch/internal/services"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type createReportBody struct {
	TargetType string `json:"TargetType"`
	TargetID   uint   `json:"TargetID"`
	Reason     string `json:"Reason"`
	Message    string `json:"Message"`
}

func CreateReport(c *fiber.Ctx) error {
	caller := c.Locals("user").(*services.JWTPayload)

	var body createReportBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}
	if body.TargetID == 0 || body.Reason == "" {
		return fiber.ErrBadRequest
	}

	reportService := services.NewReportService(repo.NewReportsRepo(), repo.NewListingsRepo(), repo.NewSeekersRepo())
	report, err := reportService.Create(caller.UserID, &services.CreateReportRequest{
		TargetType: body.TargetType,
		TargetID:   body.TargetID,
		Reason:     body.Reason,
		Message:    body.Message,
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidTargetType) {
			return fiber.ErrBadRequest
		}
		if errors.Is(err, services.ErrCannotReportSelf) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Du kan ikke rapportere dit eget opslag"})
		}
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        report.ID,
		"createdAt": report.CreatedAt,
	})
}
