package dto

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SensorDTO struct {
	ID                 uuid.UUID           `json:"id"`
	Name               string              `json:"name"`
	Color              *string             `json:"color"`
	X                  *int                `json:"x"`
	Y                  *int                `json:"y"`
	Z                  *float64            `json:"z"`
	MAC                string              `json:"mac"`
	IP                 string              `json:"ip"`
	RxAntGain          float64             `json:"rxAntGain"`
	HorRotationOffset  int                 `json:"horRotationOffset"`
	VertRotationOffset int                 `json:"vertRotationOffset"`
	CorrectionFactor24 float64             `json:"correctionFactor24"`
	CorrectionFactor5  float64             `json:"correctionFactor5"`
	CorrectionFactor6  float64             `json:"correctionFactor6"`
	IsVirtual          bool                `json:"isVirtual"`
	Diagram            *json.RawMessage    `json:"diagram"` // Тип JSON
	SensorTypeID       uuid.UUID           `json:"sensorTypeId"`
	FloorID            uuid.UUID           `json:"floorId"`
	CreatedAt          pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt          pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt          *pgtype.Timestamptz `json:"deletedAt"`
}

type SensorDetailedDTO struct {
	SensorDTO
	SensorType SensorTypeDTO     `json:"accessPointType"`
	Radios     []*SensorRadioDTO `json:"radios"`
}

type CreateSensorDTO struct {
	Name               string           `json:"name"`
	Color              *string          `json:"color"`
	X                  *int             `json:"x"`
	Y                  *int             `json:"y"`
	Z                  *float64         `json:"z"`
	MAC                string           `json:"mac"`
	IP                 string           `json:"ip"`
	RxAntGain          float64          `json:"rxAntGain"`
	HorRotationOffset  int              `json:"horRotationOffset"`
	VertRotationOffset int              `json:"vertRotationOffset"`
	CorrectionFactor24 float64          `json:"correctionFactor24"`
	CorrectionFactor5  float64          `json:"correctionFactor5"`
	CorrectionFactor6  float64          `json:"correctionFactor6"`
	IsVirtual          bool             `json:"isVirtual"`
	Diagram            *json.RawMessage `json:"diagram"` // Тип JSON
	SensorTypeID       uuid.UUID        `json:"sensorTypeId"`
	FloorID            uuid.UUID        `json:"floorId"`
}

type GetSensorDetailedDTO struct {
	ID   uuid.UUID `json:"id"`
	Page int
	Size int
}

type GetSensorsDTO struct {
	FloorID uuid.UUID `json:"floorId"`
	Page    int
	Size    int
}

type GetSensorsDetailedDTO struct {
	FloorID uuid.UUID `json:"floorId"`
	Page    int
	Size    int
}

type PatchUpdateSensorDTO struct {
	ID                 uuid.UUID        `json:"id"`
	Name               *string          `json:"name"`
	Color              *string          `json:"color"`
	X                  *int             `json:"x"`
	Y                  *int             `json:"y"`
	Z                  *float64         `json:"z"`
	MAC                *string          `json:"mac"`
	IP                 *string          `json:"ip"`
	RxAntGain          *float64         `json:"rxAntGain"`
	HorRotationOffset  *int             `json:"horRotationOffset"`
	VertRotationOffset *int             `json:"vertRotationOffset"`
	CorrectionFactor24 *float64         `json:"correctionFactor24"`
	CorrectionFactor5  *float64         `json:"correctionFactor5"`
	CorrectionFactor6  *float64         `json:"correctionFactor6"`
	IsVirtual          *bool            `json:"isVirtual"`
	Diagram            *json.RawMessage `json:"diagram"` // Тип JSON
	SensorTypeID       *uuid.UUID       `json:"sensorTypeId"`
}
