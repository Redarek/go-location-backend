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
	createAccessPointRadioTemplateURL = "/"
	getAccessPointRadioTemplateURL    = "/"
	getAccessPointRadioTemplatesURL   = "/all"

	patchUpdateAccessPointRadioTemplateURL = "/"

	softDeleteAccessPointRadioTemplateURL = "/sd"
	restoreAccessPointRadioTemplateURL    = "/restore"
)

type accessPointRadioTemplateHandler struct {
	usecase    *usecase.AccessPointRadioTemplateUsecase
	aprtMapper *mapper.AccessPointRadioTemplateMapper
}

// Регистрирует новый handler
func NewAccessPointRadioTemplateHandler(usecase *usecase.AccessPointRadioTemplateUsecase) *accessPointRadioTemplateHandler {
	return &accessPointRadioTemplateHandler{
		usecase:    usecase,
		aprtMapper: &mapper.AccessPointRadioTemplateMapper{},
	}
}

// Регистрирует маршруты для user
func (h *accessPointRadioTemplateHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createAccessPointRadioTemplateURL, h.CreateAccessPointRadioTemplate)
	router.Get(getAccessPointRadioTemplateURL, h.GetAccessPointRadioTemplate)
	router.Get(getAccessPointRadioTemplatesURL, h.GetAccessPointRadioTemplates)

	router.Patch(patchUpdateAccessPointRadioTemplateURL, h.PatchUpdateAccessPointRadioTemplate)

	router.Patch(softDeleteAccessPointRadioTemplateURL, h.SoftDeleteAccessPointRadioTemplate)
	router.Patch(restoreAccessPointRadioTemplateURL, h.RestoreAccessPointRadioTemplate)

	return router
}

func (h *accessPointRadioTemplateHandler) CreateAccessPointRadioTemplate(ctx *fiber.Ctx) error {
	var dto http_dto.CreateAccessPointRadioTemplateDTO
	err := ctx.BodyParser(&dto)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse access point radio template request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse access point radio template request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := h.aprtMapper.CreateHTTPtoDomain(&dto)

	accessPointRadioTemplateID, err := h.usecase.CreateAccessPointRadioTemplate(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			log.Info().Msg("the site with provided 'access_point_type_id' does not exist")
			return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"Invalid request body",
				"The site with provided 'access_point_type_id' does not exist",
				nil,
			))
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the access point radio template")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the access point radio template",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": accessPointRadioTemplateID})
}

func (h *accessPointRadioTemplateHandler) GetAccessPointRadioTemplate(ctx *fiber.Ctx) error {
	aprtID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// var dto http_dto.GetAccessPointRadioTemplateDTO = http_dto.GetAccessPointRadioTemplateDTO{
	// 	ID: accessPointRadioTemplateID,
	// }

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.GetAccessPointRadioTemplateDTO{
	// 	ID: dto.ID,
	// }

	aprtDomainDTO, err := h.usecase.GetAccessPointRadioTemplate(context.Background(), aprtID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the access point radio template")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the access point radio template",
			"",
			nil,
		))
	}

	// Mapping entity -> http DTO
	aprtHttpDTO := h.aprtMapper.EntityDomainToHTTP(aprtDomainDTO)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": aprtHttpDTO})
}

func (h *accessPointRadioTemplateHandler) GetAccessPointRadioTemplates(c *fiber.Ctx) error {
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

	// TODO реализовать передачу page и size
	var dto http_dto.GetAccessPointRadioTemplatesDTO = http_dto.GetAccessPointRadioTemplatesDTO{
		AccessPointTypeID: accessPointTypeID,
		Page:              1,
		Size:              100,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.GetAccessPointRadioTemplatesDTO{
		AccessPointTypeID: dto.AccessPointTypeID,
		Limit:             dto.Size,
		Offset:            (dto.Page - 1) * dto.Size,
	}

	accessPointRadioTemplates, err := h.usecase.GetAccessPointRadioTemplates(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the access point radio templates")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the access point radio templates",
			"",
			nil,
		))
	}

	var aprtHttpDTOs []*http_dto.AccessPointRadioTemplateDTO
	for _, aprtDomainDTO := range accessPointRadioTemplates {
		// Mapping entity -> http DTO
		aprtHttpDTO := h.aprtMapper.EntityDomainToHTTP(aprtDomainDTO)
		aprtHttpDTOs = append(aprtHttpDTOs, aprtHttpDTO)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": aprtHttpDTOs})
}

func (h *accessPointRadioTemplateHandler) PatchUpdateAccessPointRadioTemplate(c *fiber.Ctx) error {
	var httpDTO http_dto.PatchUpdateAccessPointRadioTemplateDTO
	err := c.BodyParser(&httpDTO)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse access point radio template request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse access point radio template request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := h.aprtMapper.UpdateHTTPtoDomain(&httpDTO)

	err = h.usecase.PatchUpdateAccessPointRadioTemplate(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The access point radio template was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The access point radio template was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to patch update the access point radio template")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the access point radio template",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *accessPointRadioTemplateHandler) SoftDeleteAccessPointRadioTemplate(c *fiber.Ctx) error {
	accessPointRadioTemplateID, err := uuid.Parse(c.Query("id"))
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
	// domainDTO := domain_dto.SoftDeleteAccessPointRadioTemplateDTO{
	// 	ID: accessPointRadioTemplateID,
	// }

	err = h.usecase.SoftDeleteAccessPointRadioTemplate(context.Background(), accessPointRadioTemplateID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The access point radio template was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The access point radio template has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to soft delete the access point radio template")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the access point radio template",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *accessPointRadioTemplateHandler) RestoreAccessPointRadioTemplate(c *fiber.Ctx) error {
	accessPointRadioTemplateID, err := uuid.Parse(c.Query("id"))
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
	// domainDTO := domain_dto.SoftDeleteAccessPointRadioTemplateDTO{
	// 	ID: accessPointRadioTemplateID,
	// }

	err = h.usecase.RestoreAccessPointRadioTemplate(context.Background(), accessPointRadioTemplateID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The access point radio template was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The access point radio template is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to restore the access point radio template")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the access point radio template",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
