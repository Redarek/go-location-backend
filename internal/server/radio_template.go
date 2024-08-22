package server

import (
	"location-backend/internal/db/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// CreateRadioTemplate creates a radio template
func (s *Fiber) CreateRadioTemplate(c *fiber.Ctx) (err error) {
	r := new(models.RadioTemplate)
	err = c.BodyParser(r)
	if err != nil {
		return err
	}

	radioTemplateID, err := s.db.CreateRadioTemplate(r)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"id": radioTemplateID,
	})
}

// GetRadioTemplate retrieves a radio template
func (s *Fiber) GetRadioTemplate(c *fiber.Ctx) (err error) {
	radioTemplateID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse radio template uuid")
		return
	}
	r, err := s.db.GetRadioTemplate(radioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get radio template")
		return
	}
	return c.JSON(fiber.Map{
		"data": r,
	})
}

// GetRadioTemplates retrieves radio templates
func (s *Fiber) GetRadioTemplates(c *fiber.Ctx) (err error) {
	accessPointTypeUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse access point type uuid")
		return
	}
	b, err := s.db.GetRadioTemplates(accessPointTypeUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get radio template")
		return
	}
	return c.JSON(fiber.Map{
		"data": b,
	})
}

// SoftDeleteRadioTemplate soft delete a radio template
func (s *Fiber) SoftDeleteRadioTemplate(c *fiber.Ctx) (err error) {
	radioTemplateID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse radio template uuid")
		return
	}
	isDeleted, err := s.db.IsRadioTemplateSoftDeleted(radioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted radio template")
		return
	}
	if !isDeleted {
		err = s.db.SoftDeleteRadioTemplate(radioTemplateID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to soft delete a radio template")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Radio template has already been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// RestoreRadioTemplate restore a radio template
func (s *Fiber) RestoreRadioTemplate(c *fiber.Ctx) (err error) {
	radioTemplateID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse radio template uuid")
		return
	}
	isDeleted, err := s.db.IsRadioTemplateSoftDeleted(radioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted radio template")
		return
	}
	if isDeleted {
		err = s.db.RestoreRadioTemplate(radioTemplateID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to restore a radio template")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Radio template has not been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// PatchUpdateRadioTemplate patch updates a radio template based on provided fields
func (s *Fiber) PatchUpdateRadioTemplate(c *fiber.Ctx) error {
	var r models.RadioTemplate
	if err := c.BodyParser(&r); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	if err := s.db.PatchUpdateRadioTemplate(&r); err != nil {
		log.Error().Err(err).Msg("Failed to update radio template")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update radio template")
	}

	return c.SendStatus(fiber.StatusOK)
}
