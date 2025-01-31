package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
	"location-backend/pkg/utils"
)

type IFloorRepo interface {
	Create(ctx context.Context, createFloorDTO *dto.CreateFloorDTO) (floorID uuid.UUID, err error)
	GetOne(ctx context.Context, floorID uuid.UUID) (floor *entity.Floor, err error)
	GetAll(ctx context.Context, buildingID uuid.UUID, limit, offset int) (floors []*entity.Floor, err error)

	Update(ctx context.Context, patchUpdateFloorDTO *dto.PatchUpdateFloorDTO) (err error)
	UpdateHeatmap(ctx context.Context, floorID uuid.UUID, heatmap string) (err error)

	IsFloorSoftDeleted(ctx context.Context, floorID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, floorID uuid.UUID) (err error)
	Restore(ctx context.Context, floorID uuid.UUID) (err error)
}

type floorService struct {
	repository IFloorRepo
}

func NewFloorService(repository IFloorRepo) *floorService {
	return &floorService{repository: repository}
}

func (s *floorService) CreateFloor(ctx context.Context, createFloorDTO *dto.CreateFloorDTO) (floorID uuid.UUID, err error) {
	floorID, err = s.repository.Create(ctx, createFloorDTO)
	return
}

func (s *floorService) GetFloor(ctx context.Context, floorID uuid.UUID) (floor *entity.Floor, err error) {
	floor, err = s.repository.GetOne(ctx, floorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return floor, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *floorService) GetFloors(ctx context.Context, dto *dto.GetFloorsDTO) (floors []*entity.Floor, err error) {
	floors, err = s.repository.GetAll(ctx, dto.BuildingID, dto.Size, utils.GetOffset(dto.Page, dto.Size))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return floors, usecase.ErrNotFound
		}

		return
	}

	return
}

// TODO PUT update
func (s *floorService) UpdateFloor(ctx context.Context, updateFloorDTO *dto.PatchUpdateFloorDTO) (err error) {
	err = s.repository.Update(ctx, updateFloorDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
		if errors.Is(err, ErrNotUpdated) {
			return usecase.ErrNotUpdated
		}

		return
	}

	return
}

func (s *floorService) UpdateFloorHeatmap(ctx context.Context, floorID uuid.UUID, heatmap string) (err error) {
	err = s.repository.UpdateHeatmap(ctx, floorID, heatmap)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
		if errors.Is(err, ErrNotUpdated) {
			return usecase.ErrNotUpdated
		}

		return
	}

	return
}

func (s *floorService) IsFloorSoftDeleted(ctx context.Context, floorID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.repository.IsFloorSoftDeleted(ctx, floorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *floorService) SoftDeleteFloor(ctx context.Context, floorID uuid.UUID) (err error) {
	err = s.repository.SoftDelete(ctx, floorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *floorService) RestoreFloor(ctx context.Context, floorID uuid.UUID) (err error) {
	err = s.repository.Restore(ctx, floorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}

		return
	}

	return
}
