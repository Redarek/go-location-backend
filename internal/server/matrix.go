package server

import (
	"image/color"
	"location-backend/internal/db/models"
	"location-backend/internal/location"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// CreateMatrix creates a matrix
func (s *Fiber) CreateMatrix(c *fiber.Ctx) (err error) {
	floorID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse floor uuid")
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	floor, err := s.db.GetFloor(floorID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get floor")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	walls, err := s.db.GetWallsDetailed(floor.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get walls detailed")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	sensors, err := s.db.GetSensors(floor.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get sensors")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if *floor.WidthInPixels == 0 || *floor.HeightInPixels == 0 {
		log.Error().Msg("Width or height of floor is 0")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	matrixInputData := location.InputData{
		Client: location.Client{
			TrSignalPower: 17,
			TrAntGain:     1,
			ZM:            1,
		},
		Walls:          s.convertWallsFromDB(walls),
		Sensors:        sensors,
		CellSizeMeters: 0.25, // TODO fix
		MinX:           0,
		MinY:           0,
		MaxX:           int(float64((float64(*floor.WidthInPixels)**floor.Scale)/1000) / 0.25), // !be careful here
		MaxY:           int(float64((float64(*floor.WidthInPixels)**floor.Scale)/1000) / 0.25), // !be careful here
	}
	log.Debug().Msgf("Matrix input data: %+v", matrixInputData)

	pointRows, matrixRows := location.CreateMatrix(floor.ID, matrixInputData)
	//log.Debug().Msgf("Point rows: %+v", pointRows)
	//log.Debug().Msgf("Matrix rows: %+v", matrixRows)

	//responseData := fiber.Map{
	//	"data": fiber.Map{
	//		"pointRows":  pointRows,
	//		"matrixRows": matrixRows,
	//	},
	//}

	const squareSize = 1 // размер квадрата в пикселях

	dc := gg.NewContext(*floor.WidthInPixels, *floor.HeightInPixels)
	for _, point := range pointRows {
		var rssi float64 = -100
		for _, matrix := range matrixRows {
			if matrix.PointID == point.ID {
				rssi = matrix.RSSI24
				break
			}
		}

		if rssi != -100 {
			normalizedValue := normalize(rssi, -100, -25)
			clr := generateColorAndOpacity(normalizedValue)

			pointX := point.X * *floor.Scale / 1000
			pointY := point.Y * *floor.Scale / 1000

			dc.DrawRectangle(pointX, pointY, squareSize, squareSize)
			dc.SetColor(clr)
			dc.Fill()
		}
	}

	// Удаление предыдущей тепловой карты
	if floor.Heatmap != nil {
		path := filepath.Join("static", *floor.Heatmap)
		err = os.Remove(path)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete previous heatmap")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		log.Debug().Msgf("Previous heatmap deleted successfully")

	}

	fileName := uuid.New().String() + ".png"
	outputPath := filepath.Join("static", fileName)

	if _, err = os.Stat("static"); os.IsNotExist(err) {
		if err = os.Mkdir("static", os.ModePerm); err != nil {
			log.Error().Err(err).Msg("Failed to create directory")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	err = dc.SavePNG(outputPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save heatmap")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err = s.db.UpdateFloorHeatmap(floor.ID, fileName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get sensors")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	log.Debug().Msgf("Heatmap saved as %v", outputPath)

	return c.SendStatus(fiber.StatusOK)
}

// GetMatrix retrieves a matrix
func (s *Fiber) GetMatrix(c *fiber.Ctx) (err error) {
	sID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse sensor uuid")
		return
	}
	sensor, err := s.db.GetSensor(sID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get sensor")
		return
	}
	return c.JSON(fiber.Map{
		"data": sensor,
	})
}

func (s *Fiber) convertWallsFromDB(walls []*models.WallDetailed) []location.Wall {
	convertedWalls := make([]location.Wall, 0, len(walls)) // Initialize a slice to store the converted walls
	for _, dbw := range walls {
		// Ensure all required pointer fields are not nil before dereferencing to prevent runtime panics
		if dbw.X1 == nil || dbw.Y1 == nil || dbw.X2 == nil || dbw.Y2 == nil || dbw.WallType.Thickness == nil || dbw.WallType.Attenuation24 == nil || dbw.WallType.Attenuation5 == nil || dbw.WallType.Attenuation6 == nil {
			log.Error().Msg("One of the required wall coordinates or wall type info is nil")
			return nil
		}
		w := location.Wall{
			ID:            dbw.ID,
			X1:            *dbw.X1,
			Y1:            *dbw.Y1,
			X2:            *dbw.X2,
			Y2:            *dbw.Y2,
			Thickness:     *dbw.WallType.Thickness,
			Attenuation24: *dbw.WallType.Attenuation24,
			Attenuation5:  *dbw.WallType.Attenuation5,
			Attenuation6:  *dbw.WallType.Attenuation6,
		}
		convertedWalls = append(convertedWalls, w)
	}
	return convertedWalls
}

func normalize(value, min, max float64) float64 {
	return (value - min) / (max - min)
}

func generateColorAndOpacity(normalizedValue float64) color.Color {
	var r, g, b, a uint8
	a = 153 // 0.6 * 255

	if normalizedValue == 0 {
		r, g, b, a = 0, 0, 0, 0
	} else if normalizedValue <= 0.33 {
		r = 0
		g = uint8((255 * normalizedValue) / 0.33)
		b = 255
	} else if normalizedValue <= 0.66 {
		r = uint8((255 * (normalizedValue - 0.33)) / 0.33)
		g = 255
		b = 0
	} else {
		r = 255
		g = uint8((255 * (1 - normalizedValue)) / 0.34)
		b = 0
	}

	return color.NRGBA{R: r, G: g, B: b, A: a}
}
