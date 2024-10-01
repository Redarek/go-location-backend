package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessPointRadioDTO struct {
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
	CreatedAt     pgtype.Timestamptz  `db:"createdAt"`
	UpdatedAt     pgtype.Timestamptz  `db:"updatedAt"`
	DeletedAt     *pgtype.Timestamptz `db:"deletedAt"`
}

type CreateAccessPointRadioDTO struct {
	Number        int       `db:"number"`
	Channel       int       `db:"channel"`
	Channel2      *int      `db:"channel2"`
	ChannelWidth  string    `db:"channel_width"`
	WiFi          string    `db:"wifi"`
	Power         int       `db:"power"`
	Bandwidth     string    `db:"bandwidth"`
	GuardInterval int       `db:"guard_interval"`
	IsActive      bool      `db:"is_active"`
	AccessPointID uuid.UUID `db:"access_point_id"`
}

type PatchUpdateAccessPointRadioDTO struct {
	ID            uuid.UUID  `db:"id"`
	Number        *int       `db:"number"`
	Channel       *int       `db:"channel"`
	Channel2      *int       `db:"channel2"`
	ChannelWidth  *string    `db:"channel_width"`
	WiFi          *string    `db:"wifi"`
	Power         *int       `db:"power"`
	Bandwidth     *string    `db:"bandwidth"`
	GuardInterval *int       `db:"guard_interval"`
	IsActive      *bool      `db:"is_active"`
	AccessPointID *uuid.UUID `db:"access_point_id"`
}

type GetAccessPointRadiosDTO struct {
	AccessPointID uuid.UUID `db:"access_point_id"`
	Limit         int
	Offset        int
}
