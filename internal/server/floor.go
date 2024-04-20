package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"location-backend/internal/db"
	"mime/multipart"
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

// UpdateFloor updates floor
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
		f.ID = uuid.Nil
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

	log.Debug().Msgf("Floor info: %+v", f)

	files := form.File["image"]
	if len(files) > 0 {
		file, err := files[0].Open()
		defer func(file multipart.File) {
			err = file.Close()
			if err != nil {
				log.Error().Err(err).Msg("Failed to defer close file")
			}
		}(file)
		if err != nil {
			log.Error().Err(err).Msg("Failed to close file")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		_, err = file.Seek(0, 0)
		img, format, err := image.Decode(file)
		if err != nil {
			log.Error().Err(err).Msg("Failed to decode file")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		newFileName := fmt.Sprintf("%s.webp", uuid.NewString())
		filePath := filepath.Join("static", newFileName)

		if err := saveImage(filePath, img, format); err != nil {
			log.Error().Err(err).Msg("Failed to save file")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		f.Image = &newFileName
		log.Debug().Msgf("Изображение сохранено: %v", filePath)
	}
	// Обновление информации об этаже в БД
	// Предполагается, что здесь вы вызываете функцию обновления БД с floor в качестве аргумента
	// Не забудьте обновлять только те поля, которые были отправлены клиентом
	err = s.db.PatchUpdateFloor(f)
	if err != nil {
		log.Error().Err(err).Msg("Failed to patch update floor")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// Отправляем ответ
	return c.SendStatus(fiber.StatusOK)
}

func saveImage(filePath string, img image.Image, format string) error {
	// Проверяем, существует ли директория `static`, если нет - создаем
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		if err := os.Mkdir("static", os.ModePerm); err != nil {
			log.Error().Err(err).Msg("Failed to create directory")
		}
	}

	// Открываем файл для записи
	fileOut, err := os.Create(filePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open file")
	}
	defer func(fileOut *os.File) {
		err = fileOut.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to defer close file")
		}
	}(fileOut)

	//// Конвертация и сохранение изображения в формате webp
	//if format != "webp" {
	//	err = webp.Encode(fileOut, img, &webp.Options{Quality: 10})
	//} else {
	//	// Если изображение уже в формате webp, просто сохраняем его
	//
	//}
	_, err = fileOut.Write(img.(*image.RGBA).Pix)

	return err
}
