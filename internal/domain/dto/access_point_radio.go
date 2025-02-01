package dto

import (
	"github.com/google/uuid"
)

type CreateAccessPointRadioDTO struct {
	Number        int       `json:"number" db:"number"`
	Channel       int       `json:"channel" db:"channel"`
	Channel2      *int      `json:"channel2" db:"channel2"`
	ChannelWidth  string    `json:"channelWidth" db:"channel_width"`
	WiFi          string    `json:"wifi" db:"wifi"`
	Power         int       `json:"power" db:"power"`
	Bandwidth     string    `json:"bandwidth" db:"bandwidth"`
	GuardInterval int       `json:"guardInterval" db:"guard_interval"`
	IsActive      bool      `json:"isActive" db:"is_active"`
	AccessPointID uuid.UUID `json:"accessPointId" db:"access_point_id"`
}

type GetAccessPointRadioDTO struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type PatchUpdateAccessPointRadioDTO struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Number        *int       `json:"number" db:"number"`
	Channel       *int       `json:"channel" db:"channel"`
	Channel2      *int       `json:"channel2" db:"channel2"`
	ChannelWidth  *string    `json:"channelWidth" db:"channel_width"`
	WiFi          *string    `json:"wifi" db:"wifi"`
	Power         *int       `json:"power" db:"power"`
	Bandwidth     *string    `json:"bandwidth" db:"bandwidth"`
	GuardInterval *int       `json:"guardInterval" db:"guard_interval"`
	IsActive      *bool      `json:"isActive" db:"is_active"`
	AccessPointID *uuid.UUID `json:"accessPointId" db:"access_point_id"`
}

type GetAccessPointRadiosDTO struct {
	AccessPointID uuid.UUID `json:"accessPointId" db:"access_point_id"`
	Page          int
	Size          int
}
