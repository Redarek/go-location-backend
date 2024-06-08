package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"location-backend/internal/db"
)

// CreateSensor creates a sensor
func (s *Fiber) CreateSensor(c *fiber.Ctx) (err error) {
	sensor := new(db.Sensor)
	err = c.BodyParser(sensor)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return err
	}
	sID, err := s.db.CreateSensor(sensor)

	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"id": sID,
	})
}

// GetSensor retrieves a sensor
func (s *Fiber) GetSensor(c *fiber.Ctx) (err error) {
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

// GetSensorDetailed retrieves a sensor detailed
func (s *Fiber) GetSensorDetailed(c *fiber.Ctx) (err error) {
	sensorID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse sensor uuid")
		return
	}
	sensor, err := s.db.GetSensorDetailed(sensorID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get sensor detailed")
		return
	}
	return c.JSON(fiber.Map{
		"data": sensor,
	})
}

// GetSensors retrieves sensors
func (s *Fiber) GetSensors(c *fiber.Ctx) (err error) {
	floorUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse floor uuid")
		return
	}
	ss, err := s.db.GetSensors(floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get sensors")
		return
	}
	return c.JSON(fiber.Map{
		"data": ss,
	})
}

// GetSensorsDetailed retrieves sensors detailed
func (s *Fiber) GetSensorsDetailed(c *fiber.Ctx) (err error) {
	floorUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse floor uuid")
		return
	}
	ss, err := s.db.GetSensorsDetailed(floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get sensors detailed")
		return
	}
	return c.JSON(fiber.Map{
		"data": ss,
	})
}

// SoftDeleteSensor soft delete a sensor
func (s *Fiber) SoftDeleteSensor(c *fiber.Ctx) (err error) {
	sID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse sensor uuid")
		return
	}
	isDeleted, err := s.db.IsSensorSoftDeleted(sID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted sensor")
		return
	}
	if !isDeleted {
		err = s.db.SoftDeleteSensor(sID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to soft delete a sensor")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Sensor has already been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// RestoreSensor restore a sensor
func (s *Fiber) RestoreSensor(c *fiber.Ctx) (err error) {
	sID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse sensor uuid")
		return
	}
	isDeleted, err := s.db.IsSensorSoftDeleted(sID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted sensor")
		return
	}
	if isDeleted {
		err = s.db.RestoreSensor(sID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to restore a sensor")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Sensor has not been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// PatchUpdateSensor patch updates a sensor based on provided fields
func (s *Fiber) PatchUpdateSensor(c *fiber.Ctx) error {
	var sensor db.Sensor
	if err := c.BodyParser(&sensor); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	if err := s.db.PatchUpdateSensor(&sensor); err != nil {
		log.Error().Err(err).Msg("Failed to update sensor")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update sensor")
	}

	return c.SendStatus(fiber.StatusOK)
}
