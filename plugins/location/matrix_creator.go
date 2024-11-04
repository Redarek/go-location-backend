package location

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/entity"
)

// // import { generateentity.MatrixPoint, MatrixPoint } from "./matrix_row_generator";

func _getMatrix(matrixRowGenerator chan MatrixPoint, floorID uuid.UUID) ([]*entity.Point, []*entity.MatrixPoint) {
	var pointRowsToInsert []*entity.Point
	var matrixRowsToInsert []*entity.MatrixPoint
	var lastID = uuid.UUID{}

	for row := range matrixRowGenerator {
		// (ID, SensorID, xM, yM, RSSI24, RSSI5, RSSI6, Distance) := row;
		var id uuid.UUID = row.id
		var sensorID uuid.UUID = row.sensorID
		xM, yM := row.xM, row.yM
		rssi24, rssi5, rssi6 := row.rssi24, row.rssi5, row.rssi6
		var distance float64 = row.distance

		if lastID != id {
			// pointRowsToInsert.push({ ID: ID, MapID: MapID, X: xM, Y: yM });
			pointRowsToInsert = append(pointRowsToInsert, &entity.Point{ID: id, FloorID: floorID, X: xM, Y: yM})
			lastID = id
		}

		// matrixRowsToInsert.push({ PointID: ID, SensorID: SensorID, RSSI24: RSSI24, RSSI5: RSSI5, RSSI6: RSSI6, Distance: Distance });
		matrixRowsToInsert = append(matrixRowsToInsert, &entity.MatrixPoint{PointID: id, SensorID: sensorID, RSSI24: rssi24, RSSI5: rssi5, RSSI6: rssi6, Distance: distance})
	}

	return pointRowsToInsert, matrixRowsToInsert
}

func CreateMatrix(floorID uuid.UUID, inputData *InputData) ([]*entity.Point, []*entity.MatrixPoint) {
	log.Info().Msg(fmt.Sprintf("Creating matrix for FloorID = %s...", floorID))
	var startTestTime time.Time = time.Now()

	// var pointRowsToInsert []entity.Point
	// var matrixRowsToInsert []entity.MatrixPoint

	// try {
	//const startFillTime: number = performance.now();

	var matrixRowGenerator chan MatrixPoint = GenerateMatrixRow(inputData)
	//const [pointSize, matrixSize]: [number, number] = await _insertIntoMatrixAsync(pointRepository, matrixRepository, matrixRowGenerator, MapID);
	pointRowsToInsert, matrixRowsToInsert := _getMatrix(matrixRowGenerator, floorID)

	//logger.info(`Created ${pointSize} points (${matrixSize} matrix points) in ${((performance.now() - startTestTime) / 1000).toFixed(2)} sec `
	//    + `(Del: ${deleteTime} sec, Fill: ${((performance.now() - startFillTime) / 1000).toFixed(2)} sec)`);
	log.Info().Msg(fmt.Sprintf("Created %d points (%d matrix points) in %v sec ",
		len(pointRowsToInsert), len(matrixRowsToInsert), time.Since(startTestTime).Seconds()))
	// }
	// catch (err) {
	//     logger.fatal(err);
	//     //return [{ 'Error': err, 'status': 'error' }, 500];
	//     throw err;
	// }

	//return [{ 'status': 'ok' }, 200];
	return pointRowsToInsert, matrixRowsToInsert
}

type InputData struct {
	Client         Client
	Walls          []*Wall
	Sensors        []*Sensor
	Floor          Floor
	CellSizeMeters float64
	MinX           int
	MinY           int
	MaxX           int
	MaxY           int
}

// type entity.Point struct {
// 	ID      int       `json:"ID"`
// 	FloorID uuid.UUID `json:"floorId"`
// 	X       float64   `json:"X"`
// 	Y       float64   `json:"Y"`
// }

// type entity.MatrixPoint struct {
// 	PointID  int       `json:"pointId"`
// 	SensorID uuid.UUID `json:"sensorId"`
// 	RSSI24   float64   `json:"RSSI24"`
// 	RSSI5    float64   `json:"RSSI5"`
// 	RSSI6    float64   `json:"RSSI6"`
// 	Distance float64   `json:"Distance"`
// }
