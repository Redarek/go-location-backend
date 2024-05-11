package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"location-backend/internal/db"
	"location-backend/internal/location"
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

	if floor.WidthInPixels == nil || floor.HeightInPixels == nil {
		log.Error().Msg("Width or height of floor is nil")
		return nil
	}

	matrixInputData := location.InputData{
		Client: location.Client{
			TrSignalPower: 1,
			TrAntGain:     1,
			ZM:            0,
		},
		Walls:          s.convertWallsFromDB(walls),
		Sensors:        sensors,
		CellSizeMeters: 0.01,
		MinX:           0,
		MinY:           0,
		MaxX:           *floor.WidthInPixels,
		MaxY:           *floor.HeightInPixels,
	}
	log.Debug().Msgf("Matrix input data: %+v", matrixInputData)

	pointRows, matrixRows := location.CreateMatrix(floor.ID, matrixInputData)
	log.Debug().Msgf("Point rows: %+v", pointRows)
	log.Debug().Msgf("Matrix rows: %+v", matrixRows)

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"pointRows":  pointRows,
			"matrixRows": matrixRows,
		},
	})
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

func (s *Fiber) convertWallsFromDB(walls []*db.WallDetailed) []location.Wall {
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
