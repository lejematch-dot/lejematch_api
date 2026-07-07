package uploads

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const maxUploadSize = 5 * 1024 * 1024 // 5MB

// allowedExtensions bruges som fallback hvis browseren ikke sender en
// genkendelig Content-Type. Filnavne fra telefoner (f.eks. iPhones
// "IMG_1234.JPG") har ofte endelsen med store bogstaver.
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

var contentTypeExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

func CreateUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("upload: FormFile error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ingen fil modtaget"})
	}
	log.Printf("upload: received file=%q size=%d contentType=%q", file.Filename, file.Size, file.Header.Get("Content-Type"))

	if file.Size > maxUploadSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Filen er for stor (maks 5MB)"})
	}

	contentType := strings.ToLower(strings.TrimSpace(file.Header.Get("Content-Type")))
	ext, ok := contentTypeExtensions[contentType]
	if !ok {
		ext = strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExtensions[ext] {
			if contentType == "image/heic" || contentType == "image/heif" || ext == ".heic" || ext == ".heif" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Billedet er i HEIC-format, som ikke understøttes. Skift til JPEG under Indstillinger → Kamera → Formater → Mest kompatibelt på din iPhone, eller vælg \"Behold som JPEG\" når du deler billedet.",
				})
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Filtypen understøttes ikke"})
		}
	}

	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Printf("upload: MkdirAll error: %v", err)
		return fiber.ErrInternalServerError
	}

	filename := uuid.NewString() + ext
	if err := c.SaveFile(file, filepath.Join("./uploads", filename)); err != nil {
		log.Printf("upload: SaveFile error for %q: %v", file.Filename, err)
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"url": "/uploads/" + filename,
	})
}
