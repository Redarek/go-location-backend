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
	"location-backend/internal/domain/entity"
	"location-backend/plugins/location"        // TODO удалить зависимость
	"location-backend/plugins/location/mapper" // TODO удалить зависимость
)

const publicFolderPath = "public"

type IMatrixService interface {
	CreateMatrix(ctx context.Context, points []*entity.Point, matrixPoints []*entity.MatrixPoint) (err error)

	SearchPoints(ctx context.Context, filter entity.SearchParameters) (points []*entity.Point, err error)

	DeletePoints(ctx context.Context, floorID uuid.UUID) (deletedCount int64, err error)
}

// TODO описать
type ILocationPlugin interface {
	// Execute() error
}

type MatrixUsecase struct {
	matrixService IMatrixService
	floorService  IFloorService
	wallService   IWallService
	sensorService ISensorService
	deviceService IDeviceService
}

// TODO придумать как компактно и универсально передавать параметры
func NewMatrixUsecase(
	matrixService IMatrixService,
	floorService IFloorService,
	wallService IWallService,
	sensorService ISensorService,
	deviceService IDeviceService,
) *MatrixUsecase {
	return &MatrixUsecase{
		matrixService: matrixService,
		floorService:  floorService,
		wallService:   wallService,
		sensorService: sensorService,
		deviceService: deviceService,
	}
}

func (u *MatrixUsecase) CreateMatrix(ctx context.Context, floorID uuid.UUID) (err error) {
	matrixInputData, err := u.getMatrixInputData(ctx, floorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get matrix input data")
		return
	}

	pointRows, matrixRows := location.CreateMatrix(floorID, matrixInputData)

	// Имя файла и путь
	fileName := uuid.New().String() + ".png"
	path := filepath.Join(publicFolderPath, fileName)
	// Создание матрицы в виде картинки
	err = createMatrixPNG(matrixInputData, pointRows, matrixRows, path)
	if err != nil {
		log.Error().Err(err).Msg("failed to create PNG from matrix")
		return
	}

	err = u.floorService.UpdateFloorHeatmap(ctx, floorID, fileName)
	if err != nil {
		log.Error().Err(err).Msg("failed to update heatmap")
		return
	}
	log.Debug().Msgf("heatmap saved as %v", path)

	// var matrix []*entity.Matrix
	// for _, point := range pointRows {
	// 	for _, matrixPoint := range matrixRows {
	// 		if matrixPoint.PointID == point.ID {
	// 			matrix = append(matrix, &entity.Matrix{
	// 				FloorID: floorID,
	// 				X:       matrixInputData.Floor.X,
	// 				Y:       matrixInputData.Floor.Y,
	// 			})
	// 		}
	// 	}
	// }
	_, err = u.matrixService.DeletePoints(ctx, floorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete points")
		return
	}
	// log.Debug().Msgf("heatmap saved as %v", path)

	err = u.matrixService.CreateMatrix(ctx, pointRows, matrixRows)
	if err != nil {
		log.Error().Err(err).Msg("failed insert matrix into database")
		return
	}

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
		Page:    0,
		Size:    0,
	}

	wallsDetailed, err := u.wallService.GetWallsDetailed(ctx, &getWallsDTO)
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
		Page:    0,
		Size:    0,
	}

	sensors, err := u.sensorService.GetSensors(ctx, &getSensorsDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			// TODO передавать данную ошибку наверх
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
	// var squareSize = floor.Scale / 1000 / floor.CellSizeMeter // ?
	matrixInputData = &location.InputData{
		Client: location.Client{
			TrSignalPower: 17,
			TrAntGain:     1,
			ZM:            1,
		},
		Walls:          wallMapper.EntitiesDomainToLocation(wallsDetailed, floor.Scale),
		Sensors:        sensorMapper.EntitiesDomainToLocation(sensors, floor.Scale),
		Floor:          *floorMapper.EntityDomainToLocation(floor),
		CellSizeMeters: floor.CellSizeMeter,
		MinX:           0,
		MinY:           0,
		// MaxX: int(float64(floor.WidthInPixels) * squareSize),  // !be careful here
		// MaxY: int(float64(floor.HeightInPixels) * squareSize), // !be careful here
		MaxX: location.PixelsToCells(floor.WidthInPixels, floor.Scale, floor.CellSizeMeter),  // !be careful here
		MaxY: location.PixelsToCells(floor.HeightInPixels, floor.Scale, floor.CellSizeMeter), // !be careful here
	}
	log.Debug().Msgf("matrix input data: %+v", matrixInputData)

	return
}

func normalize(value, min, max float64) float64 {
	return (value - min) / (max - min)
}

func generateColorAndOpacity(normalizedValue float64) color.Color {
	var r, g, b, a uint8
	// a = 153 // 0.6 * 255
	a = 255

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

func createMatrixContext(matrixInputData *location.InputData, pointRows []*entity.Point, matrixRows []*entity.MatrixPoint) (context *gg.Context) {
	// var squareSizePx = 1000 / matrixInputData.Floor.Scale * matrixInputData.Floor.CellSizeMeter // размер квадрата в пикселях
	var squareSizePx = location.MetersToPixels(matrixInputData.Floor.CellSizeMeter, matrixInputData.Floor.Scale)
	context = gg.NewContext(matrixInputData.Floor.WidthInPixels, matrixInputData.Floor.HeightInPixels)
	for _, point := range pointRows {
		var rssi float64 = -100
		for _, matrix := range matrixRows {
			if matrix.PointID == point.ID {
				rssi = matrix.RSSI24
				break
			}
		}

		if rssi != location.RSSI_INVISIBLE {
			normalizedValue := normalize(rssi, location.RSSI_INVISIBLE, 0) // -25
			clr := generateColorAndOpacity(normalizedValue)

			//! Возможно неверное преобразование
			// pointX := point.X * squareSizePx
			// pointY := point.Y * squareSizePx
			pointX := location.MetersToPixels(point.X, matrixInputData.Floor.Scale)
			pointY := location.MetersToPixels(point.Y, matrixInputData.Floor.Scale)

			context.DrawRectangle(pointX, pointY, squareSizePx, squareSizePx)
			context.SetColor(clr)
			context.Fill()
		}
	}

	return
}

func removePreviousPicture(matrixInputData *location.InputData) (err error) {
	if matrixInputData.Floor.Heatmap != nil {
		path := filepath.Join(publicFolderPath, *matrixInputData.Floor.Heatmap) // TODO fix
		err = os.Remove(path)
	}

	return
}

func saveMatrixToPNG(context *gg.Context, path string) (err error) {
	// ? проверка наличия папки public?
	if _, err = os.Stat(publicFolderPath); os.IsNotExist(err) {
		if err = os.Mkdir(publicFolderPath, os.ModePerm); err != nil {
			log.Error().Msg("failed to create directory")
			return
		}
	}

	return context.SavePNG(path)
}

func createMatrixPNG(matrixInputData *location.InputData, pointRows []*entity.Point, matrixRows []*entity.MatrixPoint, path string) (err error) {
	context := createMatrixContext(matrixInputData, pointRows, matrixRows)

	// Удаление предыдущей тепловой карты
	err = removePreviousPicture(matrixInputData)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn().Msgf("file %s does not exist", *matrixInputData.Floor.Heatmap)
		} else {
			log.Error().Msg("failed to remove previous heatmap")
			return
		}
	}
	log.Debug().Msgf("previous heatmap removed successfully")

	// Сохранение матрицы как PNG
	err = saveMatrixToPNG(context, path)
	if err != nil {
		log.Error().Msg("failed to save heatmap to PNG")
		return
	}

	return
}

func (u *MatrixUsecase) FindPoints(ctx context.Context, dto *dto.FindPointsDTO) (points []*entity.Point, err error) {
	devicesDetailed, err := u.deviceService.GetDevicesDetailedByMAC(ctx, dto.MAC, 0, 0)
	if err != nil {
		log.Error().Msg("failed to get devices detailed")
		return
	}

	// TODO delete?
	log.Debug().Msgf("detailed devices: %v", devicesDetailed)

	// TODO что тут должно быть?

	data := location.Data{
		MAC: dto.MAC,
		// TODO заполнить
	}
	location.GetParametersForFindXY(data)
	return
}
