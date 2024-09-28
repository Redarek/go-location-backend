package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type WallService interface {
	CreateWall(ctx context.Context, createWallDTO *domain_dto.CreateWallDTO) (wallID uuid.UUID, err error)
	GetWall(ctx context.Context, wallID uuid.UUID) (wall *entity.Wall, err error)
	GetWalls(ctx context.Context, dto domain_dto.GetWallsDTO) (walls []*entity.Wall, err error)
	GetWallDetailed(ctx context.Context, wallID uuid.UUID) (wallDetailed *entity.WallDetailed, err error)

	UpdateWall(ctx context.Context, updateWallDTO *domain_dto.PatchUpdateWallDTO) (err error)

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

func (u *WallUsecase) CreateWall(ctx context.Context, dto *domain_dto.CreateWallDTO) (wallID uuid.UUID, err error) {
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

func (u *WallUsecase) GetWall(ctx context.Context, wallID uuid.UUID) (wallDTO *domain_dto.WallDTO, err error) {
	wall, err := u.wallService.GetWall(ctx, wallID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Msg("failed to get wall")
			return
		}
	}

	// Mapping domain entity -> domain DTO
	wallDTO = &domain_dto.WallDTO{
		ID:         wall.ID,
		X1:         wall.X1,
		Y1:         wall.Y1,
		X2:         wall.X2,
		Y2:         wall.Y2,
		WallTypeID: wall.WallTypeID,
		FloorID:    wall.FloorID,
		CreatedAt:  wall.CreatedAt,
		UpdatedAt:  wall.UpdatedAt,
		DeletedAt:  wall.DeletedAt,
	}

	return
}

func (u *WallUsecase) GetWallDetailed(ctx context.Context, wallID uuid.UUID) (wallDetailedDTO *domain_dto.WallDetailedDTO, err error) {
	wallDetailed, err := u.wallService.GetWallDetailed(ctx, wallID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Msg("failed to get wall")
			return
		}
	}

	// Mapping domain entity -> domain DTO
	wallDetailedDTO = &domain_dto.WallDetailedDTO{
		WallDTO: domain_dto.WallDTO{
			ID:         wallDetailed.ID,
			X1:         wallDetailed.X1,
			Y1:         wallDetailed.Y1,
			X2:         wallDetailed.X2,
			Y2:         wallDetailed.Y2,
			WallTypeID: wallDetailed.WallTypeID,
			FloorID:    wallDetailed.FloorID,
			CreatedAt:  wallDetailed.CreatedAt,
			UpdatedAt:  wallDetailed.UpdatedAt,
			DeletedAt:  wallDetailed.DeletedAt,
		},
		WallTypeDTO: domain_dto.WallTypeDTO{
			ID:            wallDetailed.WallType.ID,
			Name:          wallDetailed.WallType.Name,
			Color:         wallDetailed.WallType.Color,
			Attenuation24: wallDetailed.WallType.Attenuation24,
			Attenuation5:  wallDetailed.WallType.Attenuation5,
			Attenuation6:  wallDetailed.WallType.Attenuation6,
			Thickness:     wallDetailed.WallType.Thickness,
			SiteID:        wallDetailed.WallType.SiteID,
			CreatedAt:     wallDetailed.WallType.CreatedAt,
			UpdatedAt:     wallDetailed.WallType.UpdatedAt,
			DeletedAt:     wallDetailed.WallType.DeletedAt,
		},
	}

	return
}

func (u *WallUsecase) GetWalls(ctx context.Context, dto domain_dto.GetWallsDTO) (wallsDTO []*domain_dto.WallDTO, err error) {
	walls, err := u.wallService.GetWalls(ctx, dto)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Msg("failed to get walls")
			return
		}
	}

	for _, wall := range walls {
		// Mapping domain entity -> domain DTO
		wallDTO := &domain_dto.WallDTO{
			ID:         wall.ID,
			X1:         wall.X1,
			Y1:         wall.Y1,
			X2:         wall.X2,
			Y2:         wall.Y2,
			WallTypeID: wall.WallTypeID,
			FloorID:    wall.FloorID,
			CreatedAt:  wall.CreatedAt,
			UpdatedAt:  wall.UpdatedAt,
			DeletedAt:  wall.DeletedAt,
		}

		wallsDTO = append(wallsDTO, wallDTO)
	}

	return
}

func (u *WallUsecase) PatchUpdateWall(ctx context.Context, patchUpdateDTO *domain_dto.PatchUpdateWallDTO) (err error) {
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
