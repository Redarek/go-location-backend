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
	createFloorURL = "/"
	getFloorURL    = "/"
	getFloorsURL   = "/all"

	patchUpdateFloorURL = "/"

	softDeleteFloorURL = "/sd"
	restoreFloorURL    = "/restore"
)

type floorHandler struct {
	usecase *usecase.FloorUsecase
}

// Регистрирует новый handler
func NewFloorHandler(usecase *usecase.FloorUsecase) *floorHandler {
	return &floorHandler{usecase: usecase}
}

// Регистрирует маршруты для user
func (h *floorHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createFloorURL, h.CreateFloor)
	router.Get(getFloorURL, h.GetFloor)
	router.Get(getFloorsURL, h.GetFloors)

	router.Patch(patchUpdateFloorURL, h.PatchUpdateFloor)

	router.Patch(softDeleteFloorURL, h.SoftDeleteFloor)
	router.Patch(restoreFloorURL, h.RestoreFloor)

	// TODO Get list detailed
	return router
}

func (h *floorHandler) CreateFloor(ctx *fiber.Ctx) error {
	var dto http_dto.CreateFloorDTO
	err := ctx.BodyParser(&dto)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse floor request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse floor request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := &domain_dto.CreateFloorDTO{
		Name:   dto.Name,
		Number: dto.Number,
		Image:  dto.Image,
		// Heatmap:              dto.Heatmap,
		WidthInPixels:        dto.WidthInPixels,
		HeightInPixels:       dto.HeightInPixels,
		Scale:                dto.Scale,
		CellSizeMeter:        dto.CellSizeMeter,
		NorthAreaIndentMeter: dto.NorthAreaIndentMeter,
		SouthAreaIndentMeter: dto.SouthAreaIndentMeter,
		WestAreaIndentMeter:  dto.WestAreaIndentMeter,
		EastAreaIndentMeter:  dto.EastAreaIndentMeter,
		BuildingID:           dto.BuildingID,
	}

	floorID, err := h.usecase.CreateFloor(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			log.Info().Err(err).Msg("the site with provided 'building_id' does not exist")
			return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"Invalid request body",
				"The site with provided 'building_id' does not exist",
				nil,
			))
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the floor")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the floor",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": floorID})
}

func (h *floorHandler) GetFloor(ctx *fiber.Ctx) error {
	floorID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// var dto http_dto.GetFloorDTO = http_dto.GetFloorDTO{
	// 	ID: floorID,
	// }

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.GetFloorDTO{
	// 	ID: dto.ID,
	// }

	floor, err := h.usecase.GetFloor(context.Background(), floorID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the floor")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the floor",
			"",
			nil,
		))
	}

	// Mapping domain DTO -> http DTO
	floorDTO := http_dto.FloorDTO{
		ID:                   floor.ID,
		Name:                 floor.Name,
		Number:               floor.Number,
		Image:                floor.Image,
		Heatmap:              floor.Heatmap,
		WidthInPixels:        floor.WidthInPixels,
		HeightInPixels:       floor.HeightInPixels,
		Scale:                floor.Scale,
		CellSizeMeter:        floor.CellSizeMeter,
		NorthAreaIndentMeter: floor.NorthAreaIndentMeter,
		SouthAreaIndentMeter: floor.SouthAreaIndentMeter,
		WestAreaIndentMeter:  floor.WestAreaIndentMeter,
		EastAreaIndentMeter:  floor.EastAreaIndentMeter,
		BuildingID:           floor.BuildingID,
		CreatedAt:            floor.CreatedAt,
		UpdatedAt:            floor.UpdatedAt,
		DeletedAt:            floor.DeletedAt,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": floorDTO})
}

func (h *floorHandler) GetFloors(c *fiber.Ctx) error {
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

	// TODO реализовать передачу page и size
	var dto http_dto.GetFloorsDTO = http_dto.GetFloorsDTO{
		BuildingID: buildingID,
		Page:       1,
		Size:       100,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.GetFloorsDTO{
		BuildingID: dto.BuildingID,
		Limit:      dto.Size,
		Offset:     (dto.Page - 1) * dto.Size,
	}

	floors, err := h.usecase.GetFloors(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the floor")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the floor",
			"",
			nil,
		))
	}

	var floorsDTO []http_dto.FloorDTO
	for _, floor := range floors {
		// Mapping domain DTO -> http DTO
		floorDTO := http_dto.FloorDTO{
			ID:                   floor.ID,
			Name:                 floor.Name,
			Number:               floor.Number,
			Image:                floor.Image,
			Heatmap:              floor.Heatmap,
			WidthInPixels:        floor.WidthInPixels,
			HeightInPixels:       floor.HeightInPixels,
			Scale:                floor.Scale,
			CellSizeMeter:        floor.CellSizeMeter,
			NorthAreaIndentMeter: floor.NorthAreaIndentMeter,
			SouthAreaIndentMeter: floor.SouthAreaIndentMeter,
			WestAreaIndentMeter:  floor.WestAreaIndentMeter,
			EastAreaIndentMeter:  floor.EastAreaIndentMeter,
			BuildingID:           floor.BuildingID,
			CreatedAt:            floor.CreatedAt,
			UpdatedAt:            floor.UpdatedAt,
			DeletedAt:            floor.DeletedAt,
		}

		floorsDTO = append(floorsDTO, floorDTO)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": floorsDTO})
}

func (h *floorHandler) PatchUpdateFloor(c *fiber.Ctx) error {
	var dto http_dto.PatchUpdateFloorDTO
	err := c.BodyParser(&dto)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse floor request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse floor request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := &domain_dto.PatchUpdateFloorDTO{
		ID:     dto.ID,
		Name:   dto.Name,
		Number: dto.Number,
		Image:  dto.Image,
		// Heatmap:              dto.Heatmap,
		WidthInPixels:        dto.WidthInPixels,
		HeightInPixels:       dto.HeightInPixels,
		Scale:                dto.Scale,
		CellSizeMeter:        dto.CellSizeMeter,
		NorthAreaIndentMeter: dto.NorthAreaIndentMeter,
		SouthAreaIndentMeter: dto.SouthAreaIndentMeter,
		WestAreaIndentMeter:  dto.WestAreaIndentMeter,
		EastAreaIndentMeter:  dto.EastAreaIndentMeter,
	}

	err = h.usecase.PatchUpdateFloor(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The floor was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The floor was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to patch update the floor")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the floor",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *floorHandler) SoftDeleteFloor(c *fiber.Ctx) error {
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

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.SoftDeleteFloorDTO{
	// 	ID: floorID,
	// }

	err = h.usecase.SoftDeleteFloor(context.Background(), floorID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The floor was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The floor has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to soft delete the floor")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the floor",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *floorHandler) RestoreFloor(c *fiber.Ctx) error {
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

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.SoftDeleteFloorDTO{
	// 	ID: floorID,
	// }

	err = h.usecase.RestoreFloor(context.Background(), floorID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The floor was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The floor is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to restore the floor")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the floor",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
