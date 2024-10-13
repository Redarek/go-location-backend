package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type SensorRadioTemplateService interface {
	CreateSensorRadioTemplate(ctx context.Context, createDTO *dto.CreateSensorRadioTemplateDTO) (sensorRadioTemplateID uuid.UUID, err error)
	GetSensorRadioTemplate(ctx context.Context, sensorRadioTemplateID uuid.UUID) (sensorRadioTemplate *entity.SensorRadioTemplate, err error)
	GetSensorRadioTemplates(ctx context.Context, getDTO dto.GetSensorRadioTemplatesDTO) (sensorRadioTemplates []*entity.SensorRadioTemplate, err error)

	UpdateSensorRadioTemplate(ctx context.Context, patchUpdateDTO *dto.PatchUpdateSensorRadioTemplateDTO) (err error)

	IsSensorRadioTemplateSoftDeleted(ctx context.Context, sensorRadioTemplateID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteSensorRadioTemplate(ctx context.Context, sensorRadioTemplateID uuid.UUID) (err error)
	RestoreSensorRadioTemplate(ctx context.Context, sensorRadioTemplateID uuid.UUID) (err error)
}

type SensorRadioTemplateUsecase struct {
	sensorRadioTemplateService SensorRadioTemplateService
	sensorTypeService          SensorTypeService
}

func NewSensorRadioTemplateUsecase(sensorRadioTemplateService SensorRadioTemplateService, sensorTypeService SensorTypeService) *SensorRadioTemplateUsecase {
	return &SensorRadioTemplateUsecase{
		sensorRadioTemplateService: sensorRadioTemplateService,
		sensorTypeService:          sensorTypeService,
	}
}

func (u *SensorRadioTemplateUsecase) CreateSensorRadioTemplate(ctx context.Context, createDTO *dto.CreateSensorRadioTemplateDTO) (sensorRadioTemplateID uuid.UUID, err error) {
	_, err = u.sensorTypeService.GetSensorType(ctx, createDTO.SensorTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Msg("failed to create sensor radio template: the sensor type with provided ID does not exist")
			return sensorRadioTemplateID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check sensor type existing")
		return
	}

	sensorRadioTemplateID, err = u.sensorRadioTemplateService.CreateSensorRadioTemplate(ctx, createDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create sensor radio template")
		return
	}

	log.Debug().Msgf("sensor radio template %v successfully created", sensorRadioTemplateID)
	return
}

func (u *SensorRadioTemplateUsecase) GetSensorRadioTemplate(ctx context.Context, sensorRadioTemplateID uuid.UUID) (sensorRadioTemplate *entity.SensorRadioTemplate, err error) {
	sensorRadioTemplate, err = u.sensorRadioTemplateService.GetSensorRadioTemplate(ctx, sensorRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get sensor radio template")
			return
		}
	}

	return
}

func (u *SensorRadioTemplateUsecase) GetSensorRadioTemplates(ctx context.Context, getDTO dto.GetSensorRadioTemplatesDTO) (sensorRadioTemplates []*entity.SensorRadioTemplate, err error) {
	sensorRadioTemplates, err = u.sensorRadioTemplateService.GetSensorRadioTemplates(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get sensor radio templates")
			return
		}
	}

	return
}

func (u *SensorRadioTemplateUsecase) PatchUpdateSensorRadioTemplate(ctx context.Context, patchUpdateDTO *dto.PatchUpdateSensorRadioTemplateDTO) (err error) {
	_, err = u.sensorRadioTemplateService.GetSensorRadioTemplate(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Msg("sensor radio template was not found")
			return ErrNotFound
		}
	}

	err = u.sensorRadioTemplateService.UpdateSensorRadioTemplate(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Msg("sensorRadioTemplate was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update sensor radio template")
		return
	}

	return
}

func (u *SensorRadioTemplateUsecase) SoftDeleteSensorRadioTemplate(ctx context.Context, sensorRadioTemplateID uuid.UUID) (err error) {
	isDeleted, err := u.sensorRadioTemplateService.IsSensorRadioTemplateSoftDeleted(ctx, sensorRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if sensor radio template is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.sensorRadioTemplateService.SoftDeleteSensorRadioTemplate(ctx, sensorRadioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete sensor radio template")
		return
	}

	return
}

func (u *SensorRadioTemplateUsecase) RestoreSensorRadioTemplate(ctx context.Context, sensorRadioTemplateID uuid.UUID) (err error) {
	isDeleted, err := u.sensorRadioTemplateService.IsSensorRadioTemplateSoftDeleted(ctx, sensorRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if sensor radio template is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.sensorRadioTemplateService.RestoreSensorRadioTemplate(ctx, sensorRadioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore sensor radio template")
		return
	}

	return
}
