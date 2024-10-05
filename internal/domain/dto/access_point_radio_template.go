package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessPointRadioTemplateDTO struct {
	ID                uuid.UUID           `db:"id"`
	Number            int                 `db:"number"`
	Channel           int                 `db:"channel"`
	Channel2          *int                `db:"channel2"`
	ChannelWidth      string              `db:"channel_width"`
	WiFi              string              `db:"wifi"`
	Power             int                 `db:"power"`
	Bandwidth         string              `db:"bandwidth"`
	GuardInterval     int                 `db:"guard_interval"`
	IsActive          bool                `db:"is_active"`
	AccessPointTypeID uuid.UUID           `db:"access_point_type_id"`
	CreatedAt         pgtype.Timestamptz  `db:"createdAt"`
	UpdatedAt         pgtype.Timestamptz  `db:"updatedAt"`
	DeletedAt         *pgtype.Timestamptz `db:"deletedAt"`
}

type CreateAccessPointRadioTemplateDTO struct {
	Number            int       `db:"number"`
	Channel           int       `db:"channel"`
	Channel2          *int      `db:"channel2"`
	ChannelWidth      string    `db:"channel_width"`
	WiFi              string    `db:"wifi"`
	Power             int       `db:"power"`
	Bandwidth         string    `db:"bandwidth"`
	GuardInterval     int       `db:"guard_interval"`
	IsActive          bool      `db:"is_active"`
	AccessPointTypeID uuid.UUID `db:"access_point_type_id"`
}

type PatchUpdateAccessPointRadioTemplateDTO struct {
	ID                uuid.UUID  `db:"id"`
	Number            *int       `db:"number"`
	Channel           *int       `db:"channel"`
	Channel2          *int       `db:"channel2"`
	ChannelWidth      *string    `db:"channel_width"`
	WiFi              *string    `db:"wifi"`
	Power             *int       `db:"power"`
	Bandwidth         *string    `db:"bandwidth"`
	GuardInterval     *int       `db:"guard_interval"`
	IsActive          *bool      `db:"is_active"`
	AccessPointTypeID *uuid.UUID `db:"access_point_type_id"`
}

type GetAccessPointRadioTemplatesDTO struct {
	AccessPointTypeID uuid.UUID `db:"access_point_type_id"`
	Limit             int
	Offset            int
}
