// CTRL + A, затем раскомментировать

package location

import (
	"encoding/json"
	"errors"
	"location-backend/internal/db"
	. "math"

	"github.com/google/uuid"

	// "location-backend/internal/logger"
	"github.com/rs/zerolog/log"
)

// import logger from "../../../logger"
// import Sensor from "../models/sensors.model"
// import Wall from "../models/walls.model"

// import ("consts");

// import { getHorizontalAzimuthDeg, getVerticalAzimuthDeg, getWallPathLengthThrough } from "./math_functionality";

// import type { Client, InputData } from "./matrix_creator"
//import type { XYcoordinate, XYZcoordinate } from "./math_functionality";

type MatrixPoint struct {
	id       int
	sensorId uuid.UUID
	x        int
	y        int
	x_m      float64
	y_m      float64
	rssi24   float64
	rssi5    float64
	rssi6    float64
	distance float64
}

type Client struct {
	trSignalPower int
	trAntGain     int
	zM            float64
}

// type Generator chan MatrixPoint

func GenerateMatrixRow(inputData InputData) chan MatrixPoint {
	// var client, walls, sensors, cell_size_meters, minX, minY, maxX, maxY int = inputData;
	var client Client = inputData.client
	var walls []Wall = inputData.walls
	var sensors []db.Sensor = inputData.sensors
	var minX int = inputData.minX
	var minY int = inputData.minY
	var maxX int = inputData.maxX
	var maxY int = inputData.maxY
	var cell_size_meters float64 = inputData.cell_size_meters

	var i int = 0

	ch := make(chan MatrixPoint)
	go func() {

		for y := minY; y < maxY+1; y++ {
			for x := minX; x < maxX+1; x++ {
				i++
				var matrixWithPoint MatrixPoint = MatrixPoint{
					id: i,
					// sensorId: -1,
					x:        x,
					y:        y,
					x_m:      float64(x) * cell_size_meters,
					y_m:      float64(y) * cell_size_meters,
					rssi24:   RSSI_INVISIBLE,
					rssi5:    RSSI_INVISIBLE,
					rssi6:    RSSI_INVISIBLE,
					distance: DISTANCE_INVISIBLE,
				}

				for _, sensor := range sensors {
					var distance float64 = _getDistance(x, y, client, sensor, cell_size_meters)
					var freeSpaceRSSI24, freeSpaceRSSI5, freeSpaceRSSI6 = _getFreeSpaceRSSI(x, y, client, sensor, distance)
					var wallsLoss24 float64 = 0
					var wallsLoss5 float64 = 0
					var wallsLoss6 float64 = 0

					if CALCULATE_WALLS {
						if freeSpaceRSSI24 >= RSII_CUTOFF || freeSpaceRSSI5 >= RSII_CUTOFF || freeSpaceRSSI6 >= RSII_CUTOFF {
							wallsLoss24, wallsLoss5, wallsLoss6 = _getWallsAttenuation(x, y, walls, sensor, client, cell_size_meters)
						}
					}

					var tempRSSI24 float64 = freeSpaceRSSI24 + wallsLoss24 + CORRECTION_COEFFICIENT_24
					// var rssi_24 float64 = (tempRSSI24 >= RSII_CUTOFF) ? Number(tempRSSI24.toFixed(1)) : RSSI_INVISIBLE;
					var rssi_24 float64
					var rssi_5 float64
					var rssi_6 float64
					if tempRSSI24 >= RSII_CUTOFF {
						rssi_24 = Round(tempRSSI24*10) / 10 // округление до 1 знака
					} else {
						rssi_24 = RSSI_INVISIBLE
					}

					var tempRSSI5 float64 = freeSpaceRSSI5 + wallsLoss5 + CORRECTION_COEFFICIENT_5
					// var rssi_5 float64 = (tempRSSI5 >= RSII_CUTOFF) ? Number(tempRSSI5.toFixed(1)) : RSSI_INVISIBLE;
					if tempRSSI5 >= RSII_CUTOFF {
						rssi_5 = Round(tempRSSI5*10) / 10
					} else {
						rssi_5 = RSSI_INVISIBLE
					}

					var tempRSSI6 float64 = freeSpaceRSSI6 + wallsLoss6 + CORRECTION_COEFFICIENT_6
					// var rssi_6 float64 = (tempRSSI6 >= RSII_CUTOFF) ? Number(tempRSSI6.toFixed(1)) : RSSI_INVISIBLE;
					if tempRSSI6 >= RSII_CUTOFF {
						rssi_6 = Round(tempRSSI6*10) / 10
					} else {
						rssi_6 = RSSI_INVISIBLE
					}

					distance = Round(distance*10) / 10

					matrixWithPoint.sensorId = sensor.ID
					matrixWithPoint.rssi24 = rssi_24
					matrixWithPoint.rssi5 = rssi_5
					matrixWithPoint.rssi6 = rssi_6
					matrixWithPoint.distance = distance

					ch <- matrixWithPoint
				}
			}
		}
		close(ch) // Close the channel when done
	}()
	return ch
}

/**
 * Returns the negative numbers of total walls attenuation for 2.4, 5 and 6 HHz bands.
 * @param clientX
 * @param clientY
 * @param walls
 * @param sensor
 * @param client
 * @returns
 */
func _getWallsAttenuation(clientX int, clientY int, walls []Wall, sensor db.Sensor, client Client, cell_size_meters float64) (float64, float64, float64) {
	var loss24 float64 = 0
	var loss5 float64 = 0
	var loss6 float64 = 0

	for _, wall := range walls {
		var wall_path_length_through float64 = getWallPathLengthThrough(XYZcoordinate{x: float64(clientX), y: float64(clientY), z: client.zM},
			XYZcoordinate{x: sensor.X, y: sensor.Y, z: sensor.Z},
			XYZcoordinate{x: float64(wall.X1), y: float64(wall.Y1), z: 0},
			XYZcoordinate{x: float64(wall.X2), y: float64(wall.Y2), z: 0},
			wall.Thickness,
			cell_size_meters)

		if wall_path_length_through != 0 {
			var pathDivideThickness float64 = wall_path_length_through / wall.Thickness
			loss24 -= wall.Attenuation24 * pathDivideThickness
			loss5 -= wall.Attenuation5 * pathDivideThickness
			loss6 -= wall.Attenuation6 * pathDivideThickness

			if loss24 <= RSII_CUTOFF && loss5 <= RSII_CUTOFF && loss6 <= RSII_CUTOFF {
				break
			}
		}
	}

	return loss24, loss5, loss6
}

/**
 * Returns the distance in meters between client and sensor.
 * @param clientX Client x coordinate.
 * @param clientY Client y coordinate.
 * @param client Client`s parameters.
 * @param sensor Sensor.
 * @returns Distance between client and sensor in meters.
 */
func _getDistance(clientX int, clientY int, client Client, sensor db.Sensor, cell_size_meters float64) float64 {
	return Magnitude(Vector{(float64(clientX) - sensor.X) * cell_size_meters, (float64(clientY) - sensor.Y) * cell_size_meters, (client.zM - sensor.Z)})
}

/**
 * Returns the free space pass loss in dB.
 * @param frequency Transmission frequency in GHz.
 * @param attenuation_factor Attenuation factor.
 * @param penetration_factor Penetration factor.
 * @param distance Transmission distance.
 * @returns Free space pass loss in dB.
 */
func _getFSPL(frequency int, attenuation_factor float64, penetration_factor float64, distance float64) float64 {
	if distance < 1 {
		distance = 1
	}
	return 20*Log10(float64(frequency)) + 10*attenuation_factor*Log10(distance) + penetration_factor - 24
}

func _approximateAzimuth(azimuth float64, delta float64) (int, error) {
	if delta == 0 {
		return 0, errors.New("delta cannot be zero")
	}

	return int(Floor((azimuth+Floor(delta/2))/delta)*delta) % 360, nil
}

/**
 * Returns the RSSI for 2.4, 5 and 6 HHz bands in a free space.
 * @param clientX Client x coordinate.
 * @param clientY Client y coordinate.
 * @param client Client`s parameters.
 * @param sensor Sensor.
 * @param distance Distance between client and sensors in meters.
 * @returns Tuple of RSSI for 2.4, 5 and 6 HHz bands.
 */
func _getFreeSpaceRSSI(clientX int, clientY int, client Client, sensor db.Sensor, distance float64) (float64, float64, float64) {
	var fspl24 float64 = _getFSPL(FREQUENCY24, ATTENUATION_FACTOR24, PENETRATION_FACTOR24, distance)
	var fspl5 float64 = _getFSPL(FREQUENCY5, ATTENUATION_FACTOR5, PENETRATION_FACTOR5, distance)
	var fspl6 float64 = _getFSPL(FREQUENCY6, ATTENUATION_FACTOR6, PENETRATION_FACTOR6, distance)

	var freeSpaceRSSI24 float64 = float64(client.trSignalPower) + float64(client.trAntGain) - fspl24 + sensor.CorrectionFactor24
	var freeSpaceRSSI5 float64 = float64(client.trSignalPower) + float64(client.trAntGain) - fspl5 + sensor.CorrectionFactor5
	var freeSpaceRSSI6 float64 = float64(client.trSignalPower) + float64(client.trAntGain) - fspl6 + sensor.CorrectionFactor6

	var ant_gain float64 = 2

	var diagram db.Diagram
	err := json.Unmarshal(sensor.Diagram, &diagram)
	if err != nil {
		var delta int = 0
		if _, ok := diagram.Degree["10"]; ok {
			delta = 10
		}
		if _, ok := diagram.Degree["15"]; ok {
			delta = 15
		}

		if delta == 0 {
			log.Warn().Msg(`The radiation diagram can have only the step of 10 or 15 degrees.
            Check that the radiation diagram for sensor with id = ${sensor.id} is filled out correctly.
            By default, the antenna gain of ${sensor.rx_ant_gain} will be used for all directions.`)
		} else {
			hor_azimuth, error := _approximateAzimuth(
				float64(getHorizontalAzimuthDeg(
					XYZcoordinate{x: sensor.X, y: sensor.Y, z: sensor.Z},
					XYZcoordinate{float64(clientX), float64(clientY), client.zM},
					0)),
				float64(delta))
			if error != nil {
				goto errorHandling
			}

			vert_azimuth, err := _approximateAzimuth(
				float64(getVerticalAzimuthDeg(
					XYZcoordinate{x: sensor.X, y: sensor.Y, z: sensor.Z},
					XYZcoordinate{float64(clientX), float64(clientY), client.zM},
					0)),
				float64(delta))
			if err != nil {
				goto errorHandling
			}

			ant_gain = (diagram.Degree[string(hor_azimuth)].HorGain + diagram.Degree[string(vert_azimuth)].VertGain) / 2 // окр до десятых

		}
	} else {
		freeSpaceRSSI24 += sensor.RxAntGain
		freeSpaceRSSI5 += sensor.RxAntGain
		freeSpaceRSSI6 += sensor.RxAntGain
	}

	return freeSpaceRSSI24, freeSpaceRSSI5, freeSpaceRSSI6

errorHandling:
	ant_gain = sensor.RxAntGain
	freeSpaceRSSI24 += ant_gain
	freeSpaceRSSI5 += ant_gain
	freeSpaceRSSI6 += ant_gain
	return freeSpaceRSSI24, freeSpaceRSSI5, freeSpaceRSSI6
}
