package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessPointRadioDTO struct {
	ID            uuid.UUID           `json:"id"`
	Number        int                 `json:"number"`
	Channel       int                 `json:"channel"`
	Channel2      *int                `json:"channel2"`
	ChannelWidth  string              `json:"channelWidth"`
	WiFi          string              `json:"wifi"`
	Power         int                 `json:"power"`
	Bandwidth     string              `json:"bandwidth"`
	GuardInterval int                 `json:"guardInterval"`
	IsActive      bool                `json:"isActive"`
	AccessPointID uuid.UUID           `json:"accessPointId"`
	CreatedAt     pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt     pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt     *pgtype.Timestamptz `json:"deletedAt"`
}

type CreateAccessPointRadioDTO struct {
	Number        int       `json:"number"`
	Channel       int       `json:"channel"`
	Channel2      *int      `json:"channel2"`
	ChannelWidth  string    `json:"channelWidth"`
	WiFi          string    `json:"wifi"`
	Power         int       `json:"power"`
	Bandwidth     string    `json:"bandwidth"`
	GuardInterval int       `json:"guardInterval"`
	IsActive      bool      `json:"isActive"`
	AccessPointID uuid.UUID `json:"accessPointId"`
}

type PatchUpdateAccessPointRadioDTO struct {
	ID            uuid.UUID  `json:"id"`
	Number        *int       `json:"number"`
	Channel       *int       `json:"channel"`
	Channel2      *int       `json:"channel2"`
	ChannelWidth  *string    `json:"channelWidth"`
	WiFi          *string    `json:"wifi"`
	Power         *int       `json:"power"`
	Bandwidth     *string    `json:"bandwidth"`
	GuardInterval *int       `json:"guardInterval"`
	IsActive      *bool      `json:"isActive"`
	AccessPointID *uuid.UUID `json:"accessPointId"`
}

type GetAccessPointRadiosDTO struct {
	AccessPointID uuid.UUID `json:"accessPointId"`
	Page          int
	Size          int
}
