package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type ISensorService interface {
	CreateSensor(ctx context.Context, createDTO *dto.CreateSensorDTO) (sensorID uuid.UUID, err error)
	GetSensor(ctx context.Context, sensorID uuid.UUID) (sensor *entity.Sensor, err error)
	GetSensorByMAC(ctx context.Context, mac string) (sensor *entity.Sensor, err error)
	GetSensorDetailed(ctx context.Context, getDTO *dto.GetSensorDetailedDTO) (sensorDetailed *entity.SensorDetailed, err error)
	GetSensors(ctx context.Context, getDTO *dto.GetSensorsDTO) (sensors []*entity.Sensor, err error)
	GetSensorsDetailed(ctx context.Context, dto *dto.GetSensorsDetailedDTO) (sensorsDetailed []*entity.SensorDetailed, err error)

	UpdateSensor(ctx context.Context, patchUpdateDTO *dto.PatchUpdateSensorDTO) (err error)

	IsSensorSoftDeleted(ctx context.Context, sensorID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteSensor(ctx context.Context, sensorID uuid.UUID) (err error)
	RestoreSensor(ctx context.Context, sensorID uuid.UUID) (err error)
}

type SensorUsecase struct {
	sensorService     ISensorService
	sensorTypeService ISensorTypeService
	floorService      IFloorService
}

func NewSensorUsecase(
	sensorService ISensorService,
	sensorTypeService ISensorTypeService,
	floorService IFloorService,
) *SensorUsecase {
	return &SensorUsecase{
		sensorService:     sensorService,
		sensorTypeService: sensorTypeService,
		floorService:      floorService,
	}
}

func (u *SensorUsecase) CreateSensor(ctx context.Context, createDTO *dto.CreateSensorDTO) (sensorID uuid.UUID, err error) {
	_, err = u.floorService.GetFloor(ctx, createDTO.FloorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Err(err).Msg("failed to create sensor: the floor with provided floor ID does not exist")
			return sensorID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check floor existing")
		return
	}

	_, err = u.sensorTypeService.GetSensorType(ctx, createDTO.SensorTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Err(err).Msg("failed to create sensor: the sensor type with provided sensor type ID does not exist")
			return sensorID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check sensor type existing")
		return
	}

	_, err = u.sensorService.GetSensorByMAC(ctx, createDTO.MAC)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			log.Error().Err(err).Msg("failed to check sensor MAC existing")
			return
		}
	} else {
		log.Info().Err(err).Msg("failed to create sensor: the sensor with provided MAC already exists")
		return sensorID, ErrAlreadyExists
	}

	sensorID, err = u.sensorService.CreateSensor(ctx, createDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create sensor")
		return
	}

	log.Info().Msgf("sensor %v successfully created", sensorID)
	return
}

func (u *SensorUsecase) GetSensor(ctx context.Context, sensorID uuid.UUID) (sensor *entity.Sensor, err error) {
	sensor, err = u.sensorService.GetSensor(ctx, sensorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to get sensor")
		return
	}

	return
}

func (u *SensorUsecase) GetSensorDetailed(ctx context.Context, getDTO *dto.GetSensorDetailedDTO) (sensorDetailed *entity.SensorDetailed, err error) {
	sensorDetailed, err = u.sensorService.GetSensorDetailed(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to get sensor detailed")
		return
	}

	return
}

func (u *SensorUsecase) GetSensors(ctx context.Context, getDTO *dto.GetSensorsDTO) (sensors []*entity.Sensor, err error) {
	sensors, err = u.sensorService.GetSensors(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get sensors")
			return
		}
	}

	return
}

func (u *SensorUsecase) GetSensorsDetailed(ctx context.Context, getDTO *dto.GetSensorsDetailedDTO) (sensorsDetailed []*entity.SensorDetailed, err error) {
	sensorsDetailed, err = u.sensorService.GetSensorsDetailed(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get sensors detailed")
			return
		}
	}

	return
}

func (u *SensorUsecase) PatchUpdateSensor(ctx context.Context, patchUpdateDTO *dto.PatchUpdateSensorDTO) (err error) {
	_, err = u.sensorService.GetSensor(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Error().Err(err).Msg("failed to check sensor existing")
			return ErrNotFound
		}
	}

	err = u.sensorService.UpdateSensor(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Err(err).Msg("sensor was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update sensor")
		return
	}

	return
}

func (u *SensorUsecase) SoftDeleteSensor(ctx context.Context, sensorID uuid.UUID) (err error) {
	isDeleted, err := u.sensorService.IsSensorSoftDeleted(ctx, sensorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if sensor is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.sensorService.SoftDeleteSensor(ctx, sensorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete sensor")
		return
	}

	return
}

func (u *SensorUsecase) RestoreSensor(ctx context.Context, sensorID uuid.UUID) (err error) {
	isDeleted, err := u.sensorService.IsSensorSoftDeleted(ctx, sensorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if sensor is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.sensorService.RestoreSensor(ctx, sensorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore sensor")
		return
	}

	return
}
