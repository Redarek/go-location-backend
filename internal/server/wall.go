package server

import (
	"location-backend/internal/db/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// CreateWall creates a wall
func (s *Fiber) CreateWall(c *fiber.Ctx) (err error) {
	w := new(model.Wall)
	err = c.BodyParser(w)
	if err != nil {
		return err
	}

	wallID, err := s.db.CreateWall(w)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"id": wallID,
	})
}

// GetWall retrieves a wall
func (s *Fiber) GetWall(c *fiber.Ctx) (err error) {
	wallID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse wall uuid")
		return
	}
	w, err := s.db.GetWall(wallID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get wall")
		return
	}
	return c.JSON(fiber.Map{
		"data": w,
	})
}

// GetWalls retrieves walls
func (s *Fiber) GetWalls(c *fiber.Ctx) (err error) {
	siteUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse site uuid")
		return
	}
	w, err := s.db.GetWalls(siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get walls")
		return
	}
	return c.JSON(fiber.Map{
		"data": w,
	})
}

// SoftDeleteWall soft delete a wall
func (s *Fiber) SoftDeleteWall(c *fiber.Ctx) (err error) {
	wallID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse wall uuid")
		return
	}
	isDeleted, err := s.db.IsWallSoftDeleted(wallID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted wall")
		return
	}
	if !isDeleted {
		err = s.db.SoftDeleteWall(wallID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to soft delete a wall")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Wall has already been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// RestoreWall restore a wall
func (s *Fiber) RestoreWall(c *fiber.Ctx) (err error) {
	wallID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse wall uuid")
		return
	}
	isDeleted, err := s.db.IsWallSoftDeleted(wallID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted wall")
		return
	}
	if isDeleted {
		err = s.db.RestoreWall(wallID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to restore a wall")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Wall has not been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// PatchUpdateWall patch updates a wall based on provided fields
func (s *Fiber) PatchUpdateWall(c *fiber.Ctx) error {
	var input model.Wall
	if err := c.BodyParser(&input); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	if err := s.db.PatchUpdateWall(&input); err != nil {
		log.Error().Err(err).Msg("Failed to update wall")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update wall")
	}

	return c.SendStatus(fiber.StatusOK)
}
