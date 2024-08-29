// CTRL + A, затем раскомментировать

package location

import (
	"encoding/json"
	"errors"
	"location-backend/internal/db"
	"location-backend/internal/db/model"
	. "math"

	"github.com/google/uuid"

	// "location-backend/internal/logger"
	"github.com/rs/zerolog/log"
)

// import logger from "../../../logger"
// import Sensor from "../models/Sensors.model"
// import Wall from "../models/Walls.model"

// import ("consts");

// import { getHorizontalAzimuthDeg, getVerticalAzimuthDeg, getWallPathLengthThrough } from "./math_functionality";

// import type { Client, InputData } from "./matrix_creator"
//import type { XYcoordinate, XYZcoordinate } from "./math_functionality";

type MatrixPoint struct {
	id       int
	sensorID uuid.UUID
	x        int
	y        int
	xM       float64
	yM       float64
	rssi24   float64
	rssi5    float64
	rssi6    float64
	distance float64
}

type Client struct {
	TrSignalPower int
	TrAntGain     int
	ZM            float64
}

// type Generator chan MatrixPoint

func GenerateMatrixRow(inputData InputData) chan MatrixPoint {
	// var Client, Walls, Sensors, CellSizeMeters, MinX, MinY, MaxX, MaxY int = inputData;
	var client Client = inputData.Client
	var walls []Wall = inputData.Walls
	var sensors []*model.Sensor = inputData.Sensors
	var minX int = inputData.MinX
	var minY int = inputData.MinY
	var maxX int = inputData.MaxX
	var maxY int = inputData.MaxY
	var cellSizeMeters float64 = inputData.CellSizeMeters

	var i int = 0

	ch := make(chan MatrixPoint)
	go func() {

		for y := minY; y < maxY+1; y++ {
			for x := minX; x < maxX+1; x++ {
				i++
				var matrixWithPoint MatrixPoint = MatrixPoint{
					id: i,
					// SensorID: -1,
					x:        x,
					y:        y,
					xM:       float64(x) * cellSizeMeters,
					yM:       float64(y) * cellSizeMeters,
					rssi24:   RSSI_INVISIBLE,
					rssi5:    RSSI_INVISIBLE,
					rssi6:    RSSI_INVISIBLE,
					distance: DISTANCE_INVISIBLE,
				}

				for _, sensor := range sensors {
					var distance float64 = _getDistance(matrixWithPoint.xM, matrixWithPoint.yM, client, *sensor, cellSizeMeters)
					var freeSpaceRSSI24, freeSpaceRSSI5, freeSpaceRSSI6 = _getFreeSpaceRSSI(matrixWithPoint.xM, matrixWithPoint.yM, client, *sensor, distance)
					var wallsLoss24 float64 = 0
					var wallsLoss5 float64 = 0
					var wallsLoss6 float64 = 0

					if CALCULATE_WALLS {
						if freeSpaceRSSI24 >= RSII_CUTOFF || freeSpaceRSSI5 >= RSII_CUTOFF || freeSpaceRSSI6 >= RSII_CUTOFF {
							wallsLoss24, wallsLoss5, wallsLoss6 = _getWallsAttenuation(matrixWithPoint.xM, matrixWithPoint.yM, walls, *sensor, client, cellSizeMeters)
						}
					}

					var tempRSSI24 float64 = freeSpaceRSSI24 + wallsLoss24 + CORRECTION_COEFFICIENT_24
					// var rssi_24 float64 = (tempRSSI24 >= RSII_CUTOFF) ? Number(tempRSSI24.toFixed(1)) : RSSI_INVISIBLE;
					var rssi24 float64
					var rssi5 float64
					var rssi6 float64
					if tempRSSI24 >= RSII_CUTOFF {
						rssi24 = Round(tempRSSI24*10) / 10 // округление до 1 знака
					} else {
						rssi24 = RSSI_INVISIBLE
					}

					var tempRSSI5 float64 = freeSpaceRSSI5 + wallsLoss5 + CORRECTION_COEFFICIENT_5
					// var rssi_5 float64 = (tempRSSI5 >= RSII_CUTOFF) ? Number(tempRSSI5.toFixed(1)) : RSSI_INVISIBLE;
					if tempRSSI5 >= RSII_CUTOFF {
						rssi5 = Round(tempRSSI5*10) / 10
					} else {
						rssi5 = RSSI_INVISIBLE
					}

					var tempRSSI6 float64 = freeSpaceRSSI6 + wallsLoss6 + CORRECTION_COEFFICIENT_6
					// var rssi_6 float64 = (tempRSSI6 >= RSII_CUTOFF) ? Number(tempRSSI6.toFixed(1)) : RSSI_INVISIBLE;
					if tempRSSI6 >= RSII_CUTOFF {
						rssi6 = Round(tempRSSI6*10) / 10
					} else {
						rssi6 = RSSI_INVISIBLE
					}

					distance = Round(distance*10) / 10

					matrixWithPoint.sensorID = sensor.ID
					matrixWithPoint.rssi24 = rssi24
					matrixWithPoint.rssi5 = rssi5
					matrixWithPoint.rssi6 = rssi6
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
 * Returns the negative numbers of total Walls attenuation for 2.4, 5 and 6 HHz bands.
 * @param clientX
 * @param clientY
 * @param Walls
 * @param sensor
 * @param Client
 * @returns
 */
func _getWallsAttenuation(clientX float64, clientY float64, walls []Wall, sensor model.Sensor, client Client, cellSizeMeters float64) (float64, float64, float64) {
	var loss24 float64 = 0
	var loss5 float64 = 0
	var loss6 float64 = 0

	for _, wall := range walls {
		var wallPathLengthThrough float64 = getWallPathLengthThrough(XYZcoordinate{x: clientX, y: clientY, z: client.ZM},
			XYZcoordinate{x: float64(*sensor.X), y: float64(*sensor.Y), z: *sensor.Z},
			XYZcoordinate{x: float64(wall.X1), y: float64(wall.Y1), z: 0},
			XYZcoordinate{x: float64(wall.X2), y: float64(wall.Y2), z: 0},
			wall.Thickness,
			cellSizeMeters)

		if wallPathLengthThrough != 0 {
			var pathDivideThickness float64 = wallPathLengthThrough / wall.Thickness
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
 * Returns the Distance in meters between Client and sensor.
 * @param clientX Client X coordinate.
 * @param clientY Client Y coordinate.
 * @param Client Client`s parameters.
 * @param sensor Sensor.
 * @returns Distance between Client and sensor in meters.
 */
func _getDistance(clientX float64, clientY float64, client Client, sensor model.Sensor, cellSizeMeters float64) float64 {
	return Magnitude(Vector{clientX - float64(*sensor.X), clientY - float64(*sensor.Y), client.ZM - *sensor.Z})
}

/**
 * Returns the free space pass loss in dB.
 * @param frequency Transmission frequency in GHz.
 * @param attenuation_factor Attenuation factor.
 * @param penetration_factor Penetration factor.
 * @param Distance Transmission Distance.
 * @returns Free space pass loss in dB.
 */
func _getFSPL(frequency int, attenuationFactor float64, penetrationFactor float64, distance float64) float64 {
	if distance < 1 {
		distance = 1
	}
	return 20*Log10(float64(frequency)) + 10*attenuationFactor*Log10(distance) + penetrationFactor - 24
}

func _approximateAzimuth(azimuth float64, delta float64) (int, error) {
	if delta == 0 {
		return 0, errors.New("delta cannot be zero")
	}

	return int(Floor((azimuth+Floor(delta/2))/delta)*delta) % 360, nil
}

/**
 * Returns the RSSI for 2.4, 5 and 6 HHz bands in a free space.
 * @param clientX Client X coordinate.
 * @param clientY Client Y coordinate.
 * @param Client Client`s parameters.
 * @param sensor Sensor.
 * @param Distance Distance between Client and Sensors in meters.
 * @returns Tuple of RSSI for 2.4, 5 and 6 HHz bands.
 */
func _getFreeSpaceRSSI(clientX float64, clientY float64, client Client, sensor model.Sensor, distance float64) (float64, float64, float64) {
	var fspl24 float64 = _getFSPL(FREQUENCY24, ATTENUATION_FACTOR24, PENETRATION_FACTOR24, distance)
	var fspl5 float64 = _getFSPL(FREQUENCY5, ATTENUATION_FACTOR5, PENETRATION_FACTOR5, distance)
	var fspl6 float64 = _getFSPL(FREQUENCY6, ATTENUATION_FACTOR6, PENETRATION_FACTOR6, distance)

	var freeSpaceRSSI24 float64 = float64(client.TrSignalPower) + float64(client.TrAntGain) - fspl24 + *sensor.CorrectionFactor24
	var freeSpaceRSSI5 float64 = float64(client.TrSignalPower) + float64(client.TrAntGain) - fspl5 + *sensor.CorrectionFactor5
	var freeSpaceRSSI6 float64 = float64(client.TrSignalPower) + float64(client.TrAntGain) - fspl6 + *sensor.CorrectionFactor6

	var antGain float64 = 2

	var diagram db.Diagram
	err := json.Unmarshal(*sensor.Diagram, &diagram)
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
            Check that the radiation diagram for sensor with ID = ${sensor.ID} is filled out correctly.
            By default, the antenna gain of ${sensor.rx_ant_gain} will be used for all directions.`)
		} else {
			horAzimuth, err := _approximateAzimuth(
				float64(getHorizontalAzimuthDeg(
					XYZcoordinate{x: float64(*sensor.X), y: float64(*sensor.Y), z: *sensor.Z},
					XYZcoordinate{clientX, clientY, client.ZM},
					0)),
				float64(delta))
			if err != nil {
				goto errorHandling
			}

			vertAzimuth, err := _approximateAzimuth(
				float64(getVerticalAzimuthDeg(
					XYZcoordinate{x: float64(*sensor.X), y: float64(*sensor.Y), z: *sensor.Z},
					XYZcoordinate{clientX, clientY, client.ZM},
					0)),
				float64(delta))
			if err != nil {
				goto errorHandling
			}

			antGain = (diagram.Degree[string(rune(horAzimuth))].HorGain + diagram.Degree[string(rune(vertAzimuth))].VertGain) / 2 // окр до десятых

		}
	} else {
		freeSpaceRSSI24 += *sensor.RxAntGain
		freeSpaceRSSI5 += *sensor.RxAntGain
		freeSpaceRSSI6 += *sensor.RxAntGain
	}

	return freeSpaceRSSI24, freeSpaceRSSI5, freeSpaceRSSI6

errorHandling:
	antGain = *sensor.RxAntGain
	freeSpaceRSSI24 += antGain
	freeSpaceRSSI5 += antGain
	freeSpaceRSSI6 += antGain
	return freeSpaceRSSI24, freeSpaceRSSI5, freeSpaceRSSI6
}
