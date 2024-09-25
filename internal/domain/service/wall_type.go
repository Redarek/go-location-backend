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

type WallTypeService interface {
	CreateWallType(ctx context.Context, createWallTypeDTO *dto.CreateWallTypeDTO) (wallTypeID uuid.UUID, err error)
	GetWallType(ctx context.Context, wallTypeID uuid.UUID) (wallType *entity.WallType, err error)
	GetWallTypes(ctx context.Context, dto dto.GetWallTypesDTO) (wallTypes []*entity.WallType, err error)
	// TODO get wallType list detailed

	UpdateWallType(ctx context.Context, updateWallTypeDTO *dto.PatchUpdateWallTypeDTO) (err error)

	IsWallTypeSoftDeleted(ctx context.Context, wallTypeID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteWallType(ctx context.Context, wallTypeID uuid.UUID) (err error)
	RestoreWallType(ctx context.Context, wallTypeID uuid.UUID) (err error)
}

type wallTypeService struct {
	repository repository.WallTypeRepo
}

func NewWallTypeService(repository repository.WallTypeRepo) *wallTypeService {
	return &wallTypeService{repository: repository}
}

func (s *wallTypeService) CreateWallType(ctx context.Context, createWallTypeDTO *dto.CreateWallTypeDTO) (wallTypeID uuid.UUID, err error) {
	wallTypeID, err = s.repository.Create(ctx, createWallTypeDTO)
	if err != nil {
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to create wallType")
		return
	}

	return wallTypeID, nil
}

func (s *wallTypeService) GetWallType(ctx context.Context, wallTypeID uuid.UUID) (wallType *entity.WallType, err error) {
	wallType, err = s.repository.GetOne(ctx, wallTypeID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return wallType, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve wallType")
		return
	}

	return
}

func (s *wallTypeService) GetWallTypes(ctx context.Context, dto dto.GetWallTypesDTO) (wallTypes []*entity.WallType, err error) {
	wallTypes, err = s.repository.GetAll(ctx, dto.SiteID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return wallTypes, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve wallType")
		return
	}

	return
}

// TODO PUT update
func (s *wallTypeService) UpdateWallType(ctx context.Context, updateWallTypeDTO *dto.PatchUpdateWallTypeDTO) (err error) {
	err = s.repository.Update(ctx, updateWallTypeDTO)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		if errors.Is(err, repository.ErrNotUpdated) {
			return ErrNotUpdated
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to update wallType")
		return
	}

	return
}

func (s *wallTypeService) IsWallTypeSoftDeleted(ctx context.Context, wallTypeID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.repository.IsWallTypeSoftDeleted(ctx, wallTypeID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve wallType")
		return
	}

	return
}

func (s *wallTypeService) SoftDeleteWallType(ctx context.Context, wallTypeID uuid.UUID) (err error) {
	err = s.repository.SoftDelete(ctx, wallTypeID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to soft delete wallType")
		return
	}

	return
}

func (s *wallTypeService) RestoreWallType(ctx context.Context, wallTypeID uuid.UUID) (err error) {
	err = s.repository.Restore(ctx, wallTypeID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to restore wallType")
		return
	}

	return
}
