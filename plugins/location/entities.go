package location

import (
	"encoding/json"

	"github.com/google/uuid"
)

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

type Sensor struct {
	ID   uuid.UUID
	Name string
	// Color              *string
	X   *int
	Y   *int
	Z   *float64
	MAC string
	// IP                 string
	RxAntGain          float64
	HorRotationOffset  int
	VertRotationOffset int
	CorrectionFactor24 float64
	CorrectionFactor5  float64
	CorrectionFactor6  float64
	IsVirtual          bool
	Diagram            *json.RawMessage
	SensorTypeID       uuid.UUID
	FloorID            uuid.UUID
}

type Diagram struct {
	Degree map[string]Degree `json:"degree"`
}

type Degree struct {
	HorGain  float64 `json:"hor_gain"`
	VertGain float64 `json:"vert_gain"`
}

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
	FloorID       uuid.UUID
}
