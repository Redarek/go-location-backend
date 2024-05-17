package location

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// import { Repository, Sequelize } from "sequelize-typescript";
// import { Op } from "sequelize";
// import type { Attributes } from "sequelize/types/model";

// import { findMissingTablesFromList, checkExistenceOfTablesOrGetError, getTableNamesFromModels, createTablesFromArrayAsync } from "./tables_check_functionality";
// import {
//     RSSI_INVISIBLE, INFO_AGING_TIME, C_MIN, C_MAX, EIRP,
//     FREQUENCY24, FREQUENCY5, ATTENUATION_FACTOR24, PENETRATION_FACTOR24, ATTENUATION_FACTOR5, PENETRATION_FACTOR5, FREQUENCY6, PENETRATION_FACTOR6, ATTENUATION_FACTOR6
// } from "./consts";

// const sequelize = db.sequelize;

// TODO rewrite
/**
 * @param deviceRepository Device repository
 * @param mac MAC address of desired device
 * @param mapId Map id where to search
 * @returns Returns found devices
 */
func _findDeviceDetections(deviceDetections []Device, mac string, mapId uuid.UUID) []Device {
	// ?поиск устройств
	// const deviceDetections: Device[] = await deviceRepository.findAll({
	//     attributes: [
	//         "sensor_id",
	//         "rssi",
	//         "band",
	//         "channel",
	//         "last_contact_time"
	//     ],
	//     include: [{
	//         model: sequelize.model(Sensor),
	//         attributes: [],
	//         where: {
	//             map_id: mapId
	//         }
	//     }],
	//     where: {
	//         device_mac: mac
	//         //last_contact_time: {
	//         //    gte: Sequelize.literal(`NOW() - INTERVAL '${INFO_AGING_TIME} seconds'`)
	//         //}
	//     }
	// }).then(result => {
	//     if (!result.length) {
	//         throw new MacNotFoundError();
	//     }

	//     return result;
	// });

	// TODO восстановить
	// const dateFormate: Intl.DateTimeFormatOptions = {
	//     year: 'numeric',
	//     month: '2-digit',
	//     day: '2-digit',
	//     hour: '2-digit',
	//     minute: '2-digit',
	//     second: '2-digit'
	// };

	var logMsg = "Sensors that have ever seen this device:"
	var tempDates []time.Time
	for _, detection := range deviceDetections {
		// logMsg += `\n(Sensor id: ${detection.sensor_id}, RSSI: ${detection.rssi}, Channel: ${detection.channel},
		//     Time: ${detection.last_contact_time.toLocaleDateString("ru-RU", dateFormate)})`;
		// tempDates.push(detection.lastContactTime)
		tempDates = append(tempDates, detection.lastContactTime)
	}
	// logger.debug(logMsg);

	// var maxDate time.Date = time.Date(Math.max(...tempDates.map(date => date.getTime())));
	var maxDate time.Time = time.Date(math.Max()) // TODO
	// var secondsFromLastDetection time.Time = Math.floor((Date.now() - maxDate.getTime()) / 1000); // TODO

	if secondsFromLastDetection > INFO_AGING_TIME {
		var timeAgo = "unknown"
		if secondsFromLastDetection < 60 {
			timeAgo = `${secondsFromLastDetection} seconds ago`
		} else if secondsFromLastDetection < 120 {
			timeAgo = "a minute ago"
		} else if secondsFromLastDetection < 3600 {
			timeAgo = `${Math.floor(secondsFromLastDetection / 60)} minutes ago`
		} else if secondsFromLastDetection < 7200 {
			timeAgo = "a hour ago"
		} else if secondsFromLastDetection < 86400 {
			timeAgo = `${Math.floor(secondsFromLastDetection / 3600)} hours ago`
		} else if secondsFromLastDetection < 172800 {
			timeAgo = "yesterday"
		} else {
			timeAgo = `${Math.floor(secondsFromLastDetection / 86400)} days ago`
		}

		logger.warn(`MAC ${mac} last detect ${timeAgo}`)
	}

	logMsg = "Sensors that taken into account:"
	var resultDetections []Device
	for _, detection := range deviceDetections {
		var deltaTime time.Time = Math.floor((maxDate.getTime() - detection.last_contact_time.getTime()) / 1000)
		if deltaTime < INFO_AGING_TIME {
			resultDetections.push(detection)
		}

		logMsg += `\n(Sensor id: ${detection.sensor_id}, RSSI: ${detection.rssi}, Channel: ${detection.channel}, 
            Time: ${detection.last_contact_time.toLocaleDateString("ru-RU", dateFormate)})`
	}
	logger.debug(logMsg)

	if resultDetections.length < 3 {
		logger.warn("The number of sensors is less than 3. The determining accuracy will be reduced!")
	}

	return resultDetections
}

// TODO rewrite
/**
 * Returns BandCoefficients for Band.
 * @param band Band enum.
 * @returns BandCoefficients for Band.
 */
func _getBandCoefficients(band Band) BandCoefficients {
	var frequency float64 = FREQUENCY24
	var attenuationFactor float64 = ATTENUATION_FACTOR24
	var penetrationFactor float64 = PENETRATION_FACTOR24
	var Dmax float64 = 50

	if band == Band.RSSI5 {
		frequency = FREQUENCY5
		attenuationFactor = ATTENUATION_FACTOR5
		penetrationFactor = PENETRATION_FACTOR5
		Dmax = 55
	} else if band == Band.RSSI6 { // пересчитать всё!
		frequency = FREQUENCY6
		attenuationFactor = ATTENUATION_FACTOR6
		penetrationFactor = PENETRATION_FACTOR6
		Dmax = 55 // пересчитать!
	}

	return BandCoefficients{
		frequency:         frequency,
		attenuationFactor: attenuationFactor,
		penetrationFactor: penetrationFactor,
		Dmax:              Dmax,
	}
}

/**
 * Returns delta RSSI.
 * @param bandCoefficients Coefficients for band.
 * @param rssi RSSI.
 * @returns
 */
func _getDeltaRSSI(bandCoefficients BandCoefficients, rssi float64) float64 {
	var fspl float64 = rssi - EIRP
	var d float64 = 10 * *((-fspl - bandCoefficients.penetrationFactor - 20*Math.log10(bandCoefficients.frequency) + 24) / (10 * bandCoefficients.attenuationFactor))

	var delta float64 = Number(((1 - d/bandCoefficients.Dmax) * C_MAX).toFixed(2))
	if delta <= C_MIN {
		delta = C_MIN
	} else if C_MAX < delta {
		delta = C_MAX
	}

	return delta
}

func _findCellsInMatrixAsync(matrixRepository []Matrix, deviceDetections []Device, accuracy number, mapId number) []Matrix {
	band, channel := _detectBandAndChannel(deviceDetections)
	var bandCoefficients BandCoefficients = _getBandCoefficients(band)
	var bandAccuracyCorrection float64 = _getBandAccuracyCorrection(channel)

	// var whereConditionOrPartsList []Attributes<Matrix>; // TODO

	for _, device := range deviceDetections {
		if device.rssi == RSSI_INVISIBLE {
			continue
		}

		var deltaRSSI float64 = _getDeltaRSSI(bandCoefficients, device.rssi)

		log.Debug().Msg(`Accuracy for RSSI = ${device.rssi} is ${deltaRSSI} dB`)

		var betweenFrom float64 = Number((device.rssi - deltaRSSI + bandAccuracyCorrection).toFixed(2))
		var betweenTo float64 = Number((device.rssi + deltaRSSI + bandAccuracyCorrection).toFixed(2))

		// TODO rewrite
		// var rssiConditionPart Attributes<Matrix> = { rssi24: { [Op.between]: [betweenFrom, betweenTo] } };
		// if (band == Band.RSSI5)
		//     rssiConditionPart = { rssi5: { [Op.between]: [betweenFrom, betweenTo] } };
		// else if (band == Band.RSSI6)
		//     rssiConditionPart = { rssi6: { [Op.between]: [betweenFrom, betweenTo] } };

		// const whereConditionOrPart: Attributes<Matrix> = {
		//     [Op.and]: [
		//         rssiConditionPart,
		//         { sensor_id: device.sensor_id }
		//     ]
		// };
		// whereConditionOrPartsList.push(whereConditionOrPart);
	}

	// TODO do sth with this
	// const res: Matrix[] = await matrixRepository.findAll({
	//     attributes: [
	//         "point_id",
	//         "rssi24",
	//         "rssi5",
	//         "distance",
	//         [sequelize.fn('COUNT', sequelize.col('*')), 'count']
	//     ],
	//     include: [{
	//         model: sequelize.model(Point),
	//         attributes: [
	//             "map_id",
	//             "x",
	//             "y"
	//         ],
	//         where: {
	//             map_id: mapId
	//         }
	//     }],
	//     where: {
	//         [Op.or]: whereConditionOrPartsList
	//     },
	//     group: ['point_id'],
	//     having: { count: deviceDetections.length }
	// }).then(result => {
	//     if (!result.length) {
	//         throw new MacNotFoundError();
	//     }

	//     return result;
	// });

	return res
}

/**
 * Returns band's accuracy correction.
 * @param channel Wi-Fi channel.
 * @returns
 */
func _getBandAccuracyCorrection(channel int) float64 {
	var bandAccuracyCorrectiont float64 = 0 // Поправка на частоту канала (по умолчанию: канал 6)

	if 1 <= channel && channel <= 4 {
		bandAccuracyCorrectiont = -0.1
	} else if 9 <= channel && channel <= 14 {
		bandAccuracyCorrectiont = 0.1
	}

	return bandAccuracyCorrectiont
}

/**
 * Detects band and accuracy correction for it. Returns Band.RSSI24 by default.
 * @param deviceDetections Detected devices array.
 * @returns Array of Band enum and band accuracy correction coefficient.
 */
func _detectBandAndChannel(deviceDetections []Device) (Band, number) {
	var band Band = Band.RSSI24 // Диапазон по умолчанию, если channel не указан или не удаётся определить
	var channel = 6

	for _, device := range deviceDetections {
		if device.band == null {
			continue
		} else {
			if device.band == Band.RSSI24 {
				band = Band.RSSI24
			} else if device.band == Band.RSSI5 {
				band = Band.RSSI5
			} else if device.band == Band.RSSI6 {
				band = Band.RSSI5
			}

			channel = device.channel

			break
		}
	}

	return band, channel
}

//# def _evaluate_log_accuracy_correction(distance: float, band: str):
//#     def fspl(dist: float, band: str) -> float:
//#         if band == "24":
//#             return 20 * math.log10(FREQUENCY24) + 10 * ATTENUATION_FACTOR24 * math.log10(dist) + PENETRATION_FACTOR24 - 24
//#         else:
//#             return 20 * math.log10(FREQUENCY5) + 10 * ATTENUATION_FACTOR5 * math.log10(dist) + PENETRATION_FACTOR5 - 24
//#
//#     max24 = fspl(40, "24")
//#     max5 = fspl(40, "5")
//#
//#     def normal_fspl(dist: float, band: str) -> float:
//#         if band == "24":
//#             res = fspl(dist, band) / max24
//#         else:
//#             res = fspl(dist, band) / max5
//#
//#         return res if res <= 1 else 1
//#
//#     def invert_normal_fspl(dist: float, band: str) -> float:
//#         res = 1 - normal_fspl(dist, band)
//#         return res if res != 0 else 0.01
//#
//#     return round(invert_normal_fspl(distance, band) * 20, 1)

//def _find_cells_in_matrix(db: SQLAlchemy, matrix: Table, founded_device_detections: list, accuracy: float) -> list:
//    res_length_coefficient: float = ONE_SENSOR_RESULT_LENGTH if len(founded_device_detections) == 1 else 1

//    band, band_accuracy_correction = _detect_band_and_accuracy_correction(founded_device_detections)

//    # last_marker = 0
//    min_accuracy = MIN_ACCURACY
//    last_step = ACCURACY_BIG_STEP
//    if accuracy:
//        i_accuracy = min_accuracy = accuracy
//    else:
//        i_accuracy = MAX_ACCURACY

//    result = []
//    last_result1 = None
//    # last_result2 = None
//    # last_result3 = None
//    while i_accuracy >= min_accuracy:
//        filter_query = None
//        for device in founded_device_detections:
//            if device.rssi == RSSI_INVISIBLE:
//                continue

//            temp_filter_part = getattr(matrix.c, f"rssi{device.sensor_id}_{band}").\
//                between(device.rssi - round(i_accuracy, 1) + band_accuracy_correction,
//                        device.rssi + round(i_accuracy, 1) + band_accuracy_correction)
//            filter_query = temp_filter_part if filter_query is None else filter_query & temp_filter_part

//        result = db.session.query(matrix.c.x_m, matrix.c.y_m).filter(filter_query).all()

//        # last_marker += 1
//        # last_result3 = last_result2
//        last_result2 = last_result1
//        last_result1 = result

//        if accuracy:
//            break

//        if 0 <= len(result) <= RESULT_LENGTH_BIG * res_length_coefficient:
//            if 0 <= len(result) <= RESULT_LENGTH_SMALL * res_length_coefficient:
//                if not len(result):
//                    if last_result2 is not None:
//                        result = last_result2
//                        i_accuracy += last_step
//                break

//            i_accuracy -= ACCURACY_SMALL_STEP
//            last_step = ACCURACY_SMALL_STEP
//        else:
//            i_accuracy -= ACCURACY_BIG_STEP
//            last_step = ACCURACY_BIG_STEP

//    if not accuracy:
//        i_accuracy += last_step

//    i_accuracy = round(i_accuracy, 1)

//    logger.debug(f"Accuracy is {i_accuracy} dB")
//    return result

// ?вынести
/**
 * Deletes all points with mapId from table.
 * @param pointRepository Point repository
 * @param mapId Id of the map.
 */
// func _deleteDevicePointsWithMapIdAsync(devicePointRepository []DevicePoint, mapId uuid.UUID) {
//     log.Debug.Msg(`Trying to delete points with map_id = ${mapId} from ${devicePointRepository.tableName} table.`);

//     var deleted: number = await devicePointRepository.destroy({
//         where: { map_id: mapId }
//     }).catch(err => {
//         logger.error(`Something went wrong when trying to delete points with map_id = ${mapId}`);
//         throw err;
//     })

//     logger.debug(`${deleted} rows have been deleted from ${devicePointRepository.tableName}`);
// }

// ?вынести
// /**
//  * Deletes rows from device_points by map_id, then inserts new values.
//  * @param devicePointRepository DevicePoint repository.
//  * @param matrixCells Cells to insert.
//  * @param mac Desired MAC.
//  * @param mapId Map id.
//  */
// func _insertOrUpdateDevicePointsWithMapIdAsync(devicePointRepository: Repository<DevicePoint>, matrixCells: Matrix[], mac: string, mapId: number): Promise<number> {
//     await _deleteDevicePointsWithMapIdAsync(devicePointRepository, mapId);

//     await devicePointRepository.bulkCreate(matrixCells.map(matrixCell => ({
//         point_id: matrixCell.point_id,
//         map_id: matrixCell.point.map_id,
//         // device_id: ?
//         device_mac: mac,
//         // color: default
//         x: matrixCell.point.x,
//         y: matrixCell.point.y,
//         // location_time: default
//     })), { logging: false });

//     return matrixCells.length;
// }

func findXY(mac string, mapId uuid.UUID) (string, unknown, int) {
	log.Info().Msg("MAC ${mac} search has been started...")
	var executionTimeStart time.Time = time.Time.now()

	// const sequelize: Sequelize = db.sequelize;
	mac = mac.toLowerCase().replace("-", ":")

	// не нужно
	// log.Info().Msg("Checking existence of required tables in database...");
	// try {
	//     const existingTables: string[] = await db.showAllTablesAsync();

	//     const requiredTables: string[] = getTableNamesFromModels(sequelize, [FloorMap, Point, Sensor, Device, Matrix]);
	//     checkExistenceOfTablesOrGetError(existingTables, requiredTables);
	//     logger.info(`All necessary tables exist: ${requiredTables.join(", ")}.`);

	//     const optionalTables: string[] = getTableNamesFromModels(sequelize, [DevicePoint]);
	//     const missingOptionalTables: string[] = findMissingTablesFromList(existingTables, optionalTables);
	//     if (missingOptionalTables.length) {
	//         logger.info(`The following optional tables do not exist: ${missingOptionalTables.join(", ")}. Creating...`);
	//         await createTablesFromArrayAsync(missingOptionalTables);
	//     }
	//     logger.info(`All optional tables exist: ${optionalTables.join(", ")}.`);
	// }
	// catch (err) {
	//     logger.error(`Something went wrong while checking the existence of tables:\n ${err}`);
	//     logger.fatal(err);
	//     return [{ 'Error': err, 'status': 'error' }, 500]
	// }

	// const matrixRepository: Repository<Matrix> = db.sequelize.getRepository(Matrix);
	// const deviceRepository: Repository<Device> = db.sequelize.getRepository(Device);
	// const devicePointRepository: Repository<DevicePoint> = db.sequelize.getRepository(DevicePoint);

	var deviceDetections []Device
	// try {
	deviceDetections = _findDeviceDetections(deviceRepository, mac, mapId)
	// }
	// catch (err) {
	//     if (err instanceof MacNotFoundError) {
	//         logger.warn(`MAC ${mac} not found or incorrect`);
	//         return [{ 'Error': err, 'status': 'error' }, 404];
	//     }

	//     logger.error("Unexpected error occurred while searching for devices that have detected this MAC");
	//     logger.fatal(err);
	//     return [{ 'Error': err, 'status': 'error' }, 500];
	// }

	var matrixCells []Matrix
	// try {
	matrixCells = _findCellsInMatrixAsync(matrixRepository, deviceDetections, accuracy, mapId)
	logger.info(`${matrixCells.length} cells have been found`)
	logger.debug(`Founded cells:  ${JSON.stringify(matrixCells.map(cell => ([cell.point.x, cell.point.y])))}`)
	// }
	// catch (err) {
	//     if (err instanceof MacNotFoundError) {
	//         logger.warn(`MAC ${mac} not found`);
	//         return [{ 'Error': err, 'status': 'error' }, 404];
	//     }

	//     logger.error(`Unexpected error occurred while creating ${matrixRepository.tableName}...`);
	//     logger.fatal(err);
	//     return [{ 'Error': err, 'status': 'error' }, 500];
	// }

	// try {
	var result = _insertOrUpdateDevicePointsWithMapIdAsync(devicePointRepository, matrixCells, mac, mapId)
	log.Info().Msg(`${devicePointRepository.tableName} table has been filled successfully with ${result} rows`)
	// }
	// catch (err) {
	//     logger.error(`Failed to delete some rows or fill ${devicePointRepository.tableName}`);
	//     logger.fatal(err);
	//     return [{ 'Error': err, 'status': 'error' }, 500]
	// }

	log.Debug().Msg(`Total execution time: ${((Date.now() - executionTimeStart) / 1000).toFixed(2)}`)
	return 200, err
}

const (
	RSSI24 = "24"
	RSSI5  = "5"
	RSSI6  = "6"
)

type Data struct {
	mac   string
	mapId uuid.UUID
}

type Device struct {
	id              uuid.UUID
	sensorId        uuid.UUID
	rssi            float64
	band            string
	channel         int
	lastContactTime time.Time
}

type BandCoefficients struct {
	frequency         float64
	attenuationFactor float64
	penetrationFactor float64
	Dmax              float64
}

// class MacNotFoundError extends Error { }
