package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateSensorDTO struct {
	Name               string           `json:"name" db:"name"`
	Color              *string          `json:"color" db:"color"`
	X                  *int             `json:"x" db:"x"`
	Y                  *int             `json:"y" db:"y"`
	Z                  *float64         `json:"z" db:"z"`
	MAC                string           `json:"mac" db:"mac"`
	IP                 string           `json:"ip" db:"ip"`
	RxAntGain          float64          `json:"rxAntGain" db:"rx_ant_gain"`
	HorRotationOffset  int              `json:"horRotationOffset" db:"hor_rotation_offset"`
	VertRotationOffset int              `json:"vertRotationOffset" db:"vert_rotation_offset"`
	CorrectionFactor24 float64          `json:"correctionFactor24" db:"correction_factor_24"`
	CorrectionFactor5  float64          `json:"correctionFactor5" db:"correction_factor_5"`
	CorrectionFactor6  float64          `json:"correctionFactor6" db:"correction_factor_6"`
	IsVirtual          bool             `json:"isVirtual" db:"is_virtual"`
	Diagram            *json.RawMessage `json:"diagram" db:"diagram"` // Тип JSON
	SensorTypeID       uuid.UUID        `json:"sensorTypeId" db:"sensor_type_id"`
	FloorID            uuid.UUID        `json:"floorId" db:"floor_id"`
}

type GetSensorDTO struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type GetSensorDetailedDTO struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Page int
	Size int
}

type GetSensorsDTO struct {
	FloorID uuid.UUID `json:"floorId" db:"floor_id"`
	Page    int
	Size    int
}

type GetSensorsDetailedDTO struct {
	FloorID uuid.UUID `json:"floorId" db:"floor_id"`
	Page    int
	Size    int
}

type PatchUpdateSensorDTO struct {
	ID                 uuid.UUID        `json:"id" db:"id"`
	Name               *string          `json:"name" db:"name"`
	Color              *string          `json:"color" db:"color"`
	X                  *int             `json:"x" db:"x"`
	Y                  *int             `json:"y" db:"y"`
	Z                  *float64         `json:"z" db:"z"`
	MAC                *string          `json:"mac" db:"mac"`
	IP                 *string          `json:"ip" db:"ip"`
	RxAntGain          *float64         `json:"rxAntGain" db:"rx_ant_gain"`
	HorRotationOffset  *int             `json:"horRotationOffset" db:"hor_rotation_offset"`
	VertRotationOffset *int             `json:"vertRotationOffset" db:"vert_rotation_offset"`
	CorrectionFactor24 *float64         `json:"correctionFactor24" db:"correction_factor_24"`
	CorrectionFactor5  *float64         `json:"correctionFactor5" db:"correction_factor_5"`
	CorrectionFactor6  *float64         `json:"correctionFactor6" db:"correction_factor_6"`
	IsVirtual          *bool            `json:"isVirtual" db:"is_virtual"`
	Diagram            *json.RawMessage `json:"diagram" db:"diagram"` // Тип JSON
	SensorTypeID       *uuid.UUID       `json:"sensorTypeId" db:"sensor_type_id"`
}
