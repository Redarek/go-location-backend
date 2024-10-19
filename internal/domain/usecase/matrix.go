package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/plugins/location"
	"location-backend/plugins/location/mapper"
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

// TODO описать
type ILocationPlugin interface {
	// Execute() error
}

type MatrixUsecase struct {
	floorService  IFloorService
	wallService   IWallService
	sensorService ISensorService
}

func NewMatrixUsecase(
	floorService IFloorService,
	wallService IWallService,
	sensorService ISensorService,
) *MatrixUsecase {
	return &MatrixUsecase{
		floorService:  floorService,
		wallService:   wallService,
		sensorService: sensorService,
	}
}

func (u *MatrixUsecase) CreateMatrix(ctx context.Context, floorID uuid.UUID) (err error) {
	matrixInputData, err := u.getMatrixInputData(ctx, floorID)
	log.Debug().Msgf("matrix input data: %+v", matrixInputData) // TODO remove

	return
}

func (u *MatrixUsecase) getMatrixInputData(ctx context.Context, floorID uuid.UUID) (matrixInputData *location.InputData, err error) {
	floor, err := u.floorService.GetFloor(ctx, floorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return
		}

		log.Error().Err(err).Msg("failed to get floor")
		return
	}

	getWallsDTO := dto.GetWallsDTO{
		FloorID: floorID,
		Limit:   0,
		Offset:  0,
	}

	wallsDetailed, err := u.wallService.GetWallsDetailed(ctx, getWallsDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Debug().Msg("no any wall were found on the floor")
		}

		log.Error().Err(err).Msg("failed to get walls detailed")
		return
	}

	getSensorsDTO := dto.GetSensorsDTO{
		FloorID: floorID,
		Limit:   0,
		Offset:  0,
	}

	sensors, err := u.sensorService.GetSensors(ctx, getSensorsDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Error().Err(err).Msg("no any sensor were found on the floor")
			return
		}

		log.Error().Err(err).Msg("failed to get sensors")
		return
	}

	wallMapper := mapper.GetWallMapper()
	sensorMapper := mapper.GetSensorMapper()

	matrixInputData = &location.InputData{
		Client: location.Client{
			TrSignalPower: 17,
			TrAntGain:     1,
			ZM:            1,
		},
		Walls:          wallMapper.EntitiesDomainToLocation(wallsDetailed),
		Sensors:        sensorMapper.EntitiesDomainToLocation(sensors),
		CellSizeMeters: floor.CellSizeMeter,
		MinX:           0,
		MinY:           0,
		MaxX:           int(float64((float64(floor.WidthInPixels)*floor.Scale)/1000) / floor.CellSizeMeter), // !be careful here
		MaxY:           int(float64((float64(floor.WidthInPixels)*floor.Scale)/1000) / floor.CellSizeMeter), // !be careful here
	}
	log.Debug().Msgf("matrix input data: %+v", matrixInputData)

	return
}
