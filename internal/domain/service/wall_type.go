package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
)

type WallTypeRepo interface {
	Create(ctx context.Context, dto *dto.CreateWallTypeDTO) (wallTypeID uuid.UUID, err error)
	GetOne(ctx context.Context, wallTypeID uuid.UUID) (wallType *entity.WallType, err error)
	// GetOneDetailed(ctx context.Context, wallTypeID uuid.UUID) (wallType *entity.WallTypeDetailed, err error) // TODO

	GetAll(ctx context.Context, siteID uuid.UUID, limit, offset int) (wallTypes []*entity.WallType, err error)

	Update(ctx context.Context, updateWallTypeDTO *dto.PatchUpdateWallTypeDTO) (err error)

	IsWallTypeSoftDeleted(ctx context.Context, wallTypeID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, wallTypeID uuid.UUID) (err error)
	Restore(ctx context.Context, wallTypeID uuid.UUID) (err error)
}

type wallTypeService struct {
	repository WallTypeRepo
}

func NewWallTypeService(repository WallTypeRepo) *wallTypeService {
	return &wallTypeService{repository: repository}
}

func (s *wallTypeService) CreateWallType(ctx context.Context, createWallTypeDTO *dto.CreateWallTypeDTO) (wallTypeID uuid.UUID, err error) {
	wallTypeID, err = s.repository.Create(ctx, createWallTypeDTO)
	return
}

func (s *wallTypeService) GetWallType(ctx context.Context, wallTypeID uuid.UUID) (wallType *entity.WallType, err error) {
	wallType, err = s.repository.GetOne(ctx, wallTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return wallType, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *wallTypeService) GetWallTypes(ctx context.Context, dto dto.GetWallTypesDTO) (wallTypes []*entity.WallType, err error) {
	wallTypes, err = s.repository.GetAll(ctx, dto.SiteID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return wallTypes, usecase.ErrNotFound
		}

		return
	}

	return
}

// TODO PUT update
func (s *wallTypeService) UpdateWallType(ctx context.Context, updateWallTypeDTO *dto.PatchUpdateWallTypeDTO) (err error) {
	err = s.repository.Update(ctx, updateWallTypeDTO)
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

func (s *wallTypeService) IsWallTypeSoftDeleted(ctx context.Context, wallTypeID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.repository.IsWallTypeSoftDeleted(ctx, wallTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *wallTypeService) SoftDeleteWallType(ctx context.Context, wallTypeID uuid.UUID) (err error) {
	err = s.repository.SoftDelete(ctx, wallTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *wallTypeService) RestoreWallType(ctx context.Context, wallTypeID uuid.UUID) (err error) {
	err = s.repository.Restore(ctx, wallTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}

		return
	}

	return
}
