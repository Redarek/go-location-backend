package location

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/entity"
)

type Data struct {
	MAC     string
	FloorID uuid.UUID
	Devices []*Device
	Points  []*entity.Point
	Matrix  []*entity.MatrixPoint
}

// TODO переместить мб
type SearchParameters struct {
	// MAC string
	FloorID       uuid.UUID //? возможно избыточно
	Band          string
	SensorBetween map[uuid.UUID]BetweenTuple // TODO придумать нормальное имя
	DetectCount   int
}

type BetweenTuple struct {
	From float64
	To   float64
}

const (
	RSSI24 = "24"
	RSSI5  = "5"
	RSSI6  = "6"
)

type BandCoefficients struct {
	frequency         int
	attenuationFactor float64
	penetrationFactor float64
	Dmax              float64
}

func GetParametersForFindXY(data Data) (searchParameters SearchParameters, err error) {
	log.Info().Msgf("MAC %s search has been started...", data.MAC)
	var executionTimeStart time.Time = time.Now()

	mac := strings.ReplaceAll(strings.ToLower(data.MAC), "-", ":")

	var deviceDetections []*Device = findDeviceDetections(data.Devices, mac, data.FloorID)

	// TODO check
	if len(deviceDetections) == 0 {
		log.Info().Msgf("no devices with MAC %s were found", mac)
		return SearchParameters{}, fmt.Errorf("no devices with MAC %s were found", mac)
	}

	// var matrixCells []MatrixRow = getSearchParameters(matrix, deviceDetections, mapId)
	searchParameters = getSearchParameters(deviceDetections, data.FloorID)

	// log.Info().Msgf("%d cells have been found", searchParameters.DetectCount)
	// log.Debug().Msg(fmt.Sprintf("Founded cells:  ${JSON.stringify(matrixCells.map(cell => ([cell.point.x, cell.point.y])))}")) // TODO

	log.Debug().Msgf("Total execution time: %v", float64(time.Now().Unix()-executionTimeStart.Unix())/1000)

	return
}

// TODO фильтровать сразу в БД
// TODO rewrite
func findDeviceDetections(deviceDetections []*Device, mac string, floorID uuid.UUID) (filteredDevices []*Device) {
	var maxDate = time.Time{}

	var logMsg = "Sensors that have ever seen this device:"
	for _, detection := range deviceDetections {
		logMsg += fmt.Sprintf(`\n(Sensor ID: %s,\tBand: %s,\tRSSI: %.1f, Channel: %d, Time: %s`,
			detection.SensorID,
			detection.Band,
			detection.RSSI,
			detection.Channel,
			detection.LastContactTime,
		)
		// tempDates.push(detection.lastContactTime)
		// tempDates = append(tempDates, detection.lastContactTime)
		if detection.LastContactTime.Unix() > maxDate.Unix() {
			maxDate = detection.LastContactTime
		}
	}
	log.Debug().Msg(logMsg)

	var secondsFromLastDetection float64 = math.Floor(float64(time.Now().Unix()-maxDate.Unix()) / 1000)

	// TODO вынести в отдельный метод
	if int(secondsFromLastDetection) > INFO_AGING_TIME {
		var timeAgo = "unknown"
		if secondsFromLastDetection < 60 {
			timeAgo = fmt.Sprintf("%.1f seconds ago", secondsFromLastDetection)
		} else if secondsFromLastDetection < 120 {
			timeAgo = "a minute ago"
		} else if secondsFromLastDetection < 3600 {
			timeAgo = fmt.Sprintf("%.1f minutes ago", math.Floor(secondsFromLastDetection/60))
		} else if secondsFromLastDetection < 7200 {
			timeAgo = "a hour ago"
		} else if secondsFromLastDetection < 86400 {
			timeAgo = fmt.Sprintf("%.1f hours ago", math.Floor(secondsFromLastDetection/3600))
		} else if secondsFromLastDetection < 172800 {
			timeAgo = "yesterday"
		} else {
			timeAgo = fmt.Sprintf("%.1f days ago", math.Floor(secondsFromLastDetection/86400))
		}

		log.Warn().Msg(fmt.Sprintf("MAC %s detected %s", mac, timeAgo))
	}

	logMsg = "Sensors that were taken into account:\n"
	// var maxRSSI float64 = RSSI_INVISIBLE
	for _, detection := range deviceDetections {
		var deltaTime int = int(math.Floor(float64(maxDate.Unix()-detection.LastContactTime.Unix()) / 1000))
		if deltaTime < INFO_AGING_TIME {
			filteredDevices = append(filteredDevices, detection)
		}

		logMsg += fmt.Sprintf("(Sensor id: %s,\tRSSI: %d,\tChannel: %d,\tTime: %v)\n",
			detection.SensorID, &detection.RSSI, detection.Channel, detection.LastContactTime)
	}
	log.Debug().Msg(logMsg)

	if len(filteredDevices) < 3 {
		log.Warn().Msg("The number of sensors is less than 3. The determining accuracy will be reduced!")
	}

	return filteredDevices
}

type metric struct {
	counter  int
	rssiSum1 float64
	rssiSum2 float64
	rssiSum3 float64
	rssiSum4 float64
}

type floorMetric struct {
	floorID uuid.UUID
	metric  metric
}

type deviceList []*Device

// implement the sort interface
func (d deviceList) Len() int           { return len(d) }
func (d deviceList) Less(i, j int) bool { return d[i].RSSI < d[j].RSSI }
func (d deviceList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func determineFloor(filteredDevices []*Device) (result uuid.UUID) {
	// Группировка устройств по этажам: этаж – список устройств
	devicesPerFloor := make(map[uuid.UUID]deviceList)
	for _, device := range filteredDevices {
		devicesPerFloor[device.FloorID] = append(devicesPerFloor[device.FloorID], device)
	}

	// Вычисление метрик для каждого этажа
	metricsPerFloor := make(map[uuid.UUID]metric)
	for floorID, deviceList := range devicesPerFloor {
		metrics := calculateMetrics(deviceList)
		metricsPerFloor[floorID] = metrics
	}

	// Определение наиболее вероятного этажа
	var sortedMetrics []floorMetric
	// TODO переделать metricsPerFloor сразу в структуру
	for floorID, metric := range metricsPerFloor {
		sortedMetrics = append(sortedMetrics, floorMetric{floorID, metric})
	}

	// Сортировка метрик по приоритету
	sort.SliceStable(sortedMetrics, func(i, j int) bool {
		m1, m2 := sortedMetrics[i].metric, sortedMetrics[j].metric

		// Сортировка по количеству элементов (приоритеты 3, 4, 2, >4, 1)
		priority1 := getPriorityByCount(m1.counter)
		priority2 := getPriorityByCount(m2.counter)
		if priority1 != priority2 {
			return priority1 < priority2 // Меньший номер приоритета выше
		}

		// Сортировка по метрикам
		if m1.rssiSum1 != m2.rssiSum1 {
			return m1.rssiSum1 > m2.rssiSum1
		}
		if m1.rssiSum2 != m2.rssiSum2 {
			return m1.rssiSum2 > m2.rssiSum2
		}
		if m1.rssiSum3 != m2.rssiSum3 {
			return m1.rssiSum3 > m2.rssiSum3
		}
		if m1.rssiSum4 != m2.rssiSum4 {
			return m1.rssiSum4 > m2.rssiSum4
		}

		// Если всё равное, оставляем порядок как есть
		return false
	})

	// Если список пустой, вернуть ошибку
	if len(sortedMetrics) == 0 {
		return uuid.Nil
	}

	// Результат — этаж с лучшей метрикой
	return sortedMetrics[0].floorID

	// // ещё не работает
	// for floorID, deviceList := range devicesPerFloor {
	// 	sort.Sort(sort.Reverse(deviceList))
	// }

	// // floor2rssiMetric := make(map[uuid.UUID]metric)
	// for _, device := range filteredDevices {

	// }
	// //? определение метрики?
	// // el, isExists := floor2rssiMetric[detection.FloorID]
	// // if isExists {
	// // 	el.counter += 1
	// // 	el.rssiSum += detection.RSSI
	// // } else {
	// // 	floor2rssiMetric[detection.FloorID] = metric{}
	// // }
	// //? Determine floor
	// // macEl := 0
	// // for _, el := range floor2rssiMetric {

	// // }

	// return
}

func calculateMetrics(devices deviceList) metric {
	sort.Sort(sort.Reverse(devices))

	length := len(devices)
	metrics := metric{
		counter: length,
	}
	if length >= 1 {
		metrics.rssiSum1 = devices[0].RSSI
	}
	if length >= 2 {
		metrics.rssiSum2 = metrics.rssiSum1 + devices[1].RSSI
	}
	if length >= 3 {
		metrics.rssiSum3 = metrics.rssiSum2 + devices[2].RSSI
	}
	if length >= 4 {
		metrics.rssiSum4 = metrics.rssiSum3 + devices[3].RSSI
	}
	return metrics
}

// Функция для определения приоритета по количеству элементов
func getPriorityByCount(count int) int {
	switch count {
	case 3:
		return 1 // Высший приоритет
	case 4:
		return 2
	case 2:
		return 3
	default:
		if count > 4 {
			return 4
		}
		return 5 // Минимальный приоритет для списков с 1 элементом
	}
}

func getSearchParameters(deviceDetections []*Device, floorID uuid.UUID) (result SearchParameters) {
	band, channel := detectBandAndChannel(deviceDetections)
	var bandCoefficients BandCoefficients = getBandCoefficients(band)
	var bandAccuracyCorrection float64 = getBandAccuracyCorrection(channel)

	var sensorBetween map[uuid.UUID]BetweenTuple = make(map[uuid.UUID]BetweenTuple)
	for _, device := range deviceDetections {
		if device.RSSI == RSSI_INVISIBLE {
			continue
		}

		var deltaRSSI float64 = getDeltaRSSI(bandCoefficients, device.RSSI)

		log.Debug().Msgf("Accuracy for RSSI = %.2f is %.2f dB", device.RSSI, deltaRSSI)

		var betweenFrom float64 = device.RSSI - deltaRSSI + bandAccuracyCorrection // было округление до 2
		var betweenTo float64 = device.RSSI + deltaRSSI + bandAccuracyCorrection   // было округление до 2
		log.Debug().Msgf("between %.2f and %.2f", betweenFrom, betweenTo)

		sensorBetween[device.SensorID] = BetweenTuple{From: betweenFrom, To: betweenTo}
	}

	result = SearchParameters{
		FloorID:       floorID,
		Band:          band,
		SensorBetween: sensorBetween,
		DetectCount:   len(deviceDetections),
	}

	return
}

/**
 * Returns BandCoefficients for Band.
 * @param band Band enum.
 * @returns BandCoefficients for Band.
 */
func getBandCoefficients(band string) BandCoefficients {
	var frequency int = FREQUENCY24
	var attenuationFactor float64 = ATTENUATION_FACTOR24
	var penetrationFactor float64 = PENETRATION_FACTOR24
	var Dmax float64 = 50

	if band == RSSI5 {
		frequency = FREQUENCY5
		attenuationFactor = ATTENUATION_FACTOR5
		penetrationFactor = PENETRATION_FACTOR5
		Dmax = 55
	} else if band == RSSI6 { // пересчитать всё!
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
func getDeltaRSSI(bandCoefficients BandCoefficients, rssi float64) float64 {
	var fspl float64 = rssi - EIRP
	var d float64 = math.Pow(10, (-fspl-bandCoefficients.penetrationFactor-20*math.Log10(float64(bandCoefficients.frequency))+24)/(10*bandCoefficients.attenuationFactor))

	var delta float64 = (1 - d/bandCoefficients.Dmax) * C_MAX // было округление до 2
	if delta <= C_MIN {
		delta = C_MIN
	} else if C_MAX < delta {
		delta = C_MAX
	}

	return delta
}

/**
 * Returns band's accuracy correction.
 * @param channel Wi-Fi channel.
 * @returns
 */
func getBandAccuracyCorrection(channel int) float64 {
	var bandAccuracyCorrection float64 = 0 // Поправка на частоту канала (по умолчанию: канал 6)

	if 1 <= channel && channel <= 4 {
		bandAccuracyCorrection = -0.1
	} else if 9 <= channel && channel <= 14 {
		bandAccuracyCorrection = 0.1
	}

	return bandAccuracyCorrection
}

/**
 * Detects band and accuracy correction for it. Returns Band.RSSI24 by default.
 * @param deviceDetections Detected devices array.
 * @returns Array of Band enum and band accuracy correction coefficient.
 */
func detectBandAndChannel(deviceDetections []*Device) (string, int) {
	var band string = RSSI24 // Диапазон по умолчанию, если channel не указан или не удаётся определить
	var channel = 6

	for _, device := range deviceDetections {
		if device.Band == "" {
			continue
		} else {
			// TODO delete this later
			if device.Band == RSSI24 {
				band = RSSI24
			} else if device.Band == RSSI5 {
				band = RSSI5
			} else if device.Band == RSSI6 {
				band = RSSI5
			}

			channel = device.Channel

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

// class MacNotFoundError extends Error { }
