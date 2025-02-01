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
	createSensorTypeURL      = "/"
	getSensorTypeURL         = "/"
	getSensorTypeDetailedURL = "/detailed"
	getSensorTypesURL        = "/all"

	patchUpdateSensorTypeURL = "/"

	softDeleteSensorTypeURL = "/sd"
	restoreSensorTypeURL    = "/restore"
)

type sensorTypeHandler struct {
	usecase *usecase.SensorTypeUsecase
}

// Регистрирует новый handler
func NewSensorTypeHandler(usecase *usecase.SensorTypeUsecase) *sensorTypeHandler {
	return &sensorTypeHandler{
		usecase: usecase,
	}
}

// Регистрирует маршруты для user
func (h *sensorTypeHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createSensorTypeURL, h.CreateSensorType)
	router.Get(getSensorTypeURL, h.GetSensorType)
	// router.Get(getSensorTypeDetailedURL, h.GetSensorTypeDetailed)
	router.Get(getSensorTypesURL, h.GetSensorTypes)

	router.Patch(patchUpdateSensorTypeURL, h.PatchUpdateSensorType)

	router.Patch(softDeleteSensorTypeURL, h.SoftDeleteSensorType)
	router.Patch(restoreSensorTypeURL, h.RestoreSensorType)

	// TODO Get list detailed
	return router
}

func (h *sensorTypeHandler) CreateSensorType(ctx *fiber.Ctx) error {
	var dtoObj dto.CreateSensorTypeDTO
	err := ctx.BodyParser(&dtoObj)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse sensor type request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse sensor type request body",
			nil,
		))
	}

	// TODO validate

	sensorTypeID, err := h.usecase.CreateSensorType(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			log.Info().Err(err).Msg("the site with provided 'site_id' does not exist")
			return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"Invalid request body",
				"The site with provided 'site_id' does not exist",
				nil,
			))
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the sensor type")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the sensor type",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": sensorTypeID})
}

func (h *sensorTypeHandler) GetSensorType(ctx *fiber.Ctx) error {
	sensorTypeID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	var dtoObj dto.GetSensorTypeDTO = dto.GetSensorTypeDTO{
		ID: sensorTypeID,
	}

	// TODO validate

	sensorType, err := h.usecase.GetSensorType(context.Background(), dtoObj.ID)
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

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensorType})
}

// func (h *sensorTypeHandler) GetSensorTypeDetailed(ctx *fiber.Ctx) error {
// 	sensorTypeID, err := uuid.Parse(ctx.Query("id"))
// 	if err != nil {
// 		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
// 		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
// 			fiber.StatusBadRequest,
// 			"Invalid ID",
// 			"Failed to parse 'id' as UUID",
// 			nil,
// 		))
// 	}

// 	// var dto http_dto.GetSensorTypeDTO = http_dto.GetSensorTypeDTO{
// 	// 	ID: sensorTypeID,
// 	// }

// 	// TODO реализовать передачу page и size
// 	var dto http_dto.GetSensorTypeDetailedDTO = http_dto.GetSensorTypeDetailedDTO{
// 		SensorTypeID: sensorTypeID,
// 		Page:              1,
// 		Size:              100,
// 	}

// 	// TODO validate

// 	// Mapping http DTO -> domain DTO
// 	// domainDTO := domain_dto.GetSensorTypeDTO{
// 	// 	ID: dto.ID,
// 	// }

// 	domainDTO := domain_dto.GetSensorTypeDetailedDTO{
// 		ID:     dto.SensorTypeID,
// 		Limit:  dto.Size,
// 		Offset: (dto.Page - 1) * dto.Size,
// 	}

// 	sensorTypeDetailed, err := h.usecase.GetSensorTypeDetailed(context.Background(), domainDTO)
// 	if err != nil {
// 		if errors.Is(err, usecase.ErrNotFound) {
// 			ctx.Status(fiber.StatusNoContent)
// 			return nil
// 		}

// 		log.Error().Msg("an unexpected error has occurred while trying to retrieve the sensor type detailed")
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
// 			fiber.StatusInternalServerError,
// 			"An unexpected error has occurred while trying to retrieve the sensor type detailed",
// 			"",
// 			nil,
// 		))
// 	}

// 	// Mapping access point radio template entity -> http DTO
// 	var aprtHttpDTOs []*http_dto.AccessPointRadioTemplateDTO
// 	for _, aprtHttpDTO := range sensorTypeDetailed.RadioTemplates {
// 		aprtHttpDTOs = append(aprtHttpDTOs, (*http_dto.AccessPointRadioTemplateDTO)(aprtHttpDTO))
// 	}

// 	// Mapping entity -> http DTO
// 	sensorTypeDetailedDTO := http_dto.SensorTypeDetailedDTO{
// 		SensorTypeDTO: (http_dto.SensorTypeDTO)(sensorTypeDetailed.SensorType),
// 		RadioTemplatesDTO:  aprtHttpDTOs,
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensorTypeDetailedDTO})
// }

func (h *sensorTypeHandler) GetSensorTypes(c *fiber.Ctx) error {
	siteID, err := uuid.Parse(c.Query("id"))
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
	var dtoObj dto.GetSensorTypesDTO = dto.GetSensorTypesDTO{
		SiteID: siteID,
		Page:   1,
		Size:   100,
	}

	// TODO validate

	sensorTypes, err := h.usecase.GetSensorTypes(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the sensorType")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the sensor type",
			"",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": sensorTypes})
}

func (h *sensorTypeHandler) PatchUpdateSensorType(c *fiber.Ctx) error {
	var dtoObj dto.PatchUpdateSensorTypeDTO
	err := c.BodyParser(&dtoObj)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse sensor type request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse sensor type request body",
			nil,
		))
	}

	// TODO validate

	err = h.usecase.PatchUpdateSensorType(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The sensor type was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The sensor type was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to patch update the sensor type")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the sensor type",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *sensorTypeHandler) SoftDeleteSensorType(c *fiber.Ctx) error {
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

	// TODO validate

	err = h.usecase.SoftDeleteSensorType(context.Background(), sensorTypeID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The sensorType was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The sensorType has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to soft delete the sensor type")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the sensor type",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *sensorTypeHandler) RestoreSensorType(c *fiber.Ctx) error {
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

	// TODO validate

	err = h.usecase.RestoreSensorType(context.Background(), sensorTypeID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The sensor type was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The sensor type is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to restore the sensor type")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the sensor type",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
