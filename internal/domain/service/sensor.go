package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
)

type ISensorRepo interface {
	Create(ctx context.Context, createSensorDTO *dto.CreateSensorDTO) (sensorID uuid.UUID, err error)
	GetOne(ctx context.Context, sensorID uuid.UUID) (sensor *entity.Sensor, err error)
	GetOneByMAC(ctx context.Context, mac string) (sensor *entity.Sensor, err error)
	GetOneDetailed(ctx context.Context, sensorID uuid.UUID) (sensorDetailed *entity.SensorDetailed, err error)
	GetAll(ctx context.Context, floorID uuid.UUID, limit, offset int) (sensors []*entity.Sensor, err error)
	GetAllDetailed(ctx context.Context, floorID uuid.UUID, limit, offset int) (sensorsDetailed []*entity.SensorDetailed, err error)

	Update(ctx context.Context, updateSensorDTO *dto.PatchUpdateSensorDTO) (err error)

	IsSensorSoftDeleted(ctx context.Context, sensorID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, sensorID uuid.UUID) (err error)
	Restore(ctx context.Context, sensorID uuid.UUID) (err error)
}

type sensorService struct {
	sensorRepo     ISensorRepo
	sensorTypeRepo ISensorTypeRepo
	// sensorRadioRepo SensorRadioRepo
}

func NewSensorService(
	sensorRepo ISensorRepo,
	sensorTypeRepo ISensorTypeRepo,
	// sensorRadioRepo SensorRadioRepo,
) *sensorService {
	return &sensorService{
		sensorRepo:     sensorRepo,
		sensorTypeRepo: sensorTypeRepo,
		// sensorRadioRepo: sensorRadioRepo,
	}
}

func (s *sensorService) CreateSensor(ctx context.Context, createSensorDTO *dto.CreateSensorDTO) (sensorID uuid.UUID, err error) {
	sensorID, err = s.sensorRepo.Create(ctx, createSensorDTO)
	return
}

func (s *sensorService) GetSensor(ctx context.Context, sensorID uuid.UUID) (sensor *entity.Sensor, err error) {
	sensor, err = s.sensorRepo.GetOne(ctx, sensorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return sensor, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *sensorService) GetSensorByMAC(ctx context.Context, mac string) (sensor *entity.Sensor, err error) {
	sensor, err = s.sensorRepo.GetOneByMAC(ctx, mac)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return sensor, usecase.ErrNotFound
		}

		return
	}

	return
}

// func (s *sensorService) GetSensorDetailed(ctx context.Context, dto dto.GetSensorDetailedDTO) (sensorDetailed *entity.SensorDetailed, err error) {
// 	sensor, err := s.sensorRepo.GetOne(ctx, dto.ID)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			log.Debug().Msg("access point was not found")
// 			return sensorDetailed, usecase.ErrNotFound
// 		}
// 		// TODO улучшить лог
// 		log.Error().Msg("failed to retrieve access point")
// 		return
// 	}

// 	sensorType, err := s.sensorTypeRepo.GetOne(ctx, sensor.SensorTypeID)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			log.Error().Msg("access point type was not found")
// 			return
// 		}

// 		// TODO улучшить лог
// 		log.Error().Msg("failed to retrieve access point type")
// 		return
// 	}

// 	sensorRadios, err := s.sensorRadioRepo.GetAll(ctx, sensor.ID, dto.Limit, dto.Offset)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			log.Debug().Msg("access point radios were not found")
// 			err = nil
// 		} else {
// 			// TODO улучшить лог
// 			log.Error().Msg("failed to retrieve access point radios")
// 			return
// 		}
// 	}

// 	sensorDetailed = &entity.SensorDetailed{
// 		Sensor:     *sensor,
// 		SensorType: *sensorType,
// 		Radios:          sensorRadios,
// 	}

// 	return
// }

func (s *sensorService) GetSensorDetailed(ctx context.Context, dto dto.GetSensorDetailedDTO) (apDetailed *entity.SensorDetailed, err error) {
	apDetailed, err = s.sensorRepo.GetOneDetailed(ctx, dto.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *sensorService) GetSensors(ctx context.Context, dto dto.GetSensorsDTO) (sensors []*entity.Sensor, err error) {
	sensors, err = s.sensorRepo.GetAll(ctx, dto.FloorID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return sensors, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *sensorService) GetSensorsDetailed(ctx context.Context, dto dto.GetSensorsDetailedDTO) (sensors []*entity.SensorDetailed, err error) {
	sensors, err = s.sensorRepo.GetAllDetailed(ctx, dto.FloorID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return sensors, usecase.ErrNotFound
		}

		return
	}

	return
}

// TODO PUT update
func (s *sensorService) UpdateSensor(ctx context.Context, updateSensorDTO *dto.PatchUpdateSensorDTO) (err error) {
	err = s.sensorRepo.Update(ctx, updateSensorDTO)
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

func (s *sensorService) IsSensorSoftDeleted(ctx context.Context, sensorID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.sensorRepo.IsSensorSoftDeleted(ctx, sensorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *sensorService) SoftDeleteSensor(ctx context.Context, sensorID uuid.UUID) (err error) {
	err = s.sensorRepo.SoftDelete(ctx, sensorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *sensorService) RestoreSensor(ctx context.Context, sensorID uuid.UUID) (err error) {
	err = s.sensorRepo.Restore(ctx, sensorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}

		return
	}

	return
}
