package v1

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/usecase"
	"location-backend/pkg/httperrors"
)

const (
	createAccessPointURL      = "/"
	getAccessPointURL         = "/"
	getAccessPointDetailedURL = "/detailed"
	getAccessPointsURL        = "/all"

	patchUpdateAccessPointURL = "/"

	softDeleteAccessPointURL = "/sd"
	restoreAccessPointURL    = "/restore"
)

type accessPointHandler struct {
	usecase *usecase.AccessPointUsecase
	// aptMapper  *mapper.AccessPointMapper
	// aprtMapper *mapper.AccessPointRadioTemplateMapper
}

// Регистрирует новый handler
func NewAccessPointHandler(usecase *usecase.AccessPointUsecase) *accessPointHandler {
	return &accessPointHandler{
		usecase: usecase,
		// aptMapper:  &mapper.AccessPointMapper{},
		// aprtMapper: &mapper.AccessPointRadioTemplateMapper{},
	}
}

// Регистрирует маршруты для user
func (h *accessPointHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	// Create
	router.Post(createAccessPointURL, h.CreateAccessPoint)

	// Get
	router.Get(getAccessPointURL, h.GetAccessPoint)
	// router.Get(getAccessPointDetailedURL, h.GetAccessPointDetailed)
	router.Get(getAccessPointsURL, h.GetAccessPoints)

	// Update
	router.Patch(patchUpdateAccessPointURL, h.PatchUpdateAccessPoint)

	// Delete
	router.Patch(softDeleteAccessPointURL, h.SoftDeleteAccessPoint)
	router.Patch(restoreAccessPointURL, h.RestoreAccessPoint)

	return router
}

func (h *accessPointHandler) CreateAccessPoint(ctx *fiber.Ctx) error {
	var httpDTO http_dto.CreateAccessPointDTO
	err := ctx.BodyParser(&httpDTO)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse access point request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse access point request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := (*domain_dto.CreateAccessPointDTO)(&httpDTO)
	// domainDTO := h.aptMapper.CreateHTTPtoDomain(&httpDTO) // Избыточный метод

	accessPointID, err := h.usecase.CreateAccessPoint(context.Background(), domainDTO)
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

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the access point")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the access point",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": accessPointID})
}

func (h *accessPointHandler) GetAccessPoint(ctx *fiber.Ctx) error {
	accessPointID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// var dto http_dto.GetAccessPointDTO = http_dto.GetAccessPointDTO{
	// 	ID: accessPointID,
	// }

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.GetAccessPointDTO{
	// 	ID: dto.ID,
	// }

	ap, err := h.usecase.GetAccessPoint(context.Background(), accessPointID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the access point type")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the access point type",
			"",
			nil,
		))
	}

	// Mapping entity -> http DTO
	// apHttpDTO := h.aptMapper.EntityDomainToHTTP(aptDomainDTO)
	apHttpDTO := (http_dto.AccessPointDTO)(*ap)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": apHttpDTO})
}

// func (h *accessPointHandler) GetAccessPointDetailed(ctx *fiber.Ctx) error {
// 	accessPointID, err := uuid.Parse(ctx.Query("id"))
// 	if err != nil {
// 		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
// 		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
// 			fiber.StatusBadRequest,
// 			"Invalid ID",
// 			"Failed to parse 'id' as UUID",
// 			nil,
// 		))
// 	}

// 	// var dto http_dto.GetAccessPointDTO = http_dto.GetAccessPointDTO{
// 	// 	ID: accessPointID,
// 	// }

// 	// TODO реализовать передачу page и size
// 	var dto http_dto.GetAccessPointDetailedDTO = http_dto.GetAccessPointDetailedDTO{
// 		AccessPointID: accessPointID,
// 		Page:              1,
// 		Size:              100,
// 	}

// 	// TODO validate

// 	// Mapping http DTO -> domain DTO
// 	// domainDTO := domain_dto.GetAccessPointDTO{
// 	// 	ID: dto.ID,
// 	// }

// 	domainDTO := domain_dto.GetAccessPointDetailedDTO{
// 		AccessPointID: dto.AccessPointID,
// 		Limit:             dto.Size,
// 		Offset:            (dto.Page - 1) * dto.Size,
// 	}

// 	accessPointDetailed, err := h.usecase.GetAccessPointDetailed(context.Background(), domainDTO)
// 	if err != nil {
// 		if errors.Is(err, usecase.ErrNotFound) {
// 			ctx.Status(fiber.StatusNoContent)
// 			return nil
// 		}

// 		log.Error().Msg("an unexpected error has occurred while trying to retrieve the access point type detailed")
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
// 			fiber.StatusInternalServerError,
// 			"An unexpected error has occurred while trying to retrieve the access point type detailed",
// 			"",
// 			nil,
// 		))
// 	}

// 	// Mapping access point radio template entity -> http DTO
// 	var aprtHttpDTOs []*http_dto.AccessPointRadioTemplateDTO
// 	for _, aprtHttpDTO := range accessPointDetailed.RadioTemplates {
// 		aprtHttpDTOs = append(aprtHttpDTOs, (*http_dto.AccessPointRadioTemplateDTO)(aprtHttpDTO))
// 	}

// 	// Mapping entity -> http DTO
// 	accessPointDetailedDTO := http_dto.AccessPointDetailedDTO{
// 		AccessPointDTO: (http_dto.AccessPointDTO)(accessPointDetailed.AccessPoint),
// 		RadioTemplatesDTO:  ([]*http_dto.AccessPointRadioTemplateDTO)(aprtHttpDTOs),
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": accessPointDetailedDTO})
// }

func (h *accessPointHandler) GetAccessPoints(c *fiber.Ctx) error {
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
	var dto http_dto.GetAccessPointsDTO = http_dto.GetAccessPointsDTO{
		FloorID: floorID,
		Page:    1,
		Size:    100,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.GetAccessPointsDTO{
		FloorID: dto.FloorID,
		Limit:   dto.Size,
		Offset:  (dto.Page - 1) * dto.Size,
	}

	aptDomainDTOs, err := h.usecase.GetAccessPoints(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the access point")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the access point",
			"",
			nil,
		))
	}

	// Mapping entity -> http DTO
	var apHttpDTOs []http_dto.AccessPointDTO
	for _, aptDomainDTO := range aptDomainDTOs {
		accessPointDTO := (http_dto.AccessPointDTO)(*aptDomainDTO)
		apHttpDTOs = append(apHttpDTOs, accessPointDTO)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": apHttpDTOs})
}

func (h *accessPointHandler) PatchUpdateAccessPoint(c *fiber.Ctx) error {
	var httpDTO http_dto.PatchUpdateAccessPointDTO
	err := c.BodyParser(&httpDTO)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse access point request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse access point request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := (*domain_dto.PatchUpdateAccessPointDTO)(&httpDTO)

	err = h.usecase.PatchUpdateAccessPoint(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The access point was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The access point was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to patch update the access point")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the access point",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *accessPointHandler) SoftDeleteAccessPoint(c *fiber.Ctx) error {
	accessPointID, err := uuid.Parse(c.Query("id"))
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
	// domainDTO := domain_dto.SoftDeleteAccessPointDTO{
	// 	ID: accessPointID,
	// }

	err = h.usecase.SoftDeleteAccessPoint(context.Background(), accessPointID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The access point was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The access point has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to soft delete the access point")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the access point",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *accessPointHandler) RestoreAccessPoint(c *fiber.Ctx) error {
	accessPointID, err := uuid.Parse(c.Query("id"))
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
	// domainDTO := domain_dto.SoftDeleteAccessPointDTO{
	// 	ID: accessPointID,
	// }

	err = h.usecase.RestoreAccessPoint(context.Background(), accessPointID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The access point was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The access point is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to restore the access point")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the access point",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
