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
	createWallTypeURL = "/"
	getWallTypeURL    = "/"
	getWallTypesURL   = "/all"

	patchUpdateWallTypeURL = "/"

	softDeleteWallTypeURL = "/sd"
	restoreWallTypeURL    = "/restore"
)

type wallTypeHandler struct {
	usecase usecase.WallTypeUsecase
}

// Регистрирует новый handler
func NewWallTypeHandler(usecase usecase.WallTypeUsecase) *wallTypeHandler {
	return &wallTypeHandler{usecase: usecase}
}

// Регистрирует маршруты для user
func (h *wallTypeHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createWallTypeURL, h.CreateWallType)
	router.Get(getWallTypeURL, h.GetWallType)
	router.Get(getWallTypesURL, h.GetWallTypes)

	router.Patch(patchUpdateWallTypeURL, h.PatchUpdateWallType)

	router.Patch(softDeleteWallTypeURL, h.SoftDeleteWallType)
	router.Patch(restoreWallTypeURL, h.RestoreWallType)

	// TODO Get list detailed
	return router
}

func (h *wallTypeHandler) CreateWallType(ctx *fiber.Ctx) error {
	var dto http_dto.CreateWallTypeDTO
	err := ctx.BodyParser(&dto)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse wallType request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse wallType request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := &domain_dto.CreateWallTypeDTO{
		Name:          dto.Name,
		Color:         dto.Color,
		Attenuation24: dto.Attenuation24, Attenuation5: dto.Attenuation5, Attenuation6: dto.Attenuation6,
		Thickness: dto.Thickness,
		SiteID:    dto.SiteID,
	}

	wallTypeID, err := h.usecase.CreateWallType(context.Background(), domainDTO)
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

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the wallType")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the wallType",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": wallTypeID})
}

func (h *wallTypeHandler) GetWallType(ctx *fiber.Ctx) error {
	wallTypeID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// var dto http_dto.GetWallTypeDTO = http_dto.GetWallTypeDTO{
	// 	ID: wallTypeID,
	// }

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.GetWallTypeDTO{
	// 	ID: dto.ID,
	// }

	wallType, err := h.usecase.GetWallType(context.Background(), wallTypeID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the wallType")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the wallType",
			"",
			nil,
		))
	}

	// Mapping domain DTO -> http DTO
	wallTypeDTO := http_dto.WallTypeDTO{
		ID: wallType.ID,
		CreateWallTypeDTO: http_dto.CreateWallTypeDTO{
			Name:          wallType.Name,
			Color:         wallType.Color,
			Attenuation24: wallType.Attenuation24, Attenuation5: wallType.Attenuation5, Attenuation6: wallType.Attenuation6,
			Thickness: wallType.Thickness,
			SiteID:    wallType.SiteID,
		},
		CreatedAt: wallType.CreatedAt,
		UpdatedAt: wallType.UpdatedAt,
		DeletedAt: wallType.DeletedAt,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": wallTypeDTO})
}

func (h *wallTypeHandler) GetWallTypes(c *fiber.Ctx) error {
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
	var dto http_dto.GetWallTypesDTO = http_dto.GetWallTypesDTO{
		SiteID: siteID,
		Page:   1,
		Size:   100,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.GetWallTypesDTO{
		SiteID: dto.SiteID,
		Limit:  dto.Size,
		Offset: (dto.Page - 1) * dto.Size,
	}

	wallTypes, err := h.usecase.GetWallTypes(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the wallType")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the wallType",
			"",
			nil,
		))
	}

	var wallTypesDTO []http_dto.WallTypeDTO
	for _, wallType := range wallTypes {
		// Mapping domain DTO -> http DTO
		wallTypeDTO := http_dto.WallTypeDTO{
			ID: wallType.ID,
			CreateWallTypeDTO: http_dto.CreateWallTypeDTO{
				Name:          wallType.Name,
				Color:         wallType.Color,
				Attenuation24: wallType.Attenuation24, Attenuation5: wallType.Attenuation5, Attenuation6: wallType.Attenuation6,
				Thickness: wallType.Thickness,
				SiteID:    wallType.SiteID,
			},
			CreatedAt: wallType.CreatedAt,
			UpdatedAt: wallType.UpdatedAt,
			DeletedAt: wallType.DeletedAt,
		}

		wallTypesDTO = append(wallTypesDTO, wallTypeDTO)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": wallTypesDTO})
}

func (h *wallTypeHandler) PatchUpdateWallType(c *fiber.Ctx) error {
	var dto http_dto.PatchUpdateWallTypeDTO
	err := c.BodyParser(&dto)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse wallType request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse wallType request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := &domain_dto.PatchUpdateWallTypeDTO{
		ID:            dto.ID,
		Name:          dto.Name,
		Color:         dto.Color,
		Attenuation24: dto.Attenuation24, Attenuation5: dto.Attenuation5, Attenuation6: dto.Attenuation6,
		Thickness: dto.Thickness,
	}

	err = h.usecase.PatchUpdateWallType(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The wallType was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The wallType was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to patch update the wallType")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the wallType",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *wallTypeHandler) SoftDeleteWallType(c *fiber.Ctx) error {
	wallTypeID, err := uuid.Parse(c.Query("id"))
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
	// domainDTO := domain_dto.SoftDeleteWallTypeDTO{
	// 	ID: wallTypeID,
	// }

	err = h.usecase.SoftDeleteWallType(context.Background(), wallTypeID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The wallType was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The wallType has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to soft delete the wallType")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the wallType",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *wallTypeHandler) RestoreWallType(c *fiber.Ctx) error {
	wallTypeID, err := uuid.Parse(c.Query("id"))
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
	// domainDTO := domain_dto.SoftDeleteWallTypeDTO{
	// 	ID: wallTypeID,
	// }

	err = h.usecase.RestoreWallType(context.Background(), wallTypeID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The wallType was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The wallType is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to restore the wallType")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the wallType",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
