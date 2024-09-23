package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	repository "location-backend/internal/adapters/db/postgres"
	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type BuildingService interface {
	CreateBuilding(ctx context.Context, userCreate dto.CreateBuildingDTO) (buildingID uuid.UUID, err error)
	GetBuilding(ctx context.Context, buildingID uuid.UUID) (building *entity.Building, err error)
	GetBuildings(ctx context.Context, dto dto.GetBuildingsDTO) (buildings []*entity.Building, err error)
	// TODO get building list detailed

	UpdateBuilding(ctx context.Context, updateBuildingDTO dto.PatchUpdateBuildingDTO) (err error)

	IsBuildingSoftDeleted(ctx context.Context, buildingID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteBuilding(ctx context.Context, buildingID uuid.UUID) (err error)
	RestoreBuilding(ctx context.Context, buildingID uuid.UUID) (err error)
}

type buildingService struct {
	repository repository.BuildingRepo
}

func NewBuildingService(repository repository.BuildingRepo) *buildingService {
	return &buildingService{repository: repository}
}

func (s *buildingService) CreateBuilding(ctx context.Context, createBuildingDTO dto.CreateBuildingDTO) (buildingID uuid.UUID, err error) {
	buildingID, err = s.repository.Create(ctx, createBuildingDTO)
	if err != nil {
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to create building")
		return
	}

	return buildingID, nil
}

func (s *buildingService) GetBuilding(ctx context.Context, buildingID uuid.UUID) (building *entity.Building, err error) {
	building, err = s.repository.GetOne(ctx, buildingID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return building, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve building")
		return
	}

	return
}

func (s *buildingService) GetBuildings(ctx context.Context, dto dto.GetBuildingsDTO) (buildings []*entity.Building, err error) {
	buildings, err = s.repository.GetAll(ctx, dto.SiteID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return buildings, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve building")
		return
	}

	return
}

// TODO PUT update
func (s *buildingService) UpdateBuilding(ctx context.Context, updateBuildingDTO dto.PatchUpdateBuildingDTO) (err error) {
	err = s.repository.Update(ctx, updateBuildingDTO)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		if errors.Is(err, repository.ErrNotUpdated) {
			return ErrNotUpdated
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to update building")
		return
	}

	return
}

func (s *buildingService) IsBuildingSoftDeleted(ctx context.Context, buildingID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.repository.IsBuildingSoftDeleted(ctx, buildingID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve building")
		return
	}

	return
}

func (s *buildingService) SoftDeleteBuilding(ctx context.Context, buildingID uuid.UUID) (err error) {
	err = s.repository.SoftDelete(ctx, buildingID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to soft delete building")
		return
	}

	return
}

func (s *buildingService) RestoreBuilding(ctx context.Context, buildingID uuid.UUID) (err error) {
	err = s.repository.Restore(ctx, buildingID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to restore building")
		return
	}

	return
}
