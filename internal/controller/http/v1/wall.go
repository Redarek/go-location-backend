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
	createWallURL      = "/"
	getWallURL         = "/"
	getWallDetailedURL = "/detailed"
	getWallsURL        = "/all"

	patchUpdateWallURL = "/"

	softDeleteWallURL = "/sd"
	restoreWallURL    = "/restore"
)

type wallHandler struct {
	usecase *usecase.WallUsecase
}

// Регистрирует новый handler
func NewWallHandler(usecase *usecase.WallUsecase) *wallHandler {
	return &wallHandler{usecase: usecase}
}

// Регистрирует маршруты для user
func (h *wallHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createWallURL, h.CreateWall)
	router.Get(getWallURL, h.GetWall)
	router.Get(getWallDetailedURL, h.GetWallDetailed)
	router.Get(getWallsURL, h.GetWalls)

	router.Patch(patchUpdateWallURL, h.PatchUpdateWall)

	router.Patch(softDeleteWallURL, h.SoftDeleteWall)
	router.Patch(restoreWallURL, h.RestoreWall)

	// TODO Get list detailed
	return router
}

func (h *wallHandler) CreateWall(ctx *fiber.Ctx) error {
	var dtoObj dto.CreateWallDTO
	err := ctx.BodyParser(&dtoObj)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse wall request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse wall request body",
			nil,
		))
	}

	// TODO validate

	wallID, err := h.usecase.CreateWall(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			log.Info().Msg("the site with provided 'floor_id' does not exist")
			return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"Invalid request body",
				"The site with provided 'floor_id' does not exist",
				nil,
			))
		}

		log.Error().Msg("an unexpected error has occurred while trying to create the wall")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the wall",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": wallID})
}

func (h *wallHandler) GetWall(ctx *fiber.Ctx) error {
	wallID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	var dtoObj dto.GetWallDTO = dto.GetWallDTO{
		ID: wallID,
	}

	// TODO validate

	wall, err := h.usecase.GetWall(context.Background(), dtoObj.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Msg("an unexpected error has occurred while trying to retrieve the wall")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the wall",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": wall})
}

func (h *wallHandler) GetWallDetailed(ctx *fiber.Ctx) error {
	wallID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	var dtoObj dto.GetWallDTO = dto.GetWallDTO{
		ID: wallID,
	}

	// TODO validate

	wallDetailed, err := h.usecase.GetWallDetailed(context.Background(), dtoObj.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Msg("an unexpected error has occurred while trying to retrieve the wall detailed")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the wall detailed",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": wallDetailed})
}

func (h *wallHandler) GetWalls(c *fiber.Ctx) error {
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
	var dtoObj dto.GetWallsDTO = dto.GetWallsDTO{
		FloorID: floorID,
		Page:    1,
		Size:    100,
	}

	// TODO validate

	walls, err := h.usecase.GetWalls(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Msg("an unexpected error has occurred while trying to retrieve the wall")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the wall",
			"",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": walls})
}

func (h *wallHandler) PatchUpdateWall(c *fiber.Ctx) error {
	var dtoObj dto.PatchUpdateWallDTO
	err := c.BodyParser(&dtoObj)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse wall request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse wall request body",
			nil,
		))
	}

	// TODO validate

	err = h.usecase.PatchUpdateWall(context.Background(), &dtoObj)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The wall was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The wall was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Msg("an unexpected error has occurred while trying to patch update the wall")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the wall",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *wallHandler) SoftDeleteWall(c *fiber.Ctx) error {
	wallID, err := uuid.Parse(c.Query("id"))
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

	err = h.usecase.SoftDeleteWall(context.Background(), wallID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The wall was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The wall has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Msg("an unexpected error has occurred while trying to soft delete the wall")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the wall",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *wallHandler) RestoreWall(c *fiber.Ctx) error {
	wallID, err := uuid.Parse(c.Query("id"))
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

	err = h.usecase.RestoreWall(context.Background(), wallID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The wall was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The wall is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Msg("an unexpected error has occurred while trying to restore the wall")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the wall",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
