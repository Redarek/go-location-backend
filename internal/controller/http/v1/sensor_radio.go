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
	createSensorRadioURL = "/"
	getSensorRadioURL    = "/"
	getSensorRadiosURL   = "/all"

	patchUpdateSensorRadioURL = "/"

	softDeleteSensorRadioURL = "/sd"
	restoreSensorRadioURL    = "/restore"
)

type sensorRadioHandler struct {
	usecase *usecase.SensorRadioUsecase
}

// Регистрирует новый handler
func NewSensorRadioHandler(usecase *usecase.SensorRadioUsecase) *sensorRadioHandler {
	return &sensorRadioHandler{
		usecase: usecase,
	}
}

// Регистрирует маршруты для user
func (h *sensorRadioHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createSensorRadioURL, h.CreateSensorRadio)
	router.Get(getSensorRadioURL, h.GetSensorRadio)
	router.Get(getSensorRadiosURL, h.GetSensorRadios)

	router.Patch(patchUpdateSensorRadioURL, h.PatchUpdateSensorRadio)

	router.Patch(softDeleteSensorRadioURL, h.SoftDeleteSensorRadio)
	router.Patch(restoreSensorRadioURL, h.RestoreSensorRadio)

	return router
}

func (h *sensorRadioHandler) CreateSensorRadio(ctx *fiber.Ctx) error {
	var dtoObj dto.CreateSensorRadioDTO
	err := ctx.BodyParser(&dtoObj)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse sensor radio request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse sensor radio request body",
			nil,
		))
	}

	// TODO validate

	sensorRadioID, err := h.usecase.CreateSensorRadio(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			log.Info().Msg("the access point with provided 'sensor_id' does not exist")
			return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"Invalid request body",
				"The access point with provided 'sensor_id' does not exist",
				nil,
			))
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the sensor radio")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the sensor radio",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": sensorRadioID})
}

func (h *sensorRadioHandler) GetSensorRadio(ctx *fiber.Ctx) error {
	sensorRadioID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	var dtoObj dto.GetSensorRadioDTO = dto.GetSensorRadioDTO{
		ID: sensorRadioID,
	}

	// TODO validate

	sensorRadio, err := h.usecase.GetSensorRadio(context.Background(), dtoObj.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the sensor radio")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the sensor radio",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensorRadio})
}

func (h *sensorRadioHandler) GetSensorRadios(c *fiber.Ctx) error {
	sensorID, err := uuid.Parse(c.Query("id"))
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
	var dtoObj dto.GetSensorRadiosDTO = dto.GetSensorRadiosDTO{
		SensorID: sensorID,
		Page:     1,
		Size:     100,
	}

	// TODO validate

	sensorRadios, err := h.usecase.GetSensorRadios(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the sensor radios")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the sensor radios",
			"",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensorRadios})
}

func (h *sensorRadioHandler) PatchUpdateSensorRadio(c *fiber.Ctx) error {
	var dtoObj dto.PatchUpdateSensorRadioDTO
	err := c.BodyParser(&dtoObj)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse sensor radio request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse sensor radio request body",
			nil,
		))
	}

	// TODO validate

	err = h.usecase.PatchUpdateSensorRadio(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The sensor radio was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The sensor radio was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to patch update the sensor radio")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the sensor radio",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *sensorRadioHandler) SoftDeleteSensorRadio(c *fiber.Ctx) error {
	sensorRadioID, err := uuid.Parse(c.Query("id"))
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

	err = h.usecase.SoftDeleteSensorRadio(context.Background(), sensorRadioID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The sensor radio was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The sensor radio has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to soft delete the sensor radio")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the sensor radio",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *sensorRadioHandler) RestoreSensorRadio(c *fiber.Ctx) error {
	sensorRadioID, err := uuid.Parse(c.Query("id"))
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

	err = h.usecase.RestoreSensorRadio(context.Background(), sensorRadioID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The sensor radio was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The sensor radio is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to restore the sensor radio")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the sensor radio",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
