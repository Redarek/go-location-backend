package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type WallRepo interface {
	Create(ctx context.Context, dto *dto.CreateWallDTO) (wallID uuid.UUID, err error)
	GetOne(ctx context.Context, wallID uuid.UUID) (wall *entity.Wall, err error)

	GetAll(ctx context.Context, floorID uuid.UUID, limit, offset int) (walls []*entity.Wall, err error)

	Update(ctx context.Context, updateWallDTO *dto.PatchUpdateWallDTO) (err error)

	IsWallSoftDeleted(ctx context.Context, wallID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, wallID uuid.UUID) (err error)
	Restore(ctx context.Context, wallID uuid.UUID) (err error)
}

type WallService interface {
	CreateWall(ctx context.Context, createWallDTO *dto.CreateWallDTO) (wallID uuid.UUID, err error)
	GetWall(ctx context.Context, wallID uuid.UUID) (wall *entity.Wall, err error)
	GetWalls(ctx context.Context, dto dto.GetWallsDTO) (walls []*entity.Wall, err error)
	GetWallDetailed(ctx context.Context, wallID uuid.UUID) (wallDetailed *entity.WallDetailed, err error)

	UpdateWall(ctx context.Context, updateWallDTO *dto.PatchUpdateWallDTO) (err error)

	IsWallSoftDeleted(ctx context.Context, wallID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteWall(ctx context.Context, wallID uuid.UUID) (err error)
	RestoreWall(ctx context.Context, wallID uuid.UUID) (err error)
}

type wallService struct {
	wallRepo     WallRepo
	wallTypeRepo WallTypeRepo
}

func NewWallService(wallRepo WallRepo, wallTypeRepo WallTypeRepo) *wallService {
	return &wallService{
		wallRepo:     wallRepo,
		wallTypeRepo: wallTypeRepo,
	}
}

func (s *wallService) CreateWall(ctx context.Context, createWallDTO *dto.CreateWallDTO) (wallID uuid.UUID, err error) {
	wallID, err = s.wallRepo.Create(ctx, createWallDTO)
	if err != nil {
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to create wall")
		return
	}

	return wallID, nil
}

func (s *wallService) GetWall(ctx context.Context, wallID uuid.UUID) (wall *entity.Wall, err error) {
	wall, err = s.wallRepo.GetOne(ctx, wallID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return wall, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve wall")
		return
	}

	return
}

func (s *wallService) GetWallDetailed(ctx context.Context, wallID uuid.UUID) (wallDetailed *entity.WallDetailed, err error) {
	wall, err := s.wallRepo.GetOne(ctx, wallID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return wallDetailed, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Msg("failed to retrieve wall")
		return
	}

	wallType, err := s.wallTypeRepo.GetOne(ctx, wall.WallTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return wallDetailed, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Msg("failed to retrieve wallType")
		return
	}

	wallDetailed = &entity.WallDetailed{
		Wall:     *wall,
		WallType: *wallType,
	}

	return
}

func (s *wallService) GetWalls(ctx context.Context, dto dto.GetWallsDTO) (walls []*entity.Wall, err error) {
	walls, err = s.wallRepo.GetAll(ctx, dto.FloorID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return walls, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve wall")
		return
	}

	return
}

// TODO PUT update
func (s *wallService) UpdateWall(ctx context.Context, updateWallDTO *dto.PatchUpdateWallDTO) (err error) {
	err = s.wallRepo.Update(ctx, updateWallDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		if errors.Is(err, ErrNotUpdated) {
			return ErrNotUpdated
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to update wall")
		return
	}

	return
}

func (s *wallService) IsWallSoftDeleted(ctx context.Context, wallID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.wallRepo.IsWallSoftDeleted(ctx, wallID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve wall")
		return
	}

	return
}

func (s *wallService) SoftDeleteWall(ctx context.Context, wallID uuid.UUID) (err error) {
	err = s.wallRepo.SoftDelete(ctx, wallID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to soft delete wall")
		return
	}

	return
}

func (s *wallService) RestoreWall(ctx context.Context, wallID uuid.UUID) (err error) {
	err = s.wallRepo.Restore(ctx, wallID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to restore wall")
		return
	}

	return
}
