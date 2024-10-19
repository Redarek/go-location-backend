package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type ISensorTypeService interface {
	CreateSensorType(ctx context.Context, createDTO *dto.CreateSensorTypeDTO) (sensorTypeID uuid.UUID, err error)
	GetSensorType(ctx context.Context, sensorTypeID uuid.UUID) (sensorType *entity.SensorType, err error)
	// GetSensorTypeDetailed(ctx context.Context, getDTO dto.GetSensorTypeDetailedDTO) (sensorTypeDetailed *entity.SensorTypeDetailed, err error)
	GetSensorTypes(ctx context.Context, getDTO dto.GetSensorTypesDTO) (sensorTypes []*entity.SensorType, err error)

	UpdateSensorType(ctx context.Context, patchUpdateDTO *dto.PatchUpdateSensorTypeDTO) (err error)

	IsSensorTypeSoftDeleted(ctx context.Context, sensorTypeID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteSensorType(ctx context.Context, sensorTypeID uuid.UUID) (err error)
	RestoreSensorType(ctx context.Context, sensorTypeID uuid.UUID) (err error)
}

type SensorTypeUsecase struct {
	sensorTypeService ISensorTypeService
	// accessPointRadioTemplateService AccessPointRadioTemplateService
	siteService ISiteService
}

func NewSensorTypeUsecase(
	sensorTypeService ISensorTypeService,
	// accessPointRadioTemplateService AccessPointRadioTemplateService,
	siteService ISiteService,
) *SensorTypeUsecase {
	return &SensorTypeUsecase{
		sensorTypeService: sensorTypeService,
		// accessPointRadioTemplateService: accessPointRadioTemplateService,
		siteService: siteService,
	}
}

func (u *SensorTypeUsecase) CreateSensorType(ctx context.Context, createDTO *dto.CreateSensorTypeDTO) (sensorTypeID uuid.UUID, err error) {
	_, err = u.siteService.GetSite(ctx, createDTO.SiteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Err(err).Msg("failed to create sensor type: the site with provided site ID does not exist")
			return sensorTypeID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check site existing")
		return
	}

	sensorTypeID, err = u.sensorTypeService.CreateSensorType(ctx, createDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create sensor type")
		return
	}

	log.Info().Msgf("sensor type %v successfully created", sensorTypeID)
	return
}

func (u *SensorTypeUsecase) GetSensorType(ctx context.Context, sensorTypeID uuid.UUID) (sensorType *entity.SensorType, err error) {
	sensorType, err = u.sensorTypeService.GetSensorType(ctx, sensorTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get sensor type")
			return
		}
	}

	return
}

// func (u *SensorTypeUsecase) GetSensorTypeDetailed(ctx context.Context, getDTO dto.GetSensorTypeDetailedDTO) (sensorTypeDetailed *entity.SensorTypeDetailed, err error) {
// 	sensorTypeDetailed, err = u.sensorTypeService.GetSensorTypeDetailed(ctx, getDTO)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			return nil, ErrNotFound
// 		}

// 		log.Error().Err(err).Msg("failed to get access point type detailed")
// 		return
// 	}

// 	return
// }

func (u *SensorTypeUsecase) GetSensorTypes(ctx context.Context, getDTO dto.GetSensorTypesDTO) (sensorTypes []*entity.SensorType, err error) {
	sensorTypes, err = u.sensorTypeService.GetSensorTypes(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to get sensor types")
		return
	}

	return
}

func (u *SensorTypeUsecase) PatchUpdateSensorType(ctx context.Context, patchUpdateDTO *dto.PatchUpdateSensorTypeDTO) (err error) {
	_, err = u.sensorTypeService.GetSensorType(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Error().Err(err).Msg("failed to check sensor type existing")
			return ErrNotFound
		}
	}

	err = u.sensorTypeService.UpdateSensorType(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Err(err).Msg("sensorType was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update sensor type")
		return
	}

	return
}

func (u *SensorTypeUsecase) SoftDeleteSensorType(ctx context.Context, sensorTypeID uuid.UUID) (err error) {
	isDeleted, err := u.sensorTypeService.IsSensorTypeSoftDeleted(ctx, sensorTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if sensor type is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.sensorTypeService.SoftDeleteSensorType(ctx, sensorTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete sensor type")
		return
	}

	return
}

func (u *SensorTypeUsecase) RestoreSensorType(ctx context.Context, sensorTypeID uuid.UUID) (err error) {
	isDeleted, err := u.sensorTypeService.IsSensorTypeSoftDeleted(ctx, sensorTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if sensor type is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.sensorTypeService.RestoreSensorType(ctx, sensorTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore sensor type")
		return
	}

	return
}
