package location

import (
	"fmt"
	"location-backend/internal/db"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// // import Sensor from "../models/sensors.model";
// // import Wall from "../models/walls.model";

// // import { generateMatrixRow, MatrixPoint } from "./matrix_row_generator";

func _getMatrix(matrixRowGenerator chan MatrixPoint, mapId uuid.UUID) ([]PointRow, []MatrixRow) {
	var pointRowsToInsert []PointRow
	var matrixRowsToInsert []MatrixRow
	var lastId = -1

	for row := range matrixRowGenerator {

		// (id, sensorId, x_m, y_m, rssi24, rssi5, rssi6, distance) := row;
		var id int = row.id
		var sensorId uuid.UUID = row.sensorId
		x_m, y_m := row.x_m, row.y_m
		rssi24, rssi5, rssi6 := row.rssi24, row.rssi5, row.rssi6
		var distance float64 = row.distance

		if lastId != id {
			// pointRowsToInsert.push({ id: id, map_id: mapId, x: x_m, y: y_m });
			pointRowsToInsert = append(pointRowsToInsert, PointRow{id: id, map_id: mapId, x: x_m, y: y_m})
			lastId = id
		}

		// matrixRowsToInsert.push({ point_id: id, sensor_id: sensorId, rssi24: rssi24, rssi5: rssi5, rssi6: rssi6, distance: distance });
		matrixRowsToInsert = append(matrixRowsToInsert, MatrixRow{point_id: id, sensor_id: sensorId, rssi24: rssi24, rssi5: rssi5, rssi6: rssi6, distance: distance})
	}

	return pointRowsToInsert, matrixRowsToInsert
}

func CreateMatrix(mapId uuid.UUID, inputData InputData) ([]PointRow, []MatrixRow) {
	log.Info().Msg(`Creating matrix for map_id = ${mapId}...`)
	var startTestTime time.Time = time.Now()

	// var pointRowsToInsert []PointRow
	// var matrixRowsToInsert []MatrixRow

	// try {
	//const startFillTime: number = performance.now();

	var matrixRowGenerator chan MatrixPoint = GenerateMatrixRow(inputData)
	//const [pointSize, matrixSize]: [number, number] = await _insertIntoMatrixAsync(pointRepository, matrixRepository, matrixRowGenerator, mapId);
	pointRowsToInsert, matrixRowsToInsert := _getMatrix(matrixRowGenerator, mapId)

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
// //     trSignalPower: number,
// //     trAntGain: number,
// //     zM: number
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
	client           Client
	walls            []Wall
	sensors          []db.Sensor
	cell_size_meters float64
	minX             int
	minY             int
	maxX             int
	maxY             int
}

type PointRow struct {
	id     int
	map_id uuid.UUID
	x      float64
	y      float64
}

type MatrixRow struct {
	point_id  int
	sensor_id uuid.UUID
	rssi24    float64
	rssi5     float64
	rssi6     float64
	distance  float64
}
