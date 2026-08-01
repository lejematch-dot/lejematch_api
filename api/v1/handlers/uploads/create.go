package uploads

import (
	"bytes"
	"errors"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const maxUploadSize = 25 * 1024 * 1024 // 25MB — moderne telefonkameraer producerer ofte billeder på 5-12MB

// maxImageDimension er det største vi nogensinde viser et billede i (lightbox
// på boligopslag) — telefonfotos er ofte 3000-4000px, hvilket er unødvendigt
// stort og gør sider langsommere at indlæse uden nogen synlig kvalitetsgevinst.
const maxImageDimension = 1920

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Filen er for stor (maks 25MB)"})
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
	dstPath := filepath.Join("./uploads", filename)

	resized := false
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
		if err := resizeAndSave(file, dstPath, ext); err != nil {
			log.Printf("upload: resize failed for %q, falling back to original: %v", file.Filename, err)
		} else {
			resized = true
		}
	}
	if !resized {
		if err := c.SaveFile(file, dstPath); err != nil {
			log.Printf("upload: SaveFile error for %q: %v", file.Filename, err)
			return fiber.ErrInternalServerError
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"url": "/uploads/" + filename,
	})
}

var errResizeNotSmaller = errors.New("resized image not smaller than original")

// resizeAndSave skalerer billedet ned hvis det er større end
// maxImageDimension i bredde eller højde, og gemmer det på dstPath.
// AutoOrientation retter telefonfotos der ellers ville blive gemt på siden
// (EXIF-rotation går tabt når pixels skrives ud igen).
//
// Nogle billeder er allerede komprimeret hårdt af afsenderen og bliver ikke
// mindre af at blive genkodet — i så fald returneres errResizeNotSmaller, så
// den oprindelige fil bruges i stedet (kalderen falder tilbage til den).
func resizeAndSave(fh *multipart.FileHeader, dstPath, ext string) error {
	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	img, err := imaging.Decode(src, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	if bounds.Dx() > maxImageDimension || bounds.Dy() > maxImageDimension {
		img = imaging.Fit(img, maxImageDimension, maxImageDimension, imaging.Lanczos)
	}

	var buf bytes.Buffer
	if ext == ".png" {
		err = png.Encode(&buf, img)
	} else {
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 82})
	}
	if err != nil {
		return err
	}

	if int64(buf.Len()) >= fh.Size {
		return errResizeNotSmaller
	}

	return os.WriteFile(dstPath, buf.Bytes(), 0644)
}
