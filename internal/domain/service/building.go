package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
)

type IBuildingRepo interface {
	Create(ctx context.Context, createBuildingDTO *dto.CreateBuildingDTO) (buildingID uuid.UUID, err error)
	GetOne(ctx context.Context, buildingID uuid.UUID) (building *entity.Building, err error)
	// GetOneDetailed(ctx context.Context, buildingID uuid.UUID) (building *entity.BuildingDetailed, err error) // TODO

	GetAll(ctx context.Context, siteID uuid.UUID, limit, offset int) (buildings []*entity.Building, err error)

	Update(ctx context.Context, updateBuildingDTO *dto.PatchUpdateBuildingDTO) (err error)

	IsBuildingSoftDeleted(ctx context.Context, buildingID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, buildingID uuid.UUID) (err error)
	Restore(ctx context.Context, buildingID uuid.UUID) (err error)
}

type buildingService struct {
	repository IBuildingRepo
}

func NewBuildingService(repository IBuildingRepo) *buildingService {
	return &buildingService{repository: repository}
}

func (s *buildingService) CreateBuilding(ctx context.Context, createBuildingDTO *dto.CreateBuildingDTO) (buildingID uuid.UUID, err error) {
	buildingID, err = s.repository.Create(ctx, createBuildingDTO)
	return
}

func (s *buildingService) GetBuilding(ctx context.Context, buildingID uuid.UUID) (building *entity.Building, err error) {
	building, err = s.repository.GetOne(ctx, buildingID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return building, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *buildingService) GetBuildings(ctx context.Context, dto dto.GetBuildingsDTO) (buildings []*entity.Building, err error) {
	buildings, err = s.repository.GetAll(ctx, dto.SiteID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return buildings, usecase.ErrNotFound
		}

		return
	}

	return
}

// TODO PUT update
func (s *buildingService) UpdateBuilding(ctx context.Context, updateBuildingDTO *dto.PatchUpdateBuildingDTO) (err error) {
	err = s.repository.Update(ctx, updateBuildingDTO)
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

func (s *buildingService) IsBuildingSoftDeleted(ctx context.Context, buildingID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.repository.IsBuildingSoftDeleted(ctx, buildingID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *buildingService) SoftDeleteBuilding(ctx context.Context, buildingID uuid.UUID) (err error) {
	err = s.repository.SoftDelete(ctx, buildingID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *buildingService) RestoreBuilding(ctx context.Context, buildingID uuid.UUID) (err error) {
	err = s.repository.Restore(ctx, buildingID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}

		return
	}

	return
}
