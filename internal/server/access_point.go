package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"location-backend/internal/db"
)

// CreateAccessPoint creates an access point
func (s *Fiber) CreateAccessPoint(c *fiber.Ctx) (err error) {
	r := new(db.AccessPoint)
	err = c.BodyParser(r)
	if err != nil {
		return err
	}

	accessPointID, err := s.db.CreateAccessPoint(r)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"id": accessPointID,
	})
}

// GetAccessPoint retrieves an access point
func (s *Fiber) GetAccessPoint(c *fiber.Ctx) (err error) {
	accessPointID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse access point uuid")
		return
	}
	ap, err := s.db.GetAccessPoint(accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get access point")
		return
	}
	return c.JSON(fiber.Map{
		"data": ap,
	})
}

// GetAccessPoints retrieves access points
func (s *Fiber) GetAccessPoints(c *fiber.Ctx) (err error) {
	floorUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse floor uuid")
		return
	}
	aps, err := s.db.GetAccessPoints(floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get access points")
		return
	}
	return c.JSON(fiber.Map{
		"data": aps,
	})
}

// GetAccessPointsDetailed retrieves detailed access points
func (s *Fiber) GetAccessPointsDetailed(c *fiber.Ctx) (err error) {
	floorUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse floor uuid")
		return
	}
	aps, err := s.db.GetAccessPointsDetailed(floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get access points")
		return
	}
	return c.JSON(fiber.Map{
		"data": aps,
	})
}

// SoftDeleteAccessPoint soft delete an access point
func (s *Fiber) SoftDeleteAccessPoint(c *fiber.Ctx) (err error) {
	accessPointID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse access point uuid")
		return
	}
	isDeleted, err := s.db.IsAccessPointSoftDeleted(accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted access point")
		return
	}
	if !isDeleted {
		err = s.db.SoftDeleteAccessPoint(accessPointID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to soft delete an access point")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Access point has already been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// RestoreAccessPoint restore an access point
func (s *Fiber) RestoreAccessPoint(c *fiber.Ctx) (err error) {
	accessPointID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse access point uuid")
		return
	}
	isDeleted, err := s.db.IsAccessPointSoftDeleted(accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted access point")
		return
	}
	if isDeleted {
		err = s.db.RestoreAccessPoint(accessPointID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to restore an access point")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Access point has not been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// PatchUpdateAccessPoint patch updates an access point based on provided fields
func (s *Fiber) PatchUpdateAccessPoint(c *fiber.Ctx) error {
	var input db.AccessPoint
	if err := c.BodyParser(&input); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	if err := s.db.PatchUpdateAccessPoint(&input); err != nil {
		log.Error().Err(err).Msg("Failed to update access point")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update access point")
	}

	return c.SendStatus(fiber.StatusOK)
}

// SetRadioState upsert a radio state (on/off)
func (s *Fiber) SetRadioState(c *fiber.Ctx) (err error) {
	rs := new(db.RadioState)
	err = c.BodyParser(rs)
	if err != nil {
		return err
	}

	rsID, err := s.db.SetRadioState(rs)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"id": rsID,
	})
}
