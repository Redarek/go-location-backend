package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateSensorDTO struct {
	Name               string           `db:"name"`
	Color              *string          `db:"color"`
	X                  *int             `db:"x"`
	Y                  *int             `db:"y"`
	Z                  *float64         `db:"z"`
	MAC                string           `db:"mac"`
	IP                 string           `db:"ip"`
	RxAntGain          float64          `db:"rx_ant_gain"`
	HorRotationOffset  int              `db:"hor_rotation_offset"`
	VertRotationOffset int              `db:"vert_rotation_offset"`
	CorrectionFactor24 float64          `db:"correction_factor_24"`
	CorrectionFactor5  float64          `db:"correction_factor_5"`
	CorrectionFactor6  float64          `db:"correction_factor_6"`
	IsVirtual          bool             `db:"is_virtual"`
	Diagram            *json.RawMessage `db:"diagram"` // Тип JSON
	SensorTypeID       uuid.UUID        `db:"sensor_type_id"`
	FloorID            uuid.UUID        `db:"floor_id"`
}

type GetSensorDetailedDTO struct {
	ID     uuid.UUID `db:"id"`
	Limit  int
	Offset int
}

type GetSensorsDTO struct {
	FloorID uuid.UUID `db:"floor_id"`
	Limit   int
	Offset  int
}

type GetSensorsDetailedDTO struct {
	FloorID uuid.UUID `db:"floor_id"`
	Limit   int
	Offset  int
}

type PatchUpdateSensorDTO struct {
	ID                 uuid.UUID        `db:"id"`
	Name               *string          `db:"name"`
	Color              *string          `db:"color"`
	X                  *int             `db:"x"`
	Y                  *int             `db:"y"`
	Z                  *float64         `db:"z"`
	MAC                *string          `db:"mac"`
	IP                 *string          `db:"ip"`
	RxAntGain          *float64         `db:"rx_ant_gain"`
	HorRotationOffset  *int             `db:"hor_rotation_offset"`
	VertRotationOffset *int             `db:"vert_rotation_offset"`
	CorrectionFactor24 *float64         `db:"correction_factor_24"`
	CorrectionFactor5  *float64         `db:"correction_factor_5"`
	CorrectionFactor6  *float64         `db:"correction_factor_6"`
	IsVirtual          *bool            `db:"is_virtual"`
	Diagram            *json.RawMessage `db:"diagram"` // Тип JSON
	SensorTypeID       *uuid.UUID       `db:"sensor_type_id"`
}
