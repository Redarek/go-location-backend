package model

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Sensor struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	X          *int      `json:"x" db:"x"`
	Y          *int      `json:"y" db:"y"`
	Z          *float64  `json:"z" db:"z"`
	MAC        string    `json:"mac" db:"mac"` //   "sensor_mac" VARCHAR(17) [unique, not null]
	IP         string    `json:"ip" db:"ip"`   //   "sensor_ip" VARCHAR(64) [not null]
	Alias      string    `json:"alias" db:"alias"`
	Interface0 *string   `json:"interface0" db:"interface_0"` //   "interface_0" VARCHAR(45) [not null]
	Interface1 *string   `json:"interface1" db:"interface_1"` //   "interface_1" VARCHAR(45)
	Interface2 *string   `json:"interface2" db:"interface_2"` //   "interface_2" VARCHAR(45)
	// TODO "state" sensors_state_enum [not null, default: "DOWN"]
	// TODO "state_change" DATETIME [not null]
	// TODO "packets_captured" INTEGER [not null, default: 0]
	// TODO  "uptime" TIME [not null]
	// TODO  "logs_path" VARCHAR(45)
	// TODO  "approved" TINYINT(1) [not null, default: FALSE]
	// TODO "mode" VARCHAR(45)
	// TODO "type" TINYINT(1)
	// TODO  "primary_channel_freq" FLOAT
	// TODO  "primary_channel_width" VARCHAR(45)
	// TODO  "primary_interval" FLOAT
	// TODO  "secondary_interval" FLOAT
	RxAntGain          *float64           `json:"rxAntGain" db:"rx_ant_gain"`                   //   "rx_ant_gain" FLOAT [not null, default: 0]
	HorRotationOffset  *int               `json:"horRotationOffset" db:"hor_rotation_offset"`   //   "hor_rotation_offset" INTEGER [not null, default: 0]
	VertRotationOffset *int               `json:"vertRotationOffset" db:"vert_rotation_offset"` //   "vert_rotation_offset" INTEGER [not null, default: 0]
	CorrectionFactor24 *float64           `json:"correctionFactor24" db:"correction_factor_24"` //   "correction_factor24" INTEGER [not null, default: 0]  -> FLOAT
	CorrectionFactor5  *float64           `json:"correctionFactor5" db:"correction_factor_5"`   //   "correction_factor5" INTEGER [not null, default: 0] -> FLOAT
	CorrectionFactor6  *float64           `json:"correctionFactor6" db:"correction_factor_6"`   //   "correction_factor6" INTEGER [not null, default: 0float64 -> FLOAT
	IsVirtual          bool               `json:"isVirtual" db:"is_virtual"`
	Diagram            *json.RawMessage   `json:"diagram" db:"diagram"` // Тип JSON
	SensorTypeID       uuid.UUID          `json:"sensorTypeId" db:"sensor_type_id"`
	FloorID            uuid.UUID          `json:"floorId" db:"floor_id"` //  "map_id" INTEGER
	CreatedAt          pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt          pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
	DeletedAt          pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
}
