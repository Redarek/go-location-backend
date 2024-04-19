package location


import "github.com/rs/zerolog/log";

// import Sensor from "../models/sensors.model";
// import Wall from "../models/walls.model";

// import { generateMatrixRow, MatrixPoint } from "./matrix_row_generator";


func _getMatrix(matrixRowGenerator: Generator<MatrixPoint>, mapId: number): [PointRow[], MatrixRow[]] {
    const pointRowsToInsert: PointRow[] = [];
    const matrixRowsToInsert: MatrixRow[] = [];
    let lastId = -1;

    for (const row of matrixRowGenerator) {

        const { id, sensorId, x_m, y_m, rssi24, rssi5, rssi6, distance } = row;

        if (lastId != id) {
            pointRowsToInsert.push({ id: id, map_id: mapId, x: x_m, y: y_m });
            lastId = id;
        }

        matrixRowsToInsert.push({ point_id: id, sensor_id: sensorId, rssi24: rssi24, rssi5: rssi5, rssi6: rssi6, distance: distance });
    }  

    return [pointRowsToInsert, matrixRowsToInsert];
}


function createMatrix(mapId: number, inputData: InputData): [PointRow[], MatrixRow[]] {
    logger.info(`Creating matrix for map_id = ${mapId}...`);
    const startTestTime: number = performance.now();

    let [pointRowsToInsert, matrixRowsToInsert]: [PointRow[], MatrixRow[]] = [[],[]];

    try {
        //const startFillTime: number = performance.now();

        const matrixRowGenerator: Generator<MatrixPoint> = generateMatrixRow(inputData);
        //const [pointSize, matrixSize]: [number, number] = await _insertIntoMatrixAsync(pointRepository, matrixRepository, matrixRowGenerator, mapId);
        [pointRowsToInsert, matrixRowsToInsert] = _getMatrix(matrixRowGenerator, mapId);

        //logger.info(`Created ${pointSize} points (${matrixSize} matrix points) in ${((performance.now() - startTestTime) / 1000).toFixed(2)} sec `
        //    + `(Del: ${deleteTime} sec, Fill: ${((performance.now() - startFillTime) / 1000).toFixed(2)} sec)`);
        logger.info(`Created ${pointRowsToInsert.length} points (${matrixRowsToInsert.length} matrix points) in ${((performance.now() - startTestTime) / 1000).toFixed(2)} sec `);
    }
    catch (err) {
        logger.fatal(err);
        //return [{ 'Error': err, 'status': 'error' }, 500];
        throw err;
    }

    //return [{ 'status': 'ok' }, 200];
    return [pointRowsToInsert, matrixRowsToInsert];
}


default createMatrix;

type Client = {
    trSignalPower: number,
    trAntGain: number,
    zM: number
}

type InputData = {
    client: Client,
    walls: Wall[],
    sensors: Sensor[],
    cell_size_meters: number,
    minX: number,
    minY: number,
    maxX: number,
    maxY: number
}

type PointRow = {
    id: number,
    map_id: number,
    x: number,
    y: number,
}

type MatrixRow = {
    point_id: number,
    sensor_id: number,
    rssi24: number,
    rssi5: number,
    rssi6: number,
    distance: number
}
