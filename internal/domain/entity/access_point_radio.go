package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessPointRadio struct {
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
	AccessPointID uuid.UUID           `db:"access_point_id"`
	CreatedAt     pgtype.Timestamptz  `db:"created_at"`
	UpdatedAt     pgtype.Timestamptz  `db:"updated_at"`
	DeletedAt     *pgtype.Timestamptz `db:"deleted_at"`
}
