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
	createAccessPointRadioURL = "/"
	getAccessPointRadioURL    = "/"
	getAccessPointRadiosURL   = "/all"

	patchUpdateAccessPointRadioURL = "/"

	softDeleteAccessPointRadioURL = "/sd"
	restoreAccessPointRadioURL    = "/restore"
)

type accessPointRadioHandler struct {
	usecase *usecase.AccessPointRadioUsecase
}

// Регистрирует новый handler
func NewAccessPointRadioHandler(usecase *usecase.AccessPointRadioUsecase) *accessPointRadioHandler {
	return &accessPointRadioHandler{
		usecase: usecase,
	}
}

// Регистрирует маршруты для user
func (h *accessPointRadioHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createAccessPointRadioURL, h.CreateAccessPointRadio)
	router.Get(getAccessPointRadioURL, h.GetAccessPointRadio)
	router.Get(getAccessPointRadiosURL, h.GetAccessPointRadios)

	router.Patch(patchUpdateAccessPointRadioURL, h.PatchUpdateAccessPointRadio)

	router.Patch(softDeleteAccessPointRadioURL, h.SoftDeleteAccessPointRadio)
	router.Patch(restoreAccessPointRadioURL, h.RestoreAccessPointRadio)

	return router
}

func (h *accessPointRadioHandler) CreateAccessPointRadio(ctx *fiber.Ctx) error {
	var dtoObj dto.CreateAccessPointRadioDTO
	err := ctx.BodyParser(&dtoObj)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse access point radio request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse access point radio request body",
			nil,
		))
	}

	// TODO validate

	accessPointRadioID, err := h.usecase.CreateAccessPointRadio(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			log.Info().Msg("the access point with provided 'access_point_id' does not exist")
			return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"Invalid request body",
				"The access point with provided 'access_point_id' does not exist",
				nil,
			))
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the access point radio")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the access point radio",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": accessPointRadioID})
}

func (h *accessPointRadioHandler) GetAccessPointRadio(ctx *fiber.Ctx) error {
	accessPointRadioID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	var dtoObj dto.GetAccessPointRadioDTO = dto.GetAccessPointRadioDTO{
		ID: accessPointRadioID,
	}

	// TODO validate

	accessPointRadio, err := h.usecase.GetAccessPointRadio(context.Background(), dtoObj.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the access point radio")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the access point radio",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": accessPointRadio})
}

func (h *accessPointRadioHandler) GetAccessPointRadios(c *fiber.Ctx) error {
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

	// TODO реализовать передачу page и size
	var dtoObj dto.GetAccessPointRadiosDTO = dto.GetAccessPointRadiosDTO{
		AccessPointID: accessPointID,
		Page:          1,
		Size:          100,
	}

	// TODO validate

	accessPointRadios, err := h.usecase.GetAccessPointRadios(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the access point radios")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the access point radios",
			"",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": accessPointRadios})
}

func (h *accessPointRadioHandler) PatchUpdateAccessPointRadio(c *fiber.Ctx) error {
	var dtoObj dto.PatchUpdateAccessPointRadioDTO
	err := c.BodyParser(&dtoObj)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse access point radio request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse access point radio request body",
			nil,
		))
	}

	// TODO validate

	err = h.usecase.PatchUpdateAccessPointRadio(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The access point radio was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The access point radio was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to patch update the access point radio")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the access point radio",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *accessPointRadioHandler) SoftDeleteAccessPointRadio(c *fiber.Ctx) error {
	accessPointRadioID, err := uuid.Parse(c.Query("id"))
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

	err = h.usecase.SoftDeleteAccessPointRadio(context.Background(), accessPointRadioID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The access point radio was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The access point radio has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to soft delete the access point radio")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the access point radio",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *accessPointRadioHandler) RestoreAccessPointRadio(c *fiber.Ctx) error {
	accessPointRadioID, err := uuid.Parse(c.Query("id"))
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

	err = h.usecase.RestoreAccessPointRadio(context.Background(), accessPointRadioID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The access point radio was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The access point radio is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to restore the access point radio")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the access point radio",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
