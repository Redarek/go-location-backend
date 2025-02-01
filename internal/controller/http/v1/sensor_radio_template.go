package v1

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/usecase"
	"location-backend/pkg/httperrors"
)

const (
	createSensorRadioTemplateURL = "/"
	getSensorRadioTemplateURL    = "/"
	getSensorRadioTemplatesURL   = "/all"

	patchUpdateSensorRadioTemplateURL = "/"

	softDeleteSensorRadioTemplateURL = "/sd"
	restoreSensorRadioTemplateURL    = "/restore"
)

type sensorRadioTemplateHandler struct {
	usecase *usecase.SensorRadioTemplateUsecase
}

// Регистрирует новый handler
func NewSensorRadioTemplateHandler(usecase *usecase.SensorRadioTemplateUsecase) *sensorRadioTemplateHandler {
	return &sensorRadioTemplateHandler{
		usecase: usecase,
	}
}

// Регистрирует маршруты для sensor radio template
func (h *sensorRadioTemplateHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createSensorRadioTemplateURL, h.CreateSensorRadioTemplate)
	router.Get(getSensorRadioTemplateURL, h.GetSensorRadioTemplate)
	router.Get(getSensorRadioTemplatesURL, h.GetSensorRadioTemplates)

	router.Patch(patchUpdateSensorRadioTemplateURL, h.PatchUpdateSensorRadioTemplate)

	router.Patch(softDeleteSensorRadioTemplateURL, h.SoftDeleteSensorRadioTemplate)
	router.Patch(restoreSensorRadioTemplateURL, h.RestoreSensorRadioTemplate)

	return router
}

func (h *sensorRadioTemplateHandler) CreateSensorRadioTemplate(ctx *fiber.Ctx) error {
	var dtoObj dto.CreateSensorRadioTemplateDTO
	err := ctx.BodyParser(&dtoObj)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse sensor radio template request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse sensor radio template request body",
			nil,
		))
	}

	// TODO validate

	sensorRadioTemplateID, err := h.usecase.CreateSensorRadioTemplate(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			log.Info().Msg("the site with provided 'sensor_type_id' does not exist")
			return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"Invalid request body",
				"The site with provided 'sensor_type_id' does not exist",
				nil,
			))
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the sensor radio template")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the sensor radio template",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": sensorRadioTemplateID})
}

func (h *sensorRadioTemplateHandler) GetSensorRadioTemplate(ctx *fiber.Ctx) error {
	sensorRadioTemplateID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	var dtoObj dto.GetSensorRadioTemplateDTO = dto.GetSensorRadioTemplateDTO{
		ID: sensorRadioTemplateID,
	}

	// TODO validate

	sensorRadioTemplate, err := h.usecase.GetSensorRadioTemplate(context.Background(), dtoObj.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the sensor radio template")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the sensor radio template",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensorRadioTemplate})
}

func (h *sensorRadioTemplateHandler) GetSensorRadioTemplates(c *fiber.Ctx) error {
	sensorTypeID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// TODO реализовать передачу page и size
	var dtoObj dto.GetSensorRadioTemplatesDTO = dto.GetSensorRadioTemplatesDTO{
		SensorTypeID: sensorTypeID,
		Page:         1,
		Size:         100,
	}

	// TODO validate

	sensorRadioTemplates, err := h.usecase.GetSensorRadioTemplates(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the sensor radio templates")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the sensor radio templates",
			"",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensorRadioTemplates})
}

func (h *sensorRadioTemplateHandler) PatchUpdateSensorRadioTemplate(c *fiber.Ctx) error {
	var dtoObj dto.PatchUpdateSensorRadioTemplateDTO
	err := c.BodyParser(&dtoObj)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse sensor radio template request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse sensor radio template request body",
			nil,
		))
	}

	// TODO validate

	err = h.usecase.PatchUpdateSensorRadioTemplate(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The sensor radio template was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The sensor radio template was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to patch update the sensor radio template")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the sensor radio template",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *sensorRadioTemplateHandler) SoftDeleteSensorRadioTemplate(c *fiber.Ctx) error {
	sensorRadioTemplateID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// TODO validate

	err = h.usecase.SoftDeleteSensorRadioTemplate(context.Background(), sensorRadioTemplateID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The sensor radio template was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The sensor radio template has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to soft delete the sensor radio template")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the sensor radio template",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *sensorRadioTemplateHandler) RestoreSensorRadioTemplate(c *fiber.Ctx) error {
	sensorRadioTemplateID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// TODO validate

	err = h.usecase.RestoreSensorRadioTemplate(context.Background(), sensorRadioTemplateID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The sensor radio template was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The sensor radio template is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to restore the sensor radio template")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the sensor radio template",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
