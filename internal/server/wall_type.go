package server

import (
	"location-backend/internal/db/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// CreateWallType creates a wall type
func (s *Fiber) CreateWallType(c *fiber.Ctx) (err error) {
	wt := new(model.WallType)
	err = c.BodyParser(wt)
	if err != nil {
		return err
	}

	wallTypeID, err := s.db.CreateWallType(wt)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"id": wallTypeID,
	})
}

// GetWallType retrieves a wall type
func (s *Fiber) GetWallType(c *fiber.Ctx) (err error) {
	wallTypeID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse wall type uuid")
		return
	}
	wt, err := s.db.GetWallType(wallTypeID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get wall type")
		return
	}
	return c.JSON(fiber.Map{
		"data": wt,
	})
}

// GetWallTypes retrieves wall types
func (s *Fiber) GetWallTypes(c *fiber.Ctx) (err error) {
	siteUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse site uuid")
		return
	}
	wt, err := s.db.GetWallTypes(siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get wall type")
		return
	}
	return c.JSON(fiber.Map{
		"data": wt,
	})
}

// SoftDeleteWallType soft delete a wall type
func (s *Fiber) SoftDeleteWallType(c *fiber.Ctx) (err error) {
	wallTypeID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse wall type uuid")
		return
	}
	isDeleted, err := s.db.IsWallTypeSoftDeleted(wallTypeID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted wall type")
		return
	}
	if !isDeleted {
		err = s.db.SoftDeleteWallType(wallTypeID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to soft delete a wall type")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Wall type has already been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// RestoreWallType restore a wall type
func (s *Fiber) RestoreWallType(c *fiber.Ctx) (err error) {
	wallTypeID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse wall type uuid")
		return
	}
	isDeleted, err := s.db.IsWallTypeSoftDeleted(wallTypeID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted wall type")
		return
	}
	if isDeleted {
		err = s.db.RestoreWallType(wallTypeID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to restore a wall type")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Wall type has not been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// PatchUpdateWallType patch updates a wall type based on provided fields
func (s *Fiber) PatchUpdateWallType(c *fiber.Ctx) error {
	var input model.WallType
	if err := c.BodyParser(&input); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}
	log.Debug().Msgf("Updating wall type: %v", input)
	if err := s.db.PatchUpdateWallType(&input); err != nil {
		log.Error().Err(err).Msg("Failed to update wall type")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update wall type")
	}

	return c.SendStatus(fiber.StatusOK)
}
