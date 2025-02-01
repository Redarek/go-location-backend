package v1

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/usecase"
	"location-backend/pkg/httperrors"
)

const (
	createSensorURL       = "/"
	getSensorURL          = "/"
	getSensorDetailedURL  = "/detailed"
	getSensorsURL         = "/all"
	getSensorsDetailedURL = "/all/detailed"

	patchUpdateSensorURL = "/"

	softDeleteSensorURL = "/sd"
	restoreSensorURL    = "/restore"
)

type sensorHandler struct {
	usecase *usecase.SensorUsecase
}

// Регистрирует новый handler
func NewSensorHandler(usecase *usecase.SensorUsecase) *sensorHandler {
	return &sensorHandler{
		usecase: usecase,
	}
}

// Регистрирует маршруты для user
func (h *sensorHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createSensorURL, h.CreateSensor)

	router.Get(getSensorURL, h.GetSensor)
	router.Get(getSensorDetailedURL, h.GetSensorDetailed)
	router.Get(getSensorsURL, h.GetSensors)
	router.Get(getSensorsDetailedURL, h.GetSensorsDetailed)

	router.Patch(patchUpdateSensorURL, h.PatchUpdateSensor)

	router.Patch(softDeleteSensorURL, h.SoftDeleteSensor)
	router.Patch(restoreSensorURL, h.RestoreSensor)

	return router
}

func (h *sensorHandler) CreateSensor(ctx *fiber.Ctx) error {
	var dtoObj dto.CreateSensorDTO
	err := ctx.BodyParser(&dtoObj)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse sensor request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse sensor request body",
			nil,
		))
	}

	// TODO validate

	sensorID, err := h.usecase.CreateSensor(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			log.Info().Err(err).Msg("the floor with provided 'floor_id' does not exist")
			return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"Invalid request body",
				"The floor with provided 'floor_id' does not exist",
				nil,
			))
		}

		if errors.Is(err, usecase.ErrAlreadyExists) {
			log.Info().Err(err).Msg("the sensor with provided 'mac' already exists")
			return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"Invalid request body",
				"The sensor with provided 'mac' already exists",
				nil,
			))
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the sensor")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the sensor",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": sensorID})
}

func (h *sensorHandler) GetSensor(ctx *fiber.Ctx) error {
	sensorID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	var dtoObj dto.GetSensorDTO = dto.GetSensorDTO{
		ID: sensorID,
	}

	// TODO validate

	sensor, err := h.usecase.GetSensor(context.Background(), dtoObj.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the sensor type")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the sensor type",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensor})
}

func (h *sensorHandler) GetSensorDetailed(ctx *fiber.Ctx) error {
	sensorID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// TODO реализовать передачу page и size
	var dtoObj dto.GetSensorDetailedDTO = dto.GetSensorDetailedDTO{
		ID:   sensorID,
		Page: 1,
		Size: 100,
	}

	// TODO validate

	sensorDetailed, err := h.usecase.GetSensorDetailed(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Msg("an unexpected error has occurred while trying to retrieve the sensor detailed")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the sensor detailed",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensorDetailed})
}

func (h *sensorHandler) GetSensors(c *fiber.Ctx) error {
	floorID, err := uuid.Parse(c.Query("id"))
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
	var dtoObj dto.GetSensorsDTO = dto.GetSensorsDTO{
		FloorID: floorID,
		Page:    1,
		Size:    100,
	}

	// TODO validate

	sensors, err := h.usecase.GetSensors(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the sensors")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the sensors",
			"",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensors})
}

func (h *sensorHandler) GetSensorsDetailed(c *fiber.Ctx) error {
	floorID, err := uuid.Parse(c.Query("id"))
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
	var dtoObj dto.GetSensorsDetailedDTO = dto.GetSensorsDetailedDTO{
		FloorID: floorID,
		Page:    1,
		Size:    100,
	}

	// TODO validate

	sensorsDetailed, err := h.usecase.GetSensorsDetailed(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the sensors detailed")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the sensors detailed",
			"",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensorsDetailed})
}

func (h *sensorHandler) PatchUpdateSensor(c *fiber.Ctx) error {
	var dtoObj dto.PatchUpdateSensorDTO
	err := c.BodyParser(&dtoObj)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse sensor request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse sensor request body",
			nil,
		))
	}

	// TODO validate

	err = h.usecase.PatchUpdateSensor(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The sensor was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The sensor was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to patch update the sensor")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the sensor",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *sensorHandler) SoftDeleteSensor(c *fiber.Ctx) error {
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

	// TODO validate

	err = h.usecase.SoftDeleteSensor(context.Background(), sensorID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The sensor was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The sensor has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to soft delete the sensor")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the sensor",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *sensorHandler) RestoreSensor(c *fiber.Ctx) error {
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

	// TODO validate

	err = h.usecase.RestoreSensor(context.Background(), sensorID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The sensor was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The sensor is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to restore the sensor")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the sensor",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
