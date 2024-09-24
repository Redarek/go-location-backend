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
	createBuildingURL = "/"
	getBuildingURL    = "/"
	getBuildingsURL   = "/all"

	patchUpdateBuildingURL = "/"

	softDeleteBuildingURL = "/sd"
	restoreBuildingURL    = "/restore"
)

type buildingHandler struct {
	usecase usecase.BuildingUsecase
}

// Регистрирует новый handler
func NewBuildingHandler(usecase usecase.BuildingUsecase) *buildingHandler {
	return &buildingHandler{usecase: usecase}
}

// Регистрирует маршруты для user
func (h *buildingHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createBuildingURL, h.CreateBuilding)
	router.Get(getBuildingURL, h.GetBuilding)
	router.Get(getBuildingsURL, h.GetBuildings)

	router.Patch(patchUpdateBuildingURL, h.PatchUpdateBuilding)

	router.Patch(softDeleteBuildingURL, h.SoftDeleteBuilding)
	router.Patch(restoreBuildingURL, h.RestoreBuilding)

	// TODO Get list detailed
	return router
}

func (h *buildingHandler) CreateBuilding(ctx *fiber.Ctx) error {
	var dto http_dto.CreateBuildingDTO
	err := ctx.BodyParser(&dto)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse building request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse building request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := &domain_dto.CreateBuildingDTO{
		Name:        dto.Name,
		Description: dto.Description,
		Country:     dto.Country,
		City:        dto.City,
		Address:     dto.Address,
		SiteID:      dto.SiteID,
	}

	buildingID, err := h.usecase.CreateBuilding(context.Background(), domainDTO)
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

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the building")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the building",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": buildingID})
}

func (h *buildingHandler) GetBuilding(ctx *fiber.Ctx) error {
	buildingID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	var dto http_dto.GetBuildingDTO = http_dto.GetBuildingDTO{
		ID: buildingID,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.GetBuildingDTO{
	// 	ID: dto.ID,
	// }

	building, err := h.usecase.GetBuilding(context.Background(), dto.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the building")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the building",
			"",
			nil,
		))
	}

	// Mapping domain DTO -> http DTO
	buildingDTO := http_dto.BuildingDTO{
		ID:          building.ID,
		Name:        building.Name,
		Description: building.Description,
		Country:     building.Country,
		City:        building.City,
		Address:     building.Address,
		SiteID:      building.SiteID,
		CreatedAt:   building.CreatedAt,
		UpdatedAt:   building.UpdatedAt,
		DeletedAt:   building.DeletedAt,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": buildingDTO})
}

func (h *buildingHandler) GetBuildings(c *fiber.Ctx) error {
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
	var dto http_dto.GetBuildingsDTO = http_dto.GetBuildingsDTO{
		SiteID: siteID,
		Page:   1,
		Size:   100,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.GetBuildingsDTO{
		SiteID: dto.SiteID,
		Limit:  dto.Size,
		Offset: (dto.Page - 1) * dto.Size,
	}

	buildings, err := h.usecase.GetBuildings(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the building")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the building",
			"",
			nil,
		))
	}

	var buildingsDTO []http_dto.BuildingDTO
	for _, building := range buildings {
		// Mapping domain DTO -> http DTO
		buildingDTO := http_dto.BuildingDTO{
			ID:          building.ID,
			Name:        building.Name,
			Description: building.Description,
			Country:     building.Country,
			City:        building.City,
			Address:     building.Address,
			SiteID:      building.SiteID,
			CreatedAt:   building.CreatedAt,
			UpdatedAt:   building.UpdatedAt,
			DeletedAt:   building.DeletedAt,
		}

		buildingsDTO = append(buildingsDTO, buildingDTO)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": buildingsDTO})
}

func (h *buildingHandler) PatchUpdateBuilding(c *fiber.Ctx) error {
	var dto http_dto.PatchUpdateBuildingDTO
	err := c.BodyParser(&dto)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse building request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse building request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := &domain_dto.PatchUpdateBuildingDTO{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Country:     dto.Country,
		City:        dto.City,
		Address:     dto.Address,
	}

	err = h.usecase.PatchUpdateBuilding(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The building was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The building was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to patch update the building")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the building",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *buildingHandler) SoftDeleteBuilding(c *fiber.Ctx) error {
	buildingID, err := uuid.Parse(c.Query("id"))
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
	// domainDTO := domain_dto.SoftDeleteBuildingDTO{
	// 	ID: buildingID,
	// }

	err = h.usecase.SoftDeleteBuilding(context.Background(), buildingID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The building was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The building has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to soft delete the building")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the building",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *buildingHandler) RestoreBuilding(c *fiber.Ctx) error {
	buildingID, err := uuid.Parse(c.Query("id"))
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
	// domainDTO := domain_dto.SoftDeleteBuildingDTO{
	// 	ID: buildingID,
	// }

	err = h.usecase.RestoreBuilding(context.Background(), buildingID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The building was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The building is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to restore the building")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the building",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
