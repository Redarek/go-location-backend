package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SensorRadio struct {
	ID            uuid.UUID          `json:"id" db:"id"`
	Number        int                `json:"number" db:"number"`
	Channel       int                `json:"channel" db:"channel"`
	ChannelWidth  int                `json:"channelWidth" db:"channel_width"`
	WiFi          string             `json:"wifi" db:"wifi"`
	Power         int                `json:"power" db:"power"`
	Bandwidth     string             `json:"bandwidth" db:"bandwidth"`
	GuardInterval int                `json:"guardInterval" db:"guard_interval"`
	IsActive      bool               `json:"isActive" db:"is_active"`
	SensorID      uuid.UUID          `json:"sensorId" db:"sensor_id"`
	CreatedAt     pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt     pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
	DeletedAt     pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
}
