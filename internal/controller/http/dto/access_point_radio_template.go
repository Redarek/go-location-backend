package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessPointRadioTemplateDTO struct {
	ID                uuid.UUID           `json:"id"`
	Number            int                 `json:"number"`
	Channel           int                 `json:"channel"`
	Channel2          *int                `json:"channel2"`
	ChannelWidth      string              `json:"channelWidth"`
	WiFi              string              `json:"wifi"`
	Power             int                 `json:"power"`
	Bandwidth         string              `json:"bandwidth"`
	GuardInterval     int                 `json:"guardInterval"`
	AccessPointTypeID uuid.UUID           `json:"accessPointTypeId"`
	CreatedAt         pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt         pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt         *pgtype.Timestamptz `json:"deletedAt"`
}

type CreateAccessPointRadioTemplateDTO struct {
	Number            int       `json:"number"`
	Channel           int       `json:"channel"`
	Channel2          *int      `json:"channel2"`
	ChannelWidth      string    `json:"channelWidth"`
	WiFi              string    `json:"wifi"`
	Power             int       `json:"power"`
	Bandwidth         string    `json:"bandwidth"`
	GuardInterval     int       `json:"guardInterval"`
	AccessPointTypeID uuid.UUID `json:"accessPointTypeId"`
}

type PatchUpdateAccessPointRadioTemplateDTO struct {
	ID                uuid.UUID  `json:"id"`
	Number            *int       `json:"number"`
	Channel           *int       `json:"channel"`
	Channel2          *int       `json:"channel2"`
	ChannelWidth      *string    `json:"channelWidth"`
	WiFi              *string    `json:"wifi"`
	Power             *int       `json:"power"`
	Bandwidth         *string    `json:"bandwidth"`
	GuardInterval     *int       `json:"guardInterval"`
	AccessPointTypeID *uuid.UUID `json:"accessPointTypeId"`
}

type GetAccessPointRadioTemplatesDTO struct {
	AccessPointTypeID uuid.UUID `json:"accessPointTypeId"`
	Page              int
	Size              int
}
