package uploads

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const maxUploadSize = 5 * 1024 * 1024 // 5MB

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

func CreateUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ingen fil modtaget"})
	}

	if file.Size > maxUploadSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Filen er for stor (maks 5MB)"})
	}

	ext := filepath.Ext(file.Filename)
	if !allowedExtensions[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Filtypen understøttes ikke"})
	}

	if err := os.MkdirAll("./uploads", 0755); err != nil {
		return fiber.ErrInternalServerError
	}

	filename := uuid.NewString() + ext
	if err := c.SaveFile(file, filepath.Join("./uploads", filename)); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"url": "/uploads/" + filename,
	})
}
