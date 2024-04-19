package location;
// import logger from "../../../logger"
// import Sensor from "../models/sensors.model"
// import Wall from "../models/walls.model"
import ("consts");
// import { getHorizontalAzimuthDeg, getVerticalAzimuthDeg, getWallPathLengthThrough } from "./math_functionality";

// import type { Client, InputData } from "./matrix_creator"
//import type { XYcoordinate, XYZcoordinate } from "./math_functionality";

// type MatrixPoint = {
//     id: number
//     sensorId: number,
//     x: number,
//     y: number,
//     x_m: number,
//     y_m: number,
//     rssi24: number,
//     rssi5: number,
//     rssi6: number,
//     distance: number
// }


func generateMatrixRow(inputData: InputData): Generator<MatrixPoint> { // generator
    const { client, walls, sensors, cell_size_meters, minX, minY, maxX, maxY } = inputData;

    let i = 0;

    for (let y = minY; y < maxY + 1; y++) {
        for (let x = minX; x < maxX + 1; x++) {
            i++;
            const matrixWithPoint: MatrixPoint = {
                id: i,
                sensorId: -1,
                x: x,
                y: y,
                x_m: x * cell_size_meters,
                y_m: y * cell_size_meters,
                rssi24: RSSI_INVISIBLE,
                rssi5: RSSI_INVISIBLE,
                rssi6: RSSI_INVISIBLE,
                distance: DISTANCE_INVISIBLE
            }

            for (const sensor of sensors) {
                let distance: number = _getDistance(x, y, client, sensor, cell_size_meters);
                const [freeSpaceRSSI24, freeSpaceRSSI5, freeSpaceRSSI6] = _getFreeSpaceRSSI(x, y, client, sensor, distance);
                let wallsLoss24 = 0;
                let wallsLoss5 = 0;
                let wallsLoss6 = 0;

                if (CALCULATE_WALLS) {
                    if (freeSpaceRSSI24 >= RSII_CUTOFF || freeSpaceRSSI5 >= RSII_CUTOFF || freeSpaceRSSI6 >= RSII_CUTOFF) {
                        [wallsLoss24, wallsLoss5, wallsLoss6] = _getWallsAttenuation(x, y, walls, sensor, client, cell_size_meters);
                    }
                }

                const tempRSSI24: number = freeSpaceRSSI24 + wallsLoss24 + CORRECTION_COEFFICIENT_24;
                const rssi_24: number = (tempRSSI24 >= RSII_CUTOFF) ? Number(tempRSSI24.toFixed(1)) : RSSI_INVISIBLE;
                const tempRSSI5: number = freeSpaceRSSI5 + wallsLoss5 + CORRECTION_COEFFICIENT_5;
                const rssi_5: number = (tempRSSI5 >= RSII_CUTOFF) ? Number(tempRSSI5.toFixed(1)) : RSSI_INVISIBLE;
                const tempRSSI6: number = freeSpaceRSSI6 + wallsLoss6 + CORRECTION_COEFFICIENT_6;
                const rssi_6: number = (tempRSSI6 >= RSII_CUTOFF) ? Number(tempRSSI6.toFixed(1)) : RSSI_INVISIBLE;
                distance = Number(distance.toFixed(1));

                matrixWithPoint.sensorId = sensor.id;
                matrixWithPoint.rssi24 = rssi_24;
                matrixWithPoint.rssi5 = rssi_5;
                matrixWithPoint.rssi6 = rssi_6;
                matrixWithPoint.distance = distance;

                yield matrixWithPoint;
            }
        }
    }
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
function _getWallsAttenuation(clientX: number, clientY: number, walls: Wall[], sensor: Sensor, client: Client, cell_size_meters: number): [number, number, number] { 
    let loss24 = 0;
    let loss5 = 0;
    let loss6 = 0;

    for (const wall of walls) {
        const wall_path_length_through: number = getWallPathLengthThrough([clientX, clientY, client.zM],
            [sensor.x as number, sensor.y as number, sensor.z as number],
            [wall.x1, wall.y1, 0],
            [wall.x2, wall.y2, 0],
            wall.thickness,
            cell_size_meters);

        if (wall_path_length_through) {
            const pathDivideThickness: number = wall_path_length_through / wall.thickness;
            loss24 -= wall.atten24 * pathDivideThickness;
            loss5 -= wall.atten5 * pathDivideThickness;
            loss6 -= wall.atten6 * pathDivideThickness;

            if (loss24 <= RSII_CUTOFF && loss5 <= RSII_CUTOFF && loss6 <= RSII_CUTOFF)
                break;
        }
    }

    return [loss24, loss5, loss6];
}


/**
 * Returns the distance in meters between client and sensor.
 * @param clientX Client x coordinate.
 * @param clientY Client y coordinate.
 * @param client Client`s parameters.
 * @param sensor Sensor. 
 * @returns Distance between client and sensor in meters.
 */
function _getDistance(clientX: number, clientY: number, client: Client, sensor: Sensor, cell_size_meters: number): number {
    return Math.hypot((clientX - (sensor.x as number)) * cell_size_meters, (clientY - (sensor.y as number)) * cell_size_meters, (client.zM - (sensor.z as number)));
}

/**
 * Returns the free space pass loss in dB.
 * @param frequency Transmission frequency in GHz.
 * @param attenuation_factor Attenuation factor.
 * @param penetration_factor Penetration factor.
 * @param distance Transmission distance.
 * @returns Free space pass loss in dB.
 */
function _getFSPL(frequency: number, attenuation_factor: number, penetration_factor: number, distance: number): number {
    if (distance < 1) {
        distance = 1
    }
    return 20 * Math.log10(frequency) + 10 * attenuation_factor * Math.log10(distance) + penetration_factor - 24
}

function _approximateAzimuth(azimuth: number, delta: number): number {
    return (Math.floor((azimuth + Math.floor(delta / 2)) / delta) * delta) % 360;
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
function _getFreeSpaceRSSI(clientX: number, clientY: number, client: Client, sensor: Sensor, distance: number): [number, number, number] { 
    const fspl24: number = _getFSPL(FREQUENCY24, ATTENUATION_FACTOR24, PENETRATION_FACTOR24, distance);
    const fspl5: number = _getFSPL(FREQUENCY5, ATTENUATION_FACTOR5, PENETRATION_FACTOR5, distance);
    const fspl6: number = _getFSPL(FREQUENCY6, ATTENUATION_FACTOR6, PENETRATION_FACTOR6, distance);

    let freeSpaceRSSI24: number = client.trSignalPower + client.trAntGain - fspl24 + sensor.correction_factor24;
    let freeSpaceRSSI5: number = client.trSignalPower + client.trAntGain - fspl5 + sensor.correction_factor5;
    let freeSpaceRSSI6: number = client.trSignalPower + client.trAntGain - fspl6 + sensor.correction_factor6;

    if (sensor.radiation_diagram != null) {
        let delta = 0;
        if (sensor.radiation_diagram.diagram.degree[10] !== undefined) {
            delta = 10;
        }
        if (sensor.radiation_diagram.diagram.degree[15] !== undefined) {
            delta = 15;
        }

        if (!delta) {
            logger.warn(`The radiation diagram can have only the step of 10 or 15 degrees. 
            Check that the radiation diagram for sensor with id = ${sensor.id} is filled out correctly. 
            By default, the antenna gain of ${sensor.rx_ant_gain} will be used for all directions.`)
        } else {
            let ant_gain = 2;
            try {
                const hor_azimuth: number = _approximateAzimuth(
                    getHorizontalAzimuthDeg(
                        [sensor.x as number, sensor.y as number, sensor.z as number],
                        [clientX, clientY, client.zM]),
                    delta
                )
                const vert_azimuth: number = _approximateAzimuth(
                    getVerticalAzimuthDeg(
                        [sensor.x as number, sensor.y as number, sensor.z as number],
                        [clientX, clientY, client.zM]),
                    delta
                )

                ant_gain = Number(((sensor.radiation_diagram.diagram.degree[hor_azimuth].hor_gain + sensor.radiation_diagram.diagram.degree[vert_azimuth].vert_gain) / 2).toFixed(1));
            }
            catch (err) {
                ant_gain = sensor.rx_ant_gain;
                freeSpaceRSSI24 += ant_gain;
                freeSpaceRSSI5 += ant_gain;
                freeSpaceRSSI6 += ant_gain;
            }
        }
    } else {
        freeSpaceRSSI24 += sensor.rx_ant_gain
        freeSpaceRSSI5 += sensor.rx_ant_gain
        freeSpaceRSSI6 += sensor.rx_ant_gain
    }
        
    return [freeSpaceRSSI24, freeSpaceRSSI5, freeSpaceRSSI6];
}