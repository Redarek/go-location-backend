package server

import (
	"location-backend/internal/db/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// CreateRadio creates a radio
func (s *Fiber) CreateRadio(c *fiber.Ctx) (err error) {
	r := new(model.Radio)
	err = c.BodyParser(r)
	if err != nil {
		return err
	}

	radioID, err := s.db.CreateRadio(r)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"id": radioID,
	})
}

// GetRadio retrieves a radio
func (s *Fiber) GetRadio(c *fiber.Ctx) (err error) {
	radioID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse radio uuid")
		return
	}
	r, err := s.db.GetRadio(radioID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get radio")
		return
	}
	return c.JSON(fiber.Map{
		"data": r,
	})
}

// GetRadios retrieves radios
func (s *Fiber) GetRadios(c *fiber.Ctx) (err error) {
	siteUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse site uuid")
		return
	}
	b, err := s.db.GetRadios(siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get radio")
		return
	}
	return c.JSON(fiber.Map{
		"data": b,
	})
}

// SoftDeleteRadio soft delete a radio
func (s *Fiber) SoftDeleteRadio(c *fiber.Ctx) (err error) {
	radioID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse radio uuid")
		return
	}
	isDeleted, err := s.db.IsRadioSoftDeleted(radioID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted radio")
		return
	}
	if !isDeleted {
		err = s.db.SoftDeleteRadio(radioID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to soft delete a radio")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Radio has already been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// RestoreRadio restore a radio
func (s *Fiber) RestoreRadio(c *fiber.Ctx) (err error) {
	radioID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse radio uuid")
		return
	}
	isDeleted, err := s.db.IsRadioSoftDeleted(radioID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted radio")
		return
	}
	if isDeleted {
		err = s.db.RestoreRadio(radioID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to restore a radio")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Radio has not been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// PatchUpdateRadio patch updates a radio based on provided fields
func (s *Fiber) PatchUpdateRadio(c *fiber.Ctx) error {
	var input model.Radio
	if err := c.BodyParser(&input); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	if err := s.db.PatchUpdateRadio(&input); err != nil {
		log.Error().Err(err).Msg("Failed to update radio")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update radio")
	}

	return c.SendStatus(fiber.StatusOK)
}
