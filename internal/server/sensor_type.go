package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"location-backend/internal/db"
)

// CreateSensorType creates a sensor type
func (s *Fiber) CreateSensorType(c *fiber.Ctx) (err error) {
	sensorType := new(db.SensorType)
	err = c.BodyParser(sensorType)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return err
	}
	sID, err := s.db.CreateSensorType(sensorType)

	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"id": sID,
	})
}

// GetSensorType retrieves a sensor type
func (s *Fiber) GetSensorType(c *fiber.Ctx) (err error) {
	sID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse sensor type uuid")
		return
	}
	sensorType, err := s.db.GetSensorType(sID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get sensor type")
		return
	}
	return c.JSON(fiber.Map{
		"data": sensorType,
	})
}

// GetSensorTypes retrieves sensor types
func (s *Fiber) GetSensorTypes(c *fiber.Ctx) (err error) {
	floorUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse floor uuid")
		return
	}
	ss, err := s.db.GetSensorTypes(floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get sensor types")
		return
	}
	return c.JSON(fiber.Map{
		"data": ss,
	})
}

// SoftDeleteSensorType soft delete a sensor type
func (s *Fiber) SoftDeleteSensorType(c *fiber.Ctx) (err error) {
	sID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse sensor type uuid")
		return
	}
	isDeleted, err := s.db.IsSensorTypeSoftDeleted(sID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted sensor type")
		return
	}
	if !isDeleted {
		err = s.db.SoftDeleteSensorType(sID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to soft delete a sensor type")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Sensor type has already been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// RestoreSensorType restore a sensor type
func (s *Fiber) RestoreSensorType(c *fiber.Ctx) (err error) {
	sID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse sensor type uuid")
		return
	}
	isDeleted, err := s.db.IsSensorTypeSoftDeleted(sID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted sensor type")
		return
	}
	if isDeleted {
		err = s.db.RestoreSensorType(sID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to restore a sensor type")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Sensor type has not been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// PatchUpdateSensorType patch updates a sensor type based on provided fields
func (s *Fiber) PatchUpdateSensorType(c *fiber.Ctx) error {
	var sensorType db.SensorType
	if err := c.BodyParser(&sensorType); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	if err := s.db.PatchUpdateSensorType(&sensorType); err != nil {
		log.Error().Err(err).Msg("Failed to update sensor type")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update sensor type")
	}

	return c.SendStatus(fiber.StatusOK)
}
