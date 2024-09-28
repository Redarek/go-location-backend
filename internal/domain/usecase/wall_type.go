package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type WallTypeService interface {
	CreateWallType(ctx context.Context, createWallTypeDTO *domain_dto.CreateWallTypeDTO) (wallTypeID uuid.UUID, err error)
	GetWallType(ctx context.Context, wallTypeID uuid.UUID) (wallType *entity.WallType, err error)
	GetWallTypes(ctx context.Context, dto domain_dto.GetWallTypesDTO) (wallTypes []*entity.WallType, err error)
	// TODO get wallType list detailed

	UpdateWallType(ctx context.Context, updateWallTypeDTO *domain_dto.PatchUpdateWallTypeDTO) (err error)

	IsWallTypeSoftDeleted(ctx context.Context, wallTypeID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteWallType(ctx context.Context, wallTypeID uuid.UUID) (err error)
	RestoreWallType(ctx context.Context, wallTypeID uuid.UUID) (err error)
}

type WallTypeUsecase interface {
	CreateWallType(ctx context.Context, dto *domain_dto.CreateWallTypeDTO) (wallTypeID uuid.UUID, err error)
	GetWallType(ctx context.Context, wallTypeID uuid.UUID) (wallTypeDTO *domain_dto.WallTypeDTO, err error)
	GetWallTypes(ctx context.Context, dto domain_dto.GetWallTypesDTO) (wallTypesDTO []*domain_dto.WallTypeDTO, err error)
	// TODO GetWallTypesDetailed

	PatchUpdateWallType(ctx context.Context, patchUpdateDTO *domain_dto.PatchUpdateWallTypeDTO) (err error)

	SoftDeleteWallType(ctx context.Context, wallTypeID uuid.UUID) (err error)
	RestoreWallType(ctx context.Context, wallTypeID uuid.UUID) (err error)
}

type wallTypeUsecase struct {
	wallTypeService WallTypeService
	siteService     SiteService
}

func NewWallTypeUsecase(wallTypeService WallTypeService, siteService SiteService) *wallTypeUsecase {
	return &wallTypeUsecase{
		wallTypeService: wallTypeService,
		siteService:     siteService,
	}
}

func (u *wallTypeUsecase) CreateWallType(ctx context.Context, dto *domain_dto.CreateWallTypeDTO) (wallTypeID uuid.UUID, err error) {
	_, err = u.siteService.GetSite(ctx, dto.SiteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Err(err).Msg("failed to create wallType: the site with provided site ID does not exist")
			return wallTypeID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check site existing")
		return
	}

	wallTypeID, err = u.wallTypeService.CreateWallType(ctx, dto)
	if err != nil {
		log.Error().Err(err).Msg("failed to create wallType")
		return
	}

	log.Info().Msgf("wallType %v successfully created", dto.Name)
	return
}

func (u *wallTypeUsecase) GetWallType(ctx context.Context, wallTypeID uuid.UUID) (wallTypeDTO *domain_dto.WallTypeDTO, err error) {
	wallType, err := u.wallTypeService.GetWallType(ctx, wallTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get wallType")
			return
		}
	}

	// Mapping domain entity -> domain DTO
	wallTypeDTO = &domain_dto.WallTypeDTO{
		ID:            wallType.ID,
		Name:          wallType.Name,
		Color:         wallType.Color,
		Attenuation24: wallType.Attenuation24, Attenuation5: wallType.Attenuation5, Attenuation6: wallType.Attenuation6,
		Thickness: wallType.Thickness,
		SiteID:    wallType.SiteID,
		CreatedAt: wallType.CreatedAt,
		UpdatedAt: wallType.UpdatedAt,
		DeletedAt: wallType.DeletedAt,
	}

	return
}

func (u *wallTypeUsecase) GetWallTypes(ctx context.Context, dto domain_dto.GetWallTypesDTO) (wallTypesDTO []*domain_dto.WallTypeDTO, err error) {
	wallTypes, err := u.wallTypeService.GetWallTypes(ctx, dto)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get wallTypes")
			return
		}
	}

	for _, wallType := range wallTypes {
		// Mapping domain entity -> domain DTO
		wallTypeDTO := &domain_dto.WallTypeDTO{
			ID:            wallType.ID,
			Name:          wallType.Name,
			Color:         wallType.Color,
			Attenuation24: wallType.Attenuation24, Attenuation5: wallType.Attenuation5, Attenuation6: wallType.Attenuation6,
			Thickness: wallType.Thickness,
			SiteID:    wallType.SiteID,
			CreatedAt: wallType.CreatedAt,
			UpdatedAt: wallType.UpdatedAt,
			DeletedAt: wallType.DeletedAt,
		}

		wallTypesDTO = append(wallTypesDTO, wallTypeDTO)
	}

	return
}

func (u *wallTypeUsecase) PatchUpdateWallType(ctx context.Context, patchUpdateDTO *domain_dto.PatchUpdateWallTypeDTO) (err error) {
	_, err = u.wallTypeService.GetWallType(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Error().Err(err).Msg("failed to check wallType existing")
			return ErrNotFound
		}
	}

	err = u.wallTypeService.UpdateWallType(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Err(err).Msg("wallType was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update wallType")
		return
	}

	return
}

func (u *wallTypeUsecase) SoftDeleteWallType(ctx context.Context, wallTypeID uuid.UUID) (err error) {
	isDeleted, err := u.wallTypeService.IsWallTypeSoftDeleted(ctx, wallTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if wallType is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.wallTypeService.SoftDeleteWallType(ctx, wallTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete wallType")
		return
	}

	return
}

func (u *wallTypeUsecase) RestoreWallType(ctx context.Context, wallTypeID uuid.UUID) (err error) {
	isDeleted, err := u.wallTypeService.IsWallTypeSoftDeleted(ctx, wallTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if wallType is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.wallTypeService.RestoreWallType(ctx, wallTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore wallType")
		return
	}

	return
}
