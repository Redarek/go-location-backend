package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type WallService interface {
	CreateWall(ctx context.Context, createWallDTO *dto.CreateWallDTO) (wallID uuid.UUID, err error)
	GetWall(ctx context.Context, wallID uuid.UUID) (wall *entity.Wall, err error)
	GetWallDetailed(ctx context.Context, wallID uuid.UUID) (wallDetailed *entity.WallDetailed, err error)
	GetWalls(ctx context.Context, getDTO dto.GetWallsDTO) (walls []*entity.Wall, err error)
	GetWallsDetailed(ctx context.Context, dto dto.GetWallsDTO) (wallsDetailed []*entity.WallDetailed, err error)

	UpdateWall(ctx context.Context, updateWallDTO *dto.PatchUpdateWallDTO) (err error)

	IsWallSoftDeleted(ctx context.Context, wallID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteWall(ctx context.Context, wallID uuid.UUID) (err error)
	RestoreWall(ctx context.Context, wallID uuid.UUID) (err error)
}

type WallUsecase struct {
	wallService  WallService
	floorService FloorService
}

func NewWallUsecase(wallService WallService, floorService FloorService) *WallUsecase {
	return &WallUsecase{
		wallService:  wallService,
		floorService: floorService,
	}
}

func (u *WallUsecase) CreateWall(ctx context.Context, dto *dto.CreateWallDTO) (wallID uuid.UUID, err error) {
	_, err = u.floorService.GetFloor(ctx, dto.FloorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Msg("failed to create wall: the floor with provided site ID does not exist")
			return wallID, ErrNotFound
		}

		log.Error().Msg("failed to check floor existing")
		return
	}

	wallID, err = u.wallService.CreateWall(ctx, dto)
	if err != nil {
		log.Error().Msg("failed to create wall")
		return
	}

	log.Info().Msgf("wall %v successfully created", wallID)
	return
}

func (u *WallUsecase) GetWall(ctx context.Context, wallID uuid.UUID) (wall *entity.Wall, err error) {
	wall, err = u.wallService.GetWall(ctx, wallID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Msg("failed to get wall")
			return
		}
	}

	return
}

func (u *WallUsecase) GetWallDetailed(ctx context.Context, wallID uuid.UUID) (wallDetailed *entity.WallDetailed, err error) {
	wallDetailed, err = u.wallService.GetWallDetailed(ctx, wallID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}

		log.Error().Msg("failed to get wall")
		return
	}

	return
}

func (u *WallUsecase) GetWalls(ctx context.Context, dto dto.GetWallsDTO) (walls []*entity.Wall, err error) {
	walls, err = u.wallService.GetWalls(ctx, dto)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Msg("failed to get walls")
			return
		}
	}

	return
}

// TODO uncomment, it works
// func (u *WallUsecase) GetWallsDetailed(ctx context.Context, dto dto.GetWallsDTO) (wallsDetailed []*entity.WallDetailed, err error) {
// 	wallsDetailed, err = u.wallService.GetWallsDetailed(ctx, dto)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			return nil, ErrNotFound
// 		} else {
// 			log.Error().Msg("failed to get walls detailed")
// 			return
// 		}
// 	}

// 	return
// }

func (u *WallUsecase) PatchUpdateWall(ctx context.Context, patchUpdateDTO *dto.PatchUpdateWallDTO) (err error) {
	_, err = u.wallService.GetWall(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Error().Msg("failed to check wall existing")
			return ErrNotFound
		}
	}

	err = u.wallService.UpdateWall(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Err(err).Msg("wall was not updated")
			return ErrNotUpdated
		}
		log.Error().Msg("failed to patch update wall")
		return
	}

	return
}

func (u *WallUsecase) SoftDeleteWall(ctx context.Context, wallID uuid.UUID) (err error) {
	isDeleted, err := u.wallService.IsWallSoftDeleted(ctx, wallID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Msg("failed to check if wall is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.wallService.SoftDeleteWall(ctx, wallID)
	if err != nil {
		log.Error().Msg("failed to soft delete wall")
		return
	}

	return
}

func (u *WallUsecase) RestoreWall(ctx context.Context, wallID uuid.UUID) (err error) {
	isDeleted, err := u.wallService.IsWallSoftDeleted(ctx, wallID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Msg("failed to check if wall is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.wallService.RestoreWall(ctx, wallID)
	if err != nil {
		log.Error().Msg("failed to restore wall")
		return
	}

	return
}
