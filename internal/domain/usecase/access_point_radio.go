package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type IAccessPointRadioService interface {
	CreateAccessPointRadio(ctx context.Context, createDTO *dto.CreateAccessPointRadioDTO) (accessPointRadioID uuid.UUID, err error)
	GetAccessPointRadio(ctx context.Context, accessPointRadioID uuid.UUID) (accessPointRadio *entity.AccessPointRadio, err error)
	GetAccessPointRadios(ctx context.Context, getDTO dto.GetAccessPointRadiosDTO) (accessPointRadios []*entity.AccessPointRadio, err error)

	UpdateAccessPointRadio(ctx context.Context, patchUpdateDTO *dto.PatchUpdateAccessPointRadioDTO) (err error)

	IsAccessPointRadioSoftDeleted(ctx context.Context, accessPointRadioID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteAccessPointRadio(ctx context.Context, accessPointRadioID uuid.UUID) (err error)
	RestoreAccessPointRadio(ctx context.Context, accessPointRadioID uuid.UUID) (err error)
}

type AccessPointRadioUsecase struct {
	accessPointRadioService IAccessPointRadioService
	accessPointService      IAccessPointService
}

func NewAccessPointRadioUsecase(accessPointRadioService IAccessPointRadioService, accessPointService IAccessPointService) *AccessPointRadioUsecase {
	return &AccessPointRadioUsecase{
		accessPointRadioService: accessPointRadioService,
		accessPointService:      accessPointService,
	}
}

func (u *AccessPointRadioUsecase) CreateAccessPointRadio(ctx context.Context, createDTO *dto.CreateAccessPointRadioDTO) (accessPointRadioID uuid.UUID, err error) {
	_, err = u.accessPointService.GetAccessPoint(ctx, createDTO.AccessPointID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Msg("failed to create access point radio: the access point with provided ID does not exist")
			return accessPointRadioID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check access point existing")
		return
	}

	accessPointRadioID, err = u.accessPointRadioService.CreateAccessPointRadio(ctx, createDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create access point radio")
		return
	}

	log.Debug().Msgf("access point radio %v successfully created", accessPointRadioID)
	return
}

func (u *AccessPointRadioUsecase) GetAccessPointRadio(ctx context.Context, accessPointRadioID uuid.UUID) (accessPointRadio *entity.AccessPointRadio, err error) {
	accessPointRadio, err = u.accessPointRadioService.GetAccessPointRadio(ctx, accessPointRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get access point radio")
			return
		}
	}

	return
}

func (u *AccessPointRadioUsecase) GetAccessPointRadios(ctx context.Context, getDTO dto.GetAccessPointRadiosDTO) (accessPointRadios []*entity.AccessPointRadio, err error) {
	accessPointRadios, err = u.accessPointRadioService.GetAccessPointRadios(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get access point radios")
			return
		}
	}

	return
}

func (u *AccessPointRadioUsecase) PatchUpdateAccessPointRadio(ctx context.Context, patchUpdateDTO *dto.PatchUpdateAccessPointRadioDTO) (err error) {
	_, err = u.accessPointRadioService.GetAccessPointRadio(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Msg("access point radio was not found")
			return ErrNotFound
		}
	}

	err = u.accessPointRadioService.UpdateAccessPointRadio(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Msg("access point radio was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update access point radio")
		return
	}

	return
}

func (u *AccessPointRadioUsecase) SoftDeleteAccessPointRadio(ctx context.Context, accessPointRadioID uuid.UUID) (err error) {
	isDeleted, err := u.accessPointRadioService.IsAccessPointRadioSoftDeleted(ctx, accessPointRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if access point radio is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.accessPointRadioService.SoftDeleteAccessPointRadio(ctx, accessPointRadioID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete access point radio")
		return
	}

	return
}

func (u *AccessPointRadioUsecase) RestoreAccessPointRadio(ctx context.Context, accessPointRadioID uuid.UUID) (err error) {
	isDeleted, err := u.accessPointRadioService.IsAccessPointRadioSoftDeleted(ctx, accessPointRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if access point radio is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.accessPointRadioService.RestoreAccessPointRadio(ctx, accessPointRadioID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore access point radio")
		return
	}

	return
}
