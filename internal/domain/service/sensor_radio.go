package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
	"location-backend/pkg/utils"
)

type ISensorRadioRepo interface {
	Create(ctx context.Context, createSensorRadioDTO *dto.CreateSensorRadioDTO) (sensorRadioID uuid.UUID, err error)
	GetOne(ctx context.Context, sensorRadioID uuid.UUID) (sensorRadio *entity.SensorRadio, err error)
	GetAll(ctx context.Context, sensorID uuid.UUID, limit, offset int) (sensorRadios []*entity.SensorRadio, err error)

	Update(ctx context.Context, updateSensorRadioDTO *dto.PatchUpdateSensorRadioDTO) (err error)

	IsSensorRadioSoftDeleted(ctx context.Context, sensorRadioID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, sensorRadioID uuid.UUID) (err error)
	Restore(ctx context.Context, sensorRadioID uuid.UUID) (err error)
}

type sensorRadioService struct {
	repository ISensorRadioRepo
}

func NewSensorRadioService(repository ISensorRadioRepo) *sensorRadioService {
	return &sensorRadioService{repository: repository}
}

func (s *sensorRadioService) CreateSensorRadio(ctx context.Context, createSensorRadioDTO *dto.CreateSensorRadioDTO) (sensorRadioID uuid.UUID, err error) {
	sensorRadioID, err = s.repository.Create(ctx, createSensorRadioDTO)
	return
}

func (s *sensorRadioService) GetSensorRadio(ctx context.Context, sensorRadioID uuid.UUID) (sensorRadio *entity.SensorRadio, err error) {
	sensorRadio, err = s.repository.GetOne(ctx, sensorRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return sensorRadio, usecase.ErrNotFound
		}
	}

	return
}

func (s *sensorRadioService) GetSensorRadios(ctx context.Context, dto *dto.GetSensorRadiosDTO) (sensorRadios []*entity.SensorRadio, err error) {
	sensorRadios, err = s.repository.GetAll(ctx, dto.SensorID, dto.Size, utils.GetOffset(dto.Page, dto.Size))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return sensorRadios, usecase.ErrNotFound
		}
	}

	return
}

// TODO PUT update
func (s *sensorRadioService) UpdateSensorRadio(ctx context.Context, updateSensorRadioDTO *dto.PatchUpdateSensorRadioDTO) (err error) {
	err = s.repository.Update(ctx, updateSensorRadioDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
		if errors.Is(err, ErrNotUpdated) {
			return usecase.ErrNotUpdated
		}
	}

	return
}

func (s *sensorRadioService) IsSensorRadioSoftDeleted(ctx context.Context, sensorRadioID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.repository.IsSensorRadioSoftDeleted(ctx, sensorRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, usecase.ErrNotFound
		}
	}

	return
}

func (s *sensorRadioService) SoftDeleteSensorRadio(ctx context.Context, sensorRadioID uuid.UUID) (err error) {
	err = s.repository.SoftDelete(ctx, sensorRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
	}

	return
}

func (s *sensorRadioService) RestoreSensorRadio(ctx context.Context, sensorRadioID uuid.UUID) (err error) {
	err = s.repository.Restore(ctx, sensorRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
	}

	return
}
