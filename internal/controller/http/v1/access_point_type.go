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
	createAccessPointTypeURL      = "/"
	getAccessPointTypeURL         = "/"
	getAccessPointTypeDetailedURL = "/detailed"
	getAccessPointTypesURL        = "/all"

	patchUpdateAccessPointTypeURL = "/"

	softDeleteAccessPointTypeURL = "/sd"
	restoreAccessPointTypeURL    = "/restore"
)

type accessPointTypeHandler struct {
	usecase    *usecase.AccessPointTypeUsecase
	aptMapper  *mapper.AccessPointTypeMapper
	aprtMapper *mapper.AccessPointRadioTemplateMapper
}

// Регистрирует новый handler
func NewAccessPointTypeHandler(usecase *usecase.AccessPointTypeUsecase) *accessPointTypeHandler {
	return &accessPointTypeHandler{
		usecase:    usecase,
		aptMapper:  &mapper.AccessPointTypeMapper{},
		aprtMapper: &mapper.AccessPointRadioTemplateMapper{},
	}
}

// Регистрирует маршруты для user
func (h *accessPointTypeHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createAccessPointTypeURL, h.CreateAccessPointType)
	router.Get(getAccessPointTypeURL, h.GetAccessPointType)
	router.Get(getAccessPointTypeDetailedURL, h.GetAccessPointTypeDetailed)
	router.Get(getAccessPointTypesURL, h.GetAccessPointTypes)

	router.Patch(patchUpdateAccessPointTypeURL, h.PatchUpdateAccessPointType)

	router.Patch(softDeleteAccessPointTypeURL, h.SoftDeleteAccessPointType)
	router.Patch(restoreAccessPointTypeURL, h.RestoreAccessPointType)

	// TODO Get list detailed
	return router
}

func (h *accessPointTypeHandler) CreateAccessPointType(ctx *fiber.Ctx) error {
	var httpDTO http_dto.CreateAccessPointTypeDTO
	err := ctx.BodyParser(&httpDTO)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse accessPointType request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse accessPointType request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := (*domain_dto.CreateAccessPointTypeDTO)(&httpDTO)
	// domainDTO := h.aptMapper.CreateHTTPtoDomain(&httpDTO) // Избыточный метод

	accessPointTypeID, err := h.usecase.CreateAccessPointType(context.Background(), domainDTO)
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

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the accessPointType")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the accessPointType",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": accessPointTypeID})
}

func (h *accessPointTypeHandler) GetAccessPointType(ctx *fiber.Ctx) error {
	accessPointTypeID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// var dto http_dto.GetAccessPointTypeDTO = http_dto.GetAccessPointTypeDTO{
	// 	ID: accessPointTypeID,
	// }

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.GetAccessPointTypeDTO{
	// 	ID: dto.ID,
	// }

	apt, err := h.usecase.GetAccessPointType(context.Background(), accessPointTypeID)
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
	// aptHttpDTO := h.aptMapper.EntityDomainToHTTP(aptDomainDTO)
	aptHttpDTO := (http_dto.AccessPointTypeDTO)(*apt)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": aptHttpDTO})
}

func (h *accessPointTypeHandler) GetAccessPointTypeDetailed(ctx *fiber.Ctx) error {
	accessPointTypeID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// var dto http_dto.GetAccessPointTypeDTO = http_dto.GetAccessPointTypeDTO{
	// 	ID: accessPointTypeID,
	// }

	// TODO реализовать передачу page и size
	var dto http_dto.GetAccessPointTypeDetailedDTO = http_dto.GetAccessPointTypeDetailedDTO{
		AccessPointTypeID: accessPointTypeID,
		Page:              1,
		Size:              100,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.GetAccessPointTypeDTO{
	// 	ID: dto.ID,
	// }

	domainDTO := domain_dto.GetAccessPointTypeDetailedDTO{
		AccessPointTypeID: dto.AccessPointTypeID,
		Limit:             dto.Size,
		Offset:            (dto.Page - 1) * dto.Size,
	}

	accessPointTypeDetailed, err := h.usecase.GetAccessPointTypeDetailed(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Msg("an unexpected error has occurred while trying to retrieve the access point type detailed")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the access point type detailed",
			"",
			nil,
		))
	}

	// Mapping access point radio template entity -> http DTO
	var aprtHttpDTOs []*http_dto.AccessPointRadioTemplateDTO
	for _, aprtHttpDTO := range accessPointTypeDetailed.RadioTemplates {
		aprtHttpDTOs = append(aprtHttpDTOs, (*http_dto.AccessPointRadioTemplateDTO)(aprtHttpDTO))
	}

	// Mapping entity -> http DTO
	accessPointTypeDetailedDTO := http_dto.AccessPointTypeDetailedDTO{
		AccessPointTypeDTO: (http_dto.AccessPointTypeDTO)(accessPointTypeDetailed.AccessPointType),
		RadioTemplatesDTO:  aprtHttpDTOs,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": accessPointTypeDetailedDTO})
}

func (h *accessPointTypeHandler) GetAccessPointTypes(c *fiber.Ctx) error {
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
	var dto http_dto.GetAccessPointTypesDTO = http_dto.GetAccessPointTypesDTO{
		SiteID: siteID,
		Page:   1,
		Size:   100,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.GetAccessPointTypesDTO{
		SiteID: dto.SiteID,
		Limit:  dto.Size,
		Offset: (dto.Page - 1) * dto.Size,
	}

	aptDomainDTOs, err := h.usecase.GetAccessPointTypes(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the accessPointType")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the access point type",
			"",
			nil,
		))
	}

	var aptHttpDTOs []http_dto.AccessPointTypeDTO
	for _, aptDomainDTO := range aptDomainDTOs {
		// Mapping entity -> http DTO
		accessPointTypeDTO := (http_dto.AccessPointTypeDTO)(*aptDomainDTO)
		aptHttpDTOs = append(aptHttpDTOs, accessPointTypeDTO)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": aptHttpDTOs})
}

func (h *accessPointTypeHandler) PatchUpdateAccessPointType(c *fiber.Ctx) error {
	var httpDTO http_dto.PatchUpdateAccessPointTypeDTO
	err := c.BodyParser(&httpDTO)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse accessPointType request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse access point type request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := (*domain_dto.PatchUpdateAccessPointTypeDTO)(&httpDTO)

	err = h.usecase.PatchUpdateAccessPointType(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The access point type was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The access point type was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to patch update the accessPointType")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the access point type",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *accessPointTypeHandler) SoftDeleteAccessPointType(c *fiber.Ctx) error {
	accessPointTypeID, err := uuid.Parse(c.Query("id"))
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
	// domainDTO := domain_dto.SoftDeleteAccessPointTypeDTO{
	// 	ID: accessPointTypeID,
	// }

	err = h.usecase.SoftDeleteAccessPointType(context.Background(), accessPointTypeID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The accessPointType was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The accessPointType has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to soft delete the accessPointType")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the access point type",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *accessPointTypeHandler) RestoreAccessPointType(c *fiber.Ctx) error {
	accessPointTypeID, err := uuid.Parse(c.Query("id"))
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
	// domainDTO := domain_dto.SoftDeleteAccessPointTypeDTO{
	// 	ID: accessPointTypeID,
	// }

	err = h.usecase.RestoreAccessPointType(context.Background(), accessPointTypeID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The access point type was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The access point type is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to restore the accessPointType")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the access point type",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
