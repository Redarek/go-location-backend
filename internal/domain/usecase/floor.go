package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type IFloorService interface {
	CreateFloor(ctx context.Context, createFloorDTO *dto.CreateFloorDTO) (floorID uuid.UUID, err error)
	GetFloor(ctx context.Context, floorID uuid.UUID) (floor *entity.Floor, err error)
	GetFloors(ctx context.Context, getDTO dto.GetFloorsDTO) (floors []*entity.Floor, err error)
	// TODO get floor list detailed

	UpdateFloor(ctx context.Context, updateFloorDTO *dto.PatchUpdateFloorDTO) (err error)
	UpdateFloorHeatmap(ctx context.Context, floorID uuid.UUID, heatmap string) (err error)

	IsFloorSoftDeleted(ctx context.Context, floorID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteFloor(ctx context.Context, floorID uuid.UUID) (err error)
	RestoreFloor(ctx context.Context, floorID uuid.UUID) (err error)
}

type FloorUsecase struct {
	floorService    IFloorService
	buildingService IBuildingService
}

func NewFloorUsecase(floorService IFloorService, buildingService IBuildingService) *FloorUsecase {
	return &FloorUsecase{
		floorService:    floorService,
		buildingService: buildingService,
	}
}

func (u *FloorUsecase) CreateFloor(ctx context.Context, dto *dto.CreateFloorDTO) (floorID uuid.UUID, err error) {
	_, err = u.buildingService.GetBuilding(ctx, dto.BuildingID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
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

func (u *FloorUsecase) GetFloor(ctx context.Context, floorID uuid.UUID) (floor *entity.Floor, err error) {
	floor, err = u.floorService.GetFloor(ctx, floorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get floor")
			return
		}
	}

	return
}

func (u *FloorUsecase) GetFloors(ctx context.Context, dto dto.GetFloorsDTO) (floors []*entity.Floor, err error) {
	floors, err = u.floorService.GetFloors(ctx, dto)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get floors")
			return
		}
	}

	return
}

func (u *FloorUsecase) PatchUpdateFloor(ctx context.Context, patchUpdateDTO *dto.PatchUpdateFloorDTO) (err error) {
	_, err = u.floorService.GetFloor(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Error().Err(err).Msg("failed to check floor existing")
			return ErrNotFound
		}
	}

	err = u.floorService.UpdateFloor(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Err(err).Msg("floor was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update floor")
		return
	}

	return
}

func (u *FloorUsecase) SoftDeleteFloor(ctx context.Context, floorID uuid.UUID) (err error) {
	isDeleted, err := u.floorService.IsFloorSoftDeleted(ctx, floorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
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

func (u *FloorUsecase) RestoreFloor(ctx context.Context, floorID uuid.UUID) (err error) {
	isDeleted, err := u.floorService.IsFloorSoftDeleted(ctx, floorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
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
