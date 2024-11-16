package entity

import (
	"github.com/google/uuid"
)

type Point struct {
	ID      uuid.UUID `json:"id" db:"id"`
	FloorID uuid.UUID `json:"floorId" db:"floor_id"`
	X       float64   `json:"x" db:"x"`
	Y       float64   `json:"y" db:"y"`
}

type MatrixPoint struct {
	ID       int       `json:"id" db:"id"`
	PointID  uuid.UUID `json:"pointId" db:"point_id"`
	SensorID uuid.UUID `json:"sensorId" db:"sensor_id"`
	RSSI24   float64   `json:"rssi24" db:"rssi24"`
	RSSI5    float64   `json:"rssi5" db:"rssi5"`
	RSSI6    float64   `json:"rssi6" db:"rssi6"`
	Distance float64   `json:"distance" db:"distance"`
}

// type Matrix struct {
// 	ID           int           `json:"id" db:"id"`
// 	FloorID      uuid.UUID     `json:"floorId" db:"floor_id"`
// 	X            float64       `json:"x" db:"x"`
// 	Y            float64       `json:"y" db:"y"`
// 	MatrixPoints []MatrixPoint `json:"matrixPoints"`
// }

// TODO централизовать
type SearchParameters struct {
	FloorID        uuid.UUID
	Band           string
	SensorsBetween map[uuid.UUID]BetweenTuple // TODO придумать нормальное имя
	DetectCount    int
}

type BetweenTuple struct {
	From float64
	To   float64
}
