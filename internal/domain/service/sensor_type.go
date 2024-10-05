package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
)

type SensorTypeRepo interface {
	Create(ctx context.Context, createSensorTypeDTO *dto.CreateSensorTypeDTO) (sensorTypeID uuid.UUID, err error)
	GetOne(ctx context.Context, sensorTypeID uuid.UUID) (sensorType *entity.SensorType, err error)
	GetAll(ctx context.Context, siteID uuid.UUID, limit, offset int) (sensorTypes []*entity.SensorType, err error)

	Update(ctx context.Context, updateSensorTypeDTO *dto.PatchUpdateSensorTypeDTO) (err error)

	IsSensorTypeSoftDeleted(ctx context.Context, sensorTypeID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, sensorTypeID uuid.UUID) (err error)
	Restore(ctx context.Context, sensorTypeID uuid.UUID) (err error)
}

type sensorTypeService struct {
	sensorTypeRepo SensorTypeRepo
	// sensorRadioTemplateRepo SensorRadioTemplateRepo
}

func NewSensorTypeService(
	sensorTypeRepo SensorTypeRepo,
	// sensorRadioTemplateRepo SensorRadioTemplateRepo,
) *sensorTypeService {
	return &sensorTypeService{
		sensorTypeRepo: sensorTypeRepo,
		// sensorRadioTemplateRepo: sensorRadioTemplateRepo,
	}
}

func (s *sensorTypeService) CreateSensorType(ctx context.Context, createSensorTypeDTO *dto.CreateSensorTypeDTO) (sensorTypeID uuid.UUID, err error) {
	sensorTypeID, err = s.sensorTypeRepo.Create(ctx, createSensorTypeDTO)
	return
}

func (s *sensorTypeService) GetSensorType(ctx context.Context, sensorTypeID uuid.UUID) (sensorType *entity.SensorType, err error) {
	sensorType, err = s.sensorTypeRepo.GetOne(ctx, sensorTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return sensorType, usecase.ErrNotFound
		}

		return
	}

	return
}

// func (s *sensorTypeService) GetSensorTypeDetailed(ctx context.Context, dto dto.GetSensorTypeDetailedDTO) (sensorTypeDetailed *entity.SensorTypeDetailed, err error) {
// 	sensorType, err := s.sensorTypeRepo.GetOne(ctx, dto.ID)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			return sensorTypeDetailed, usecase.ErrNotFound
// 		}

// 		return
// 	}

// 	accessPointRadioTemplates, err := s.accessPointRadioTemplateRepo.GetAll(ctx, dto.ID, dto.Limit, dto.Offset)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			err = nil
// 		} else {
// 			log.Error().Msg("failed to retrieve access point radio template")
// 			return
// 		}
// 	}

// 	sensorTypeDetailed = &entity.SensorTypeDetailed{
// 		SensorType: *sensorType,
// 		RadioTemplates:  accessPointRadioTemplates,
// 	}

// 	return
// }

func (s *sensorTypeService) GetSensorTypes(ctx context.Context, dto dto.GetSensorTypesDTO) (sensorTypes []*entity.SensorType, err error) {
	sensorTypes, err = s.sensorTypeRepo.GetAll(ctx, dto.SiteID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return sensorTypes, usecase.ErrNotFound
		}

		return
	}

	return
}

// TODO PUT update
func (s *sensorTypeService) UpdateSensorType(ctx context.Context, updateSensorTypeDTO *dto.PatchUpdateSensorTypeDTO) (err error) {
	err = s.sensorTypeRepo.Update(ctx, updateSensorTypeDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
		if errors.Is(err, ErrNotUpdated) {
			return usecase.ErrNotUpdated
		}

		return
	}

	return
}

func (s *sensorTypeService) IsSensorTypeSoftDeleted(ctx context.Context, sensorTypeID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.sensorTypeRepo.IsSensorTypeSoftDeleted(ctx, sensorTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *sensorTypeService) SoftDeleteSensorType(ctx context.Context, sensorTypeID uuid.UUID) (err error) {
	err = s.sensorTypeRepo.SoftDelete(ctx, sensorTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *sensorTypeService) RestoreSensorType(ctx context.Context, sensorTypeID uuid.UUID) (err error) {
	err = s.sensorTypeRepo.Restore(ctx, sensorTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}

		return
	}

	return
}
