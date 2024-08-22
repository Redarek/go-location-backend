package server

import (
	"location-backend/internal/db/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// CreateAccessPointType creates an access point type
func (s *Fiber) CreateAccessPointType(c *fiber.Ctx) (err error) {
	apt := new(model.AccessPointType)
	err = c.BodyParser(apt)
	if err != nil {
		return err
	}

	accessPointTypeID, err := s.db.CreateAccessPointType(apt)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"id": accessPointTypeID,
	})
}

// GetAccessPointType retrieves an access point type
func (s *Fiber) GetAccessPointType(c *fiber.Ctx) (err error) {
	accessPointTypeID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse access point type uuid")
		return
	}
	apt, err := s.db.GetAccessPointType(accessPointTypeID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get access point type")
		return
	}
	return c.JSON(fiber.Map{
		"data": apt,
	})
}

// GetAccessPointTypes retrieves access point types
func (s *Fiber) GetAccessPointTypes(c *fiber.Ctx) (err error) {
	siteUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse site uuid")
		return
	}
	apt, err := s.db.GetAccessPointTypes(siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get access point types")
		return
	}
	return c.JSON(fiber.Map{
		"data": apt,
	})
}

// SoftDeleteAccessPointType soft delete an access point type
func (s *Fiber) SoftDeleteAccessPointType(c *fiber.Ctx) (err error) {
	accessPointTypeID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse access point type uuid")
		return
	}
	isDeleted, err := s.db.IsAccessPointTypeSoftDeleted(accessPointTypeID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted access point type")
		return
	}
	if !isDeleted {
		err = s.db.SoftDeleteAccessPointType(accessPointTypeID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to soft delete a access point type")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Access point type has already been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// RestoreAccessPointType restore an access point type
func (s *Fiber) RestoreAccessPointType(c *fiber.Ctx) (err error) {
	accessPointTypeID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse access point type uuid")
		return
	}
	isDeleted, err := s.db.IsAccessPointTypeSoftDeleted(accessPointTypeID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted access point type")
		return
	}
	if isDeleted {
		err = s.db.RestoreAccessPointType(accessPointTypeID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to restore an access point type")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Access point type has not been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}
