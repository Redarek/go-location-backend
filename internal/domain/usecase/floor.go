package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/service"
)

type FloorUsecase interface {
	CreateFloor(ctx context.Context, dto *domain_dto.CreateFloorDTO) (floorID uuid.UUID, err error)
	GetFloor(ctx context.Context, floorID uuid.UUID) (floorDTO *domain_dto.FloorDTO, err error)
	GetFloors(ctx context.Context, dto domain_dto.GetFloorsDTO) (floorsDTO []*domain_dto.FloorDTO, err error)
	// TODO GetFloorsDetailed

	PatchUpdateFloor(ctx context.Context, patchUpdateDTO *domain_dto.PatchUpdateFloorDTO) (err error)

	SoftDeleteFloor(ctx context.Context, floorID uuid.UUID) (err error)
	RestoreFloor(ctx context.Context, floorID uuid.UUID) (err error)
}

type floorUsecase struct {
	floorService    service.FloorService
	buildingService service.BuildingService
}

func NewFloorUsecase(floorService service.FloorService, buildingService service.BuildingService) *floorUsecase {
	return &floorUsecase{
		floorService:    floorService,
		buildingService: buildingService,
	}
}

func (u *floorUsecase) CreateFloor(ctx context.Context, dto *domain_dto.CreateFloorDTO) (floorID uuid.UUID, err error) {
	_, err = u.buildingService.GetBuilding(ctx, dto.BuildingID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			log.Info().Err(err).Msg("failed to create floor: the building with provided building ID does not exist")
			return floorID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check building existing")
		return
	}

	floorID, err = u.floorService.CreateFloor(ctx, dto)
	if err != nil {
		log.Error().Err(err).Msg("failed to create floor")
		return
	}

	log.Info().Msgf("floor %v successfully created", dto.Name)
	return
}

func (u *floorUsecase) GetFloor(ctx context.Context, floorID uuid.UUID) (floorDTO *domain_dto.FloorDTO, err error) {
	floor, err := u.floorService.GetFloor(ctx, floorID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get floor")
			return
		}
	}

	// Mapping domain entity -> domain DTO
	floorDTO = &domain_dto.FloorDTO{
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

	return
}

func (u *floorUsecase) GetFloors(ctx context.Context, dto domain_dto.GetFloorsDTO) (floorsDTO []*domain_dto.FloorDTO, err error) {
	floors, err := u.floorService.GetFloors(ctx, dto)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get floors")
			return
		}
	}

	for _, floor := range floors {
		// Mapping domain entity -> domain DTO
		floorDTO := &domain_dto.FloorDTO{
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

	return
}

func (u *floorUsecase) PatchUpdateFloor(ctx context.Context, patchUpdateDTO *domain_dto.PatchUpdateFloorDTO) (err error) {
	_, err = u.floorService.GetFloor(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			log.Error().Err(err).Msg("failed to check floor existing")
			return ErrNotFound
		}
	}

	err = u.floorService.UpdateFloor(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, service.ErrNotUpdated) {
			log.Info().Err(err).Msg("floor was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update floor")
		return
	}

	return
}

func (u *floorUsecase) SoftDeleteFloor(ctx context.Context, floorID uuid.UUID) (err error) {
	isDeleted, err := u.floorService.IsFloorSoftDeleted(ctx, floorID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if floor is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.floorService.SoftDeleteFloor(ctx, floorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete floor")
		return
	}

	return
}

func (u *floorUsecase) RestoreFloor(ctx context.Context, floorID uuid.UUID) (err error) {
	isDeleted, err := u.floorService.IsFloorSoftDeleted(ctx, floorID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if floor is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.floorService.RestoreFloor(ctx, floorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore floor")
		return
	}

	return
}
