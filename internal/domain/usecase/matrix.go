package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	// "location-backend/internal/domain/dto"
	// "location-backend/internal/domain/entity"
)

// type MatrixService interface {
// 	CreateSensor(ctx context.Context, createDTO *dto.CreateSensorDTO) (sensorID uuid.UUID, err error)
// 	GetSensor(ctx context.Context, sensorID uuid.UUID) (sensor *entity.Sensor, err error)
// 	GetSensorByMAC(ctx context.Context, mac string) (sensor *entity.Sensor, err error)
// 	GetSensorDetailed(ctx context.Context, getDTO dto.GetSensorDetailedDTO) (sensorDetailed *entity.SensorDetailed, err error)
// 	GetSensors(ctx context.Context, getDTO dto.GetSensorsDTO) (sensors []*entity.Sensor, err error)
// 	GetSensorsDetailed(ctx context.Context, dto dto.GetSensorsDetailedDTO) (sensorsDetailed []*entity.SensorDetailed, err error)

// 	UpdateSensor(ctx context.Context, patchUpdateDTO *dto.PatchUpdateSensorDTO) (err error)

// 	IsSensorSoftDeleted(ctx context.Context, sensorID uuid.UUID) (isDeleted bool, err error)
// 	SoftDeleteSensor(ctx context.Context, sensorID uuid.UUID) (err error)
// 	RestoreSensor(ctx context.Context, sensorID uuid.UUID) (err error)
// }

type MatrixUsecase struct {
	floorService  FloorService
	wallService   WallService
	sensorService SensorService
}

func NewMatrixUsecase(
	floorService FloorService,
	wallService WallService,
	sensorService SensorService,
) *MatrixUsecase {
	return &MatrixUsecase{
		floorService:  floorService,
		wallService:   wallService,
		sensorService: sensorService,
	}
}

func (u *MatrixUsecase) CreateMatrix(ctx context.Context, floorID uuid.UUID) (err error) {
	floor, err := u.floorService.GetFloor(ctx, floorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}

		log.Error().Err(err).Msg("failed to get floor")
		return
	}

	wallsDetailed, err := u.wallService.GetWallDetailed(ctx, floor.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Debug().Msg("no any wall were found on the floor")
		}

		log.Error().Err(err).Msg("failed to get walls detailed")
		return
	}

	// TODO подумать над LIMIT OFFSET
	// sensors, err := u.sensorService.GetSensors(ctx, floor.ID)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to get sensors")
	// 	return ctx.SendStatus(fiber.StatusInternalServerError)
	// }

	return
}
