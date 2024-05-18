package location

import (
	"fmt"
	"location-backend/internal/db"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// // import Sensor from "../models/Sensors.model";
// // import Wall from "../models/Walls.model";

// // import { generateMatrixRow, MatrixPoint } from "./matrix_row_generator";

func _getMatrix(matrixRowGenerator chan MatrixPoint, mapID uuid.UUID) ([]PointRow, []MatrixRow) {
	var pointRowsToInsert []PointRow
	var matrixRowsToInsert []MatrixRow
	var lastID = -1

	for row := range matrixRowGenerator {

		// (id, sensorID, xM, yM, rssi24, rssi5, rssi6, distance) := row;
		var id int = row.id
		var sensorID uuid.UUID = row.sensorID
		xM, yM := row.xM, row.yM
		rssi24, rssi5, rssi6 := row.rssi24, row.rssi5, row.rssi6
		var distance float64 = row.distance

		if lastID != id {
			// pointRowsToInsert.push({ id: id, mapID: mapID, x: xM, y: yM });
			pointRowsToInsert = append(pointRowsToInsert, PointRow{id: id, mapID: mapID, x: xM, y: yM})
			lastID = id
		}

		// matrixRowsToInsert.push({ pointID: id, sensorID: sensorID, rssi24: rssi24, rssi5: rssi5, rssi6: rssi6, distance: distance });
		matrixRowsToInsert = append(matrixRowsToInsert, MatrixRow{pointID: id, sensorID: sensorID, rssi24: rssi24, rssi5: rssi5, rssi6: rssi6, distance: distance})
	}

	return pointRowsToInsert, matrixRowsToInsert
}

func CreateMatrix(mapID uuid.UUID, inputData InputData) ([]PointRow, []MatrixRow) {
	log.Info().Msg(fmt.Sprintf("Creating matrix for map_id = %s...", mapID))
	var startTestTime time.Time = time.Now()

	// var pointRowsToInsert []PointRow
	// var matrixRowsToInsert []MatrixRow

	// try {
	//const startFillTime: number = performance.now();

	var matrixRowGenerator chan MatrixPoint = GenerateMatrixRow(inputData)
	//const [pointSize, matrixSize]: [number, number] = await _insertIntoMatrixAsync(pointRepository, matrixRepository, matrixRowGenerator, mapID);
	pointRowsToInsert, matrixRowsToInsert := _getMatrix(matrixRowGenerator, mapID)

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

// // type Client = {
// //     TrSignalPower: number,
// //     TrAntGain: number,
// //     ZM: number
// // }

type Wall struct {
	ID            uuid.UUID
	X1            int
	Y1            int
	X2            int
	Y2            int
	Thickness     float64
	Attenuation24 float64
	Attenuation5  float64
	Attenuation6  float64
}

type InputData struct {
	Client         Client
	Walls          []Wall
	Sensors        []*db.Sensor
	CellSizeMeters float64
	MinX           int
	MinY           int
	MaxX           int
	MaxY           int
}

type PointRow struct {
	id    int
	mapID uuid.UUID
	x     float64
	y     float64
}

type MatrixRow struct {
	pointID  int
	sensorID uuid.UUID
	rssi24   float64
	rssi5    float64
	rssi6    float64
	distance float64
}
