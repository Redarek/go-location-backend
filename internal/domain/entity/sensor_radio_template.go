package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SensorRadioTemplate struct {
	ID            uuid.UUID           `db:"id"`
	Number        int                 `db:"number"`
	Channel       int                 `db:"channel"`
	Channel2      *int                `db:"channel2"`
	ChannelWidth  string              `db:"channel_width"`
	WiFi          string              `db:"wifi"`
	Power         int                 `db:"power"`
	Bandwidth     string              `db:"bandwidth"`
	GuardInterval int                 `db:"guard_interval"`
	IsActive      bool                `db:"is_active"`
	SensorTypeID  uuid.UUID           `db:"sensor_type_id"`
	CreatedAt     pgtype.Timestamptz  `db:"created_at"`
	UpdatedAt     pgtype.Timestamptz  `db:"updated_at"`
	DeletedAt     *pgtype.Timestamptz `db:"deleted_at"`
}
