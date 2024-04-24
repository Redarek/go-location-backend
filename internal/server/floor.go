package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"location-backend/internal/db"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CreateFloor creates a floor
func (s *Fiber) CreateFloor(c *fiber.Ctx) (err error) {
	f := new(db.Floor)
	err = c.BodyParser(f)
	if err != nil {
		return err
	}

	floorID, err := s.db.CreateFloor(f)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"id": floorID,
	})
}

// GetFloor retrieves a floor
func (s *Fiber) GetFloor(c *fiber.Ctx) (err error) {
	floorUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse floor uuid")
		return
	}
	f, err := s.db.GetFloor(floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get floor")
		return
	}
	return c.JSON(fiber.Map{
		"data": f,
	})
}

// GetFloors retrieves floors
func (s *Fiber) GetFloors(c *fiber.Ctx) (err error) {
	buildingUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse site uuid")
		return
	}
	f, err := s.db.GetFloors(buildingUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get floors")
		return
	}
	return c.JSON(fiber.Map{
		"data": f,
	})
}

// SoftDeleteFloor soft delete a floor
func (s *Fiber) SoftDeleteFloor(c *fiber.Ctx) (err error) {
	floorUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse floor uuid")
		return
	}
	isDeleted, err := s.db.IsFloorSoftDeleted(floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted floor")
		return
	}
	if !isDeleted {
		err = s.db.SoftDeleteFloor(floorUUID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to soft delete a floor")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Floor has already been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// RestoreFloor restore a floor
func (s *Fiber) RestoreFloor(c *fiber.Ctx) (err error) {
	floorUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse floor uuid")
		return
	}
	isDeleted, err := s.db.IsFloorSoftDeleted(floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted floor")
		return
	}
	if isDeleted {
		err = s.db.RestoreFloor(floorUUID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to restore a floor")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Floor has not been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// PatchUpdateFloor updates floor
func (s *Fiber) PatchUpdateFloor(c *fiber.Ctx) error {

	form, err := c.MultipartForm()
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse multipart form")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	f := &db.Floor{}
	if id, ok := form.Value["id"]; ok && id[0] != "" {
		f.ID, err = uuid.Parse(id[0])
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse floor uuid")
			return c.Status(fiber.StatusBadRequest).SendString("Invalid floor UUID")
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("id is required")
	}
	if name, ok := form.Value["name"]; ok && name[0] != "" {
		f.Name = &name[0]
	} else {
		f.Name = nil
	}

	if number, ok := form.Value["number"]; ok && number[0] != "" {
		parsedNumber, err := strconv.Atoi(strings.TrimSpace(number[0]))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid number")
		}
		f.Number = &parsedNumber
	} else {
		f.Number = nil
	}

	if scale, ok := form.Value["scale"]; ok && scale[0] != "" {
		parsedScale, err := strconv.ParseFloat(strings.TrimSpace(scale[0]), 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid Scale")
		}
		f.Scale = &parsedScale
	} else {
		f.Scale = nil
	}

	files := form.File["image"]
	if len(files) > 0 {
		file, err := files[0].Open()
		if err != nil {
			log.Error().Err(err).Msg("Failed to open file")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer file.Close()

		fileExtension := filepath.Ext(files[0].Filename)
		newFileName := uuid.NewString() + fileExtension // Using UUID and original file extension
		filePath := filepath.Join("static", newFileName)

		img, _, err := image.Decode(file) // We no longer use format here, only for decoding
		if err != nil {
			log.Error().Err(err).Msg("Failed to decode image")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if err := saveImage(filePath, img); err != nil {
			log.Error().Err(err).Msg("Failed to save image")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		f.Image = &newFileName
		log.Debug().Msgf("Image saved: %v", filePath)
	}

	log.Debug().Msgf("Floor info: %+v", f)

	err = s.db.PatchUpdateFloor(f)
	if err != nil {
		log.Error().Err(err).Msg("Failed to patch update floor")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

// saveImage saves an image to the specified file path, handling different image formats
func saveImage(filePath string, img image.Image) error {
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		if err = os.Mkdir("static", os.ModePerm); err != nil {
			log.Error().Err(err).Msg("Failed to create directory")
			return err
		}
	}

	fileOut, err := os.Create(filePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create file")
		return err
	}
	defer fileOut.Close()

	// Determine the image format based on the file extension
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".jpeg", ".jpg":
		err = jpeg.Encode(fileOut, img, nil) // nil EncoderOptions uses default settings
	case ".png":
		err = png.Encode(fileOut, img)
	default:
		// Default to JPEG if no recognized format is specified
		err = jpeg.Encode(fileOut, img, nil)
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to encode and save image")
		return err
	}

	return nil
}
