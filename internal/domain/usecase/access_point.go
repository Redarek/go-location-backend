package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type AccessPointService interface {
	CreateAccessPoint(ctx context.Context, createDTO *dto.CreateAccessPointDTO) (accessPointID uuid.UUID, err error)
	GetAccessPoint(ctx context.Context, accessPointID uuid.UUID) (accessPoint *entity.AccessPoint, err error)
	GetAccessPointDetailed(ctx context.Context, getDTO dto.GetAccessPointDetailedDTO) (accessPointDetailed *entity.AccessPointDetailed, err error)
	GetAccessPoints(ctx context.Context, getDTO dto.GetAccessPointsDTO) (accessPoints []*entity.AccessPoint, err error)
	GetAccessPointsDetailed(ctx context.Context, dto dto.GetAccessPointsDetailedDTO) (accessPoints []*entity.AccessPointDetailed, err error)

	UpdateAccessPoint(ctx context.Context, patchUpdateDTO *dto.PatchUpdateAccessPointDTO) (err error)

	IsAccessPointSoftDeleted(ctx context.Context, accessPointID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteAccessPoint(ctx context.Context, accessPointID uuid.UUID) (err error)
	RestoreAccessPoint(ctx context.Context, accessPointID uuid.UUID) (err error)
}

type AccessPointUsecase struct {
	accessPointService     AccessPointService
	accessPointTypeService AccessPointTypeService
	floorService           FloorService
}

func NewAccessPointUsecase(
	accessPointService AccessPointService,
	accessPointTypeService AccessPointTypeService,
	floorService FloorService,
) *AccessPointUsecase {
	return &AccessPointUsecase{
		accessPointService:     accessPointService,
		accessPointTypeService: accessPointTypeService,
		floorService:           floorService,
	}
}

func (u *AccessPointUsecase) CreateAccessPoint(ctx context.Context, createDTO *dto.CreateAccessPointDTO) (accessPointID uuid.UUID, err error) {
	_, err = u.floorService.GetFloor(ctx, createDTO.FloorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Err(err).Msg("failed to create access point: the floor with provided floor ID does not exist")
			return accessPointID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check floor existing")
		return
	}

	_, err = u.accessPointTypeService.GetAccessPointType(ctx, createDTO.AccessPointTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Err(err).Msg("failed to create access point: the access point type with provided access point type ID does not exist")
			return accessPointID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check access point type existing")
		return
	}

	accessPointID, err = u.accessPointService.CreateAccessPoint(ctx, createDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create access point")
		return
	}

	log.Info().Msgf("access point %v successfully created", accessPointID)
	return
}

func (u *AccessPointUsecase) GetAccessPoint(ctx context.Context, accessPointID uuid.UUID) (accessPoint *entity.AccessPoint, err error) {
	accessPoint, err = u.accessPointService.GetAccessPoint(ctx, accessPointID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to get access point")
		return
	}

	return
}

func (u *AccessPointUsecase) GetAccessPointDetailed(ctx context.Context, getDTO dto.GetAccessPointDetailedDTO) (accessPointDetailed *entity.AccessPointDetailed, err error) {
	accessPointDetailed, err = u.accessPointService.GetAccessPointDetailed(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to get access point detailed")
		return
	}

	return
}

func (u *AccessPointUsecase) GetAccessPoints(ctx context.Context, getDTO dto.GetAccessPointsDTO) (accessPoints []*entity.AccessPoint, err error) {
	accessPoints, err = u.accessPointService.GetAccessPoints(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get access points")
			return
		}
	}

	return
}

func (u *AccessPointUsecase) GetAccessPointsDetailed(ctx context.Context, getDTO dto.GetAccessPointsDetailedDTO) (accessPointsDetailed []*entity.AccessPointDetailed, err error) {
	accessPointsDetailed, err = u.accessPointService.GetAccessPointsDetailed(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get access points detailed")
			return
		}
	}

	return
}

func (u *AccessPointUsecase) PatchUpdateAccessPoint(ctx context.Context, patchUpdateDTO *dto.PatchUpdateAccessPointDTO) (err error) {
	_, err = u.accessPointService.GetAccessPoint(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Error().Err(err).Msg("failed to check access point existing")
			return ErrNotFound
		}
	}

	err = u.accessPointService.UpdateAccessPoint(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Err(err).Msg("access point was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update access point")
		return
	}

	return
}

func (u *AccessPointUsecase) SoftDeleteAccessPoint(ctx context.Context, accessPointID uuid.UUID) (err error) {
	isDeleted, err := u.accessPointService.IsAccessPointSoftDeleted(ctx, accessPointID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if access point is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.accessPointService.SoftDeleteAccessPoint(ctx, accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete access point")
		return
	}

	return
}

func (u *AccessPointUsecase) RestoreAccessPoint(ctx context.Context, accessPointID uuid.UUID) (err error) {
	isDeleted, err := u.accessPointService.IsAccessPointSoftDeleted(ctx, accessPointID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if access point is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.accessPointService.RestoreAccessPoint(ctx, accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore access point")
		return
	}

	return
}
