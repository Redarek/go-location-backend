package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type ISensorRadioService interface {
	CreateSensorRadio(ctx context.Context, createDTO *dto.CreateSensorRadioDTO) (sensorRadioID uuid.UUID, err error)
	GetSensorRadio(ctx context.Context, sensorRadioID uuid.UUID) (sensorRadio *entity.SensorRadio, err error)
	GetSensorRadios(ctx context.Context, getDTO dto.GetSensorRadiosDTO) (sensorRadios []*entity.SensorRadio, err error)

	UpdateSensorRadio(ctx context.Context, patchUpdateDTO *dto.PatchUpdateSensorRadioDTO) (err error)

	IsSensorRadioSoftDeleted(ctx context.Context, sensorRadioID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteSensorRadio(ctx context.Context, sensorRadioID uuid.UUID) (err error)
	RestoreSensorRadio(ctx context.Context, sensorRadioID uuid.UUID) (err error)
}

type SensorRadioUsecase struct {
	sensorRadioService ISensorRadioService
	sensorService      ISensorService
}

func NewSensorRadioUsecase(sensorRadioService ISensorRadioService, sensorService ISensorService) *SensorRadioUsecase {
	return &SensorRadioUsecase{
		sensorRadioService: sensorRadioService,
		sensorService:      sensorService,
	}
}

func (u *SensorRadioUsecase) CreateSensorRadio(ctx context.Context, createDTO *dto.CreateSensorRadioDTO) (sensorRadioID uuid.UUID, err error) {
	_, err = u.sensorService.GetSensor(ctx, createDTO.SensorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Msg("failed to create sensor radio: the sensor with provided ID does not exist")
			return sensorRadioID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check sensor existing")
		return
	}

	sensorRadioID, err = u.sensorRadioService.CreateSensorRadio(ctx, createDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create sensor radio")
		return
	}

	log.Debug().Msgf("sensor radio %v successfully created", sensorRadioID)
	return
}

func (u *SensorRadioUsecase) GetSensorRadio(ctx context.Context, sensorRadioID uuid.UUID) (sensorRadio *entity.SensorRadio, err error) {
	sensorRadio, err = u.sensorRadioService.GetSensorRadio(ctx, sensorRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get sensor radio")
			return
		}
	}

	return
}

func (u *SensorRadioUsecase) GetSensorRadios(ctx context.Context, getDTO dto.GetSensorRadiosDTO) (sensorRadios []*entity.SensorRadio, err error) {
	sensorRadios, err = u.sensorRadioService.GetSensorRadios(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get sensor radios")
			return
		}
	}

	return
}

func (u *SensorRadioUsecase) PatchUpdateSensorRadio(ctx context.Context, patchUpdateDTO *dto.PatchUpdateSensorRadioDTO) (err error) {
	_, err = u.sensorRadioService.GetSensorRadio(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Msg("sensor radio was not found")
			return ErrNotFound
		}
	}

	err = u.sensorRadioService.UpdateSensorRadio(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Msg("sensor radio was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update sensor radio")
		return
	}

	return
}

func (u *SensorRadioUsecase) SoftDeleteSensorRadio(ctx context.Context, sensorRadioID uuid.UUID) (err error) {
	isDeleted, err := u.sensorRadioService.IsSensorRadioSoftDeleted(ctx, sensorRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if sensor radio is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.sensorRadioService.SoftDeleteSensorRadio(ctx, sensorRadioID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete sensor radio")
		return
	}

	return
}

func (u *SensorRadioUsecase) RestoreSensorRadio(ctx context.Context, sensorRadioID uuid.UUID) (err error) {
	isDeleted, err := u.sensorRadioService.IsSensorRadioSoftDeleted(ctx, sensorRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if sensor radio is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.sensorRadioService.RestoreSensorRadio(ctx, sensorRadioID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore sensor radio")
		return
	}

	return
}
