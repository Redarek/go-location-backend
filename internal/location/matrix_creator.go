package location

import (
	"fmt"
	"location-backend/internal/db/model"
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

		// (ID, SensorID, xM, yM, RSSI24, RSSI5, RSSI6, Distance) := row;
		var id int = row.id
		var sensorID uuid.UUID = row.sensorID
		xM, yM := row.xM, row.yM
		rssi24, rssi5, rssi6 := row.rssi24, row.rssi5, row.rssi6
		var distance float64 = row.distance

		if lastID != id {
			// pointRowsToInsert.push({ ID: ID, MapID: MapID, X: xM, Y: yM });
			pointRowsToInsert = append(pointRowsToInsert, PointRow{ID: id, MapID: mapID, X: xM, Y: yM})
			lastID = id
		}

		// matrixRowsToInsert.push({ PointID: ID, SensorID: SensorID, RSSI24: RSSI24, RSSI5: RSSI5, RSSI6: RSSI6, Distance: Distance });
		matrixRowsToInsert = append(matrixRowsToInsert, MatrixRow{PointID: id, SensorID: sensorID, RSSI24: rssi24, RSSI5: rssi5, RSSI6: rssi6, Distance: distance})
	}

	return pointRowsToInsert, matrixRowsToInsert
}

func CreateMatrix(mapID uuid.UUID, inputData InputData) ([]PointRow, []MatrixRow) {
	log.Info().Msg(fmt.Sprintf("Creating matrix for MapID = %s...", mapID))
	var startTestTime time.Time = time.Now()

	// var pointRowsToInsert []PointRow
	// var matrixRowsToInsert []MatrixRow

	// try {
	//const startFillTime: number = performance.now();

	var matrixRowGenerator chan MatrixPoint = GenerateMatrixRow(inputData)
	//const [pointSize, matrixSize]: [number, number] = await _insertIntoMatrixAsync(pointRepository, matrixRepository, matrixRowGenerator, MapID);
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
	Sensors        []*model.Sensor
	CellSizeMeters float64
	MinX           int
	MinY           int
	MaxX           int
	MaxY           int
}

type PointRow struct {
	ID    int       `json:"ID"`
	MapID uuid.UUID `json:"mapId"`
	X     float64   `json:"X"`
	Y     float64   `json:"Y"`
}

type MatrixRow struct {
	PointID  int       `json:"pointId"`
	SensorID uuid.UUID `json:"sensorId"`
	RSSI24   float64   `json:"RSSI24"`
	RSSI5    float64   `json:"RSSI5"`
	RSSI6    float64   `json:"RSSI6"`
	Distance float64   `json:"Distance"`
}
