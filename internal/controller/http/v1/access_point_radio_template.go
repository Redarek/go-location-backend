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
	createAccessPointRadioTemplateURL = "/"
	getAccessPointRadioTemplateURL    = "/"
	getAccessPointRadioTemplatesURL   = "/all"

	patchUpdateAccessPointRadioTemplateURL = "/"

	softDeleteAccessPointRadioTemplateURL = "/sd"
	restoreAccessPointRadioTemplateURL    = "/restore"
)

type accessPointRadioTemplateHandler struct {
	usecase *usecase.AccessPointRadioTemplateUsecase
}

// Регистрирует новый handler
func NewAccessPointRadioTemplateHandler(usecase *usecase.AccessPointRadioTemplateUsecase) *accessPointRadioTemplateHandler {
	return &accessPointRadioTemplateHandler{
		usecase: usecase,
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
	var dtoObj dto.CreateAccessPointRadioTemplateDTO
	err := ctx.BodyParser(&dtoObj)
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

	accessPointRadioTemplateID, err := h.usecase.CreateAccessPointRadioTemplate(context.Background(), &dtoObj)
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

	var dtoObj dto.GetAccessPointRadioTemplateDTO = dto.GetAccessPointRadioTemplateDTO{
		ID: aprtID,
	}

	// TODO validate

	aprt, err := h.usecase.GetAccessPointRadioTemplate(context.Background(), dtoObj.ID)
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

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": aprt})
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
	var dtoObj dto.GetAccessPointRadioTemplatesDTO = dto.GetAccessPointRadioTemplatesDTO{
		AccessPointTypeID: accessPointTypeID,
		Page:              1,
		Size:              100,
	}

	// TODO validate

	accessPointRadioTemplates, err := h.usecase.GetAccessPointRadioTemplates(context.Background(), &dtoObj)
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": accessPointRadioTemplates})
}

func (h *accessPointRadioTemplateHandler) PatchUpdateAccessPointRadioTemplate(c *fiber.Ctx) error {
	var dtoObj dto.PatchUpdateAccessPointRadioTemplateDTO
	err := c.BodyParser(&dtoObj)
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

	err = h.usecase.PatchUpdateAccessPointRadioTemplate(context.Background(), &dtoObj)
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
