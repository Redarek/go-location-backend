package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	// "github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
)

type SensorRadioTemplateRepo interface {
	Create(ctx context.Context, createSensorRadioTemplateDTO *dto.CreateSensorRadioTemplateDTO) (sensorRadioTemplateID uuid.UUID, err error)
	GetOne(ctx context.Context, sensorRadioTemplateID uuid.UUID) (sensorRadioTemplate *entity.SensorRadioTemplate, err error)
	GetAll(ctx context.Context, sensorTypeID uuid.UUID, limit, offset int) (sensorRadioTemplates []*entity.SensorRadioTemplate, err error)

	Update(ctx context.Context, updateSensorRadioTemplateDTO *dto.PatchUpdateSensorRadioTemplateDTO) (err error)

	IsSensorRadioTemplateSoftDeleted(ctx context.Context, sensorRadioTemplateID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, sensorRadioTemplateID uuid.UUID) (err error)
	Restore(ctx context.Context, sensorRadioTemplateID uuid.UUID) (err error)
}

type sensorRadioTemplateService struct {
	repository SensorRadioTemplateRepo
}

func NewSensorRadioTemplateService(repository SensorRadioTemplateRepo) *sensorRadioTemplateService {
	return &sensorRadioTemplateService{repository: repository}
}

func (s *sensorRadioTemplateService) CreateSensorRadioTemplate(ctx context.Context, createSensorRadioTemplateDTO *dto.CreateSensorRadioTemplateDTO) (sensorRadioTemplateID uuid.UUID, err error) {
	sensorRadioTemplateID, err = s.repository.Create(ctx, createSensorRadioTemplateDTO)
	return
}

func (s *sensorRadioTemplateService) GetSensorRadioTemplate(ctx context.Context, sensorRadioTemplateID uuid.UUID) (sensorRadioTemplate *entity.SensorRadioTemplate, err error) {
	sensorRadioTemplate, err = s.repository.GetOne(ctx, sensorRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return sensorRadioTemplate, usecase.ErrNotFound
		}
	}

	return
}

func (s *sensorRadioTemplateService) GetSensorRadioTemplates(ctx context.Context, dto dto.GetSensorRadioTemplatesDTO) (sensorRadioTemplates []*entity.SensorRadioTemplate, err error) {
	sensorRadioTemplates, err = s.repository.GetAll(ctx, dto.SensorTypeID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return sensorRadioTemplates, usecase.ErrNotFound
		}
	}

	return
}

// TODO PUT update
func (s *sensorRadioTemplateService) UpdateSensorRadioTemplate(ctx context.Context, updateSensorRadioTemplateDTO *dto.PatchUpdateSensorRadioTemplateDTO) (err error) {
	err = s.repository.Update(ctx, updateSensorRadioTemplateDTO)
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

func (s *sensorRadioTemplateService) IsSensorRadioTemplateSoftDeleted(ctx context.Context, sensorRadioTemplateID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.repository.IsSensorRadioTemplateSoftDeleted(ctx, sensorRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, usecase.ErrNotFound
		}
	}

	return
}

func (s *sensorRadioTemplateService) SoftDeleteSensorRadioTemplate(ctx context.Context, sensorRadioTemplateID uuid.UUID) (err error) {
	err = s.repository.SoftDelete(ctx, sensorRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
	}

	return
}

func (s *sensorRadioTemplateService) RestoreSensorRadioTemplate(ctx context.Context, sensorRadioTemplateID uuid.UUID) (err error) {
	err = s.repository.Restore(ctx, sensorRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
	}

	return
}
