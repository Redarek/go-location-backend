package usecase

import (
	"context"
	"errors"
	"image/color"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
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
	publicFolderPath := "public"

	matrixInputData, err := u.getMatrixInputData(ctx, floorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get matrix input data")
		return
	}

	pointRows, matrixRows := location.CreateMatrix(floorID, matrixInputData)

	var squareSizePx = 1000 / matrixInputData.Floor.Scale * matrixInputData.Floor.CellSizeMeter // размер квадрата в пикселях
	dc := gg.NewContext(matrixInputData.Floor.WidthInPixels, matrixInputData.Floor.HeightInPixels)
	for _, point := range pointRows {
		var rssi float64 = -100
		for _, matrix := range matrixRows {
			if matrix.PointID == point.ID {
				rssi = matrix.RSSI24
				break
			}
		}

		if rssi != location.RSSI_INVISIBLE {
			normalizedValue := normalize(rssi, location.RSSI_INVISIBLE, -25)
			clr := generateColorAndOpacity(normalizedValue)

			// pointY := point.Y * matrixInputData.Floor.Scale / 1000
			//! Возможно неверное преобразование
			// pointX := point.X * 1000 / matrixInputData.Floor.Scale
			// pointY := point.Y * 1000 / matrixInputData.Floor.Scale
			pointX := point.X * squareSizePx
			pointY := point.Y * squareSizePx

			dc.DrawRectangle(pointX, pointY, squareSizePx, squareSizePx)
			dc.SetColor(clr)
			dc.Fill()
		}
	}

	// Удаление предыдущей тепловой карты
	if matrixInputData.Floor.Heatmap != nil {
		path := filepath.Join(publicFolderPath, *matrixInputData.Floor.Heatmap) // TODO fix
		err = os.Remove(path)
		if err != nil {
			if os.IsNotExist(err) {
				log.Warn().Msgf("file %s does not exist", *matrixInputData.Floor.Heatmap)
			} else {
				log.Error().Err(err).Msg("failed to remove previous heatmap")
				return
			}
		}
		log.Debug().Msgf("previous heatmap removed successfully")
	}

	fileName := uuid.New().String() + ".png"
	outputPath := filepath.Join(publicFolderPath, fileName)

	// ? проверка наличия папки public?
	if _, err = os.Stat(publicFolderPath); os.IsNotExist(err) {
		if err = os.Mkdir(publicFolderPath, os.ModePerm); err != nil {
			log.Error().Err(err).Msg("failed to create directory")
			return
		}
	}

	err = dc.SavePNG(outputPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to save heatmap")
		return
	}

	err = u.floorService.UpdateFloorHeatmap(ctx, floorID, fileName)
	if err != nil {
		log.Error().Err(err).Msg("failed to update heatmap")
		return
	}
	log.Debug().Msgf("heatmap saved as %v", outputPath)

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
		} else {
			log.Error().Err(err).Msg("failed to get walls detailed")
			return
		}
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
	floorMapper := mapper.GetFloorMapper()

	// TODO убрать избыточные параметры
	var squareSize = floor.Scale / 1000 / floor.CellSizeMeter
	matrixInputData = &location.InputData{
		Client: location.Client{
			TrSignalPower: 17,
			TrAntGain:     1,
			ZM:            1,
		},
		Walls:          wallMapper.EntitiesDomainToLocation(wallsDetailed),
		Sensors:        sensorMapper.EntitiesDomainToLocation(sensors),
		Floor:          *floorMapper.EntityDomainToLocation(floor),
		CellSizeMeters: floor.CellSizeMeter,
		MinX:           0,
		MinY:           0,
		// MaxX:           int(float64(floor.WidthInPixels) * floor.Scale * 1000 / floor.CellSizeMeter),  // !be careful here
		// MaxY:           int(float64(floor.HeightInPixels) * floor.Scale * 1000 / floor.CellSizeMeter), // !be careful here
		MaxX: int(float64(floor.WidthInPixels) * squareSize),  // !be careful here
		MaxY: int(float64(floor.HeightInPixels) * squareSize), // !be careful here
	}
	log.Debug().Msgf("matrix input data: %+v", matrixInputData)

	return
}

func normalize(value, min, max float64) float64 {
	return (value - min) / (max - min)
}

func generateColorAndOpacity(normalizedValue float64) color.Color {
	var r, g, b, a uint8
	a = 153 // 0.6 * 255

	if normalizedValue == 0 {
		r, g, b, a = 0, 0, 0, 0
	} else if normalizedValue <= 0.33 {
		r = 0
		g = uint8((255 * normalizedValue) / 0.33)
		b = 255
	} else if normalizedValue <= 0.66 {
		r = uint8((255 * (normalizedValue - 0.33)) / 0.33)
		g = 255
		b = 0
	} else {
		r = 255
		g = uint8((255 * (1 - normalizedValue)) / 0.34)
		b = 0
	}

	return color.NRGBA{R: r, G: g, B: b, A: a}
}
