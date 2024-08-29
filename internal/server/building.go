package server

import (
	"location-backend/internal/db/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// CreateBuilding creates a building
func (s *Fiber) CreateBuilding(c *fiber.Ctx) (err error) {
	b := new(model.Building)
	err = c.BodyParser(b)
	if err != nil {
		return err
	}

	buildingID, err := s.db.CreateBuilding(b)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"id": buildingID,
	})
}

// GetBuilding retrieves a building
func (s *Fiber) GetBuilding(c *fiber.Ctx) (err error) {
	buildingUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse building uuid")
		return
	}
	b, err := s.db.GetBuilding(buildingUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get building")
		return
	}
	return c.JSON(fiber.Map{
		"data": b,
	})
}

// GetBuildings retrieves buildings
func (s *Fiber) GetBuildings(c *fiber.Ctx) (err error) {
	siteUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse site uuid")
		return
	}
	b, err := s.db.GetBuildings(siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get buildings")
		return
	}
	return c.JSON(fiber.Map{
		"data": b,
	})
}

// SoftDeleteBuilding soft delete a building
func (s *Fiber) SoftDeleteBuilding(c *fiber.Ctx) (err error) {
	buildingUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse building uuid")
		return
	}
	isDeleted, err := s.db.IsBuildingSoftDeleted(buildingUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted building")
		return
	}
	if !isDeleted {
		err = s.db.SoftDeleteBuilding(buildingUUID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to soft delete a building")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Building has already been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// RestoreBuilding restore a building
func (s *Fiber) RestoreBuilding(c *fiber.Ctx) (err error) {
	buildingUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse building uuid")
		return
	}
	isDeleted, err := s.db.IsBuildingSoftDeleted(buildingUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted building")
		return
	}
	if isDeleted {
		err = s.db.RestoreBuilding(buildingUUID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to restore a building")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Building has not been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// PatchUpdateBuilding patch updates a building based on provided fields
func (s *Fiber) PatchUpdateBuilding(c *fiber.Ctx) error {
	var input model.Building
	if err := c.BodyParser(&input); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	if err := s.db.PatchUpdateBuilding(&input); err != nil {
		log.Error().Err(err).Msg("Failed to update building")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update building")
	}

	return c.SendStatus(fiber.StatusOK)
}
