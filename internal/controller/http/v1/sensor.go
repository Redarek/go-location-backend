package v1

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	http_dto "location-backend/internal/controller/http/dto"
	"location-backend/internal/controller/http/mapper"
	domain_dto "location-backend/internal/domain/dto"
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
	usecase      *usecase.SensorUsecase
	sensorMapper *mapper.SensorMapper
	// aprtMapper *mapper.SensorRadioTemplateMapper
}

// Регистрирует новый handler
func NewSensorHandler(usecase *usecase.SensorUsecase) *sensorHandler {
	return &sensorHandler{
		usecase:      usecase,
		sensorMapper: &mapper.SensorMapper{},
		// aprtMapper: &mapper.SensorRadioTemplateMapper{},
	}
}

// Регистрирует маршруты для user
func (h *sensorHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	// Create
	router.Post(createSensorURL, h.CreateSensor)

	// Get
	router.Get(getSensorURL, h.GetSensor)
	router.Get(getSensorDetailedURL, h.GetSensorDetailed)
	router.Get(getSensorsURL, h.GetSensors)
	router.Get(getSensorsDetailedURL, h.GetSensorsDetailed)

	// Update
	router.Patch(patchUpdateSensorURL, h.PatchUpdateSensor)

	// Delete
	router.Patch(softDeleteSensorURL, h.SoftDeleteSensor)
	router.Patch(restoreSensorURL, h.RestoreSensor)

	return router
}

func (h *sensorHandler) CreateSensor(ctx *fiber.Ctx) error {
	var httpDTO http_dto.CreateSensorDTO
	err := ctx.BodyParser(&httpDTO)
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

	// Mapping http DTO -> domain DTO
	domainDTO := (*domain_dto.CreateSensorDTO)(&httpDTO)
	// domainDTO := h.aptMapper.CreateHTTPtoDomain(&httpDTO) // Избыточный метод

	sensorID, err := h.usecase.CreateSensor(context.Background(), domainDTO)
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

	// var dto http_dto.GetSensorDTO = http_dto.GetSensorDTO{
	// 	ID: sensorID,
	// }

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.GetSensorDTO{
	// 	ID: dto.ID,
	// }

	sensor, err := h.usecase.GetSensor(context.Background(), sensorID)
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

	// Mapping entity -> http DTO
	// apHttpDTO := h.aptMapper.EntityDomainToHTTP(aptDomainDTO)
	apHttpDTO := (http_dto.SensorDTO)(*sensor)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": apHttpDTO})
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

	// var dto http_dto.GetSensorDTO = http_dto.GetSensorDTO{
	// 	ID: sensorID,
	// }

	// TODO реализовать передачу page и size
	var dto http_dto.GetSensorDetailedDTO = http_dto.GetSensorDetailedDTO{
		ID:   sensorID,
		Page: 1,
		Size: 100,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.GetSensorDTO{
	// 	ID: dto.ID,
	// }

	domainDTO := domain_dto.GetSensorDetailedDTO{
		ID:     dto.ID,
		Limit:  dto.Size,
		Offset: (dto.Page - 1) * dto.Size,
	}

	sensorDetailed, err := h.usecase.GetSensorDetailed(context.Background(), domainDTO)
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

	// Mapping sensor radio entity -> http DTO
	var sensorRadioHttpDTOs []*http_dto.SensorRadioDTO
	for _, sensorRadioHttpDTO := range sensorDetailed.Radios {
		sensorRadioHttpDTOs = append(sensorRadioHttpDTOs, (*http_dto.SensorRadioDTO)(sensorRadioHttpDTO))
	}

	// Mapping entity -> http DTO
	sensorDetailedDTO := http_dto.SensorDetailedDTO{
		SensorDTO:  (http_dto.SensorDTO)(sensorDetailed.Sensor),
		SensorType: (http_dto.SensorTypeDTO)(sensorDetailed.SensorType),
		Radios:     sensorRadioHttpDTOs,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensorDetailedDTO})
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
	var dto http_dto.GetSensorsDTO = http_dto.GetSensorsDTO{
		FloorID: floorID,
		Page:    1,
		Size:    100,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.GetSensorsDTO{
		FloorID: dto.FloorID,
		Limit:   dto.Size,
		Offset:  (dto.Page - 1) * dto.Size,
	}

	apDomainDTOs, err := h.usecase.GetSensors(context.Background(), domainDTO)
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

	// Mapping entity -> http DTO
	var apHttpDTOs []http_dto.SensorDTO
	for _, apDomainDTO := range apDomainDTOs {
		sensorDTO := (http_dto.SensorDTO)(*apDomainDTO)
		apHttpDTOs = append(apHttpDTOs, sensorDTO)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": apHttpDTOs})
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
	var dto http_dto.GetSensorsDetailedDTO = http_dto.GetSensorsDetailedDTO{
		FloorID: floorID,
		Page:    1,
		Size:    100,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.GetSensorsDetailedDTO{
		FloorID: dto.FloorID,
		Limit:   dto.Size,
		Offset:  (dto.Page - 1) * dto.Size,
	}

	sensorsDetailed, err := h.usecase.GetSensorsDetailed(context.Background(), domainDTO)
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

	// Mapping entity -> http DTO
	sensorsDetailedDTO := h.sensorMapper.DetailedToHTTPList(sensorsDetailed)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensorsDetailedDTO})
}

func (h *sensorHandler) PatchUpdateSensor(c *fiber.Ctx) error {
	var httpDTO http_dto.PatchUpdateSensorDTO
	err := c.BodyParser(&httpDTO)
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

	// Mapping http DTO -> domain DTO
	domainDTO := (*domain_dto.PatchUpdateSensorDTO)(&httpDTO)

	err = h.usecase.PatchUpdateSensor(context.Background(), domainDTO)
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

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.SoftDeleteSensorDTO{
	// 	ID: sensorID,
	// }

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

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.SoftDeleteSensorDTO{
	// 	ID: sensorID,
	// }

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
