package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SensorRadioDTO struct {
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
	SensorID      uuid.UUID           `json:"sensorId"`
	CreatedAt     pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt     pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt     *pgtype.Timestamptz `json:"deletedAt"`
}

type CreateSensorRadioDTO struct {
	Number        int       `json:"number"`
	Channel       int       `json:"channel"`
	Channel2      *int      `json:"channel2"`
	ChannelWidth  string    `json:"channelWidth"`
	WiFi          string    `json:"wifi"`
	Power         int       `json:"power"`
	Bandwidth     string    `json:"bandwidth"`
	GuardInterval int       `json:"guardInterval"`
	IsActive      bool      `json:"isActive"`
	SensorID      uuid.UUID `json:"sensorId"`
}

type PatchUpdateSensorRadioDTO struct {
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
	SensorID      *uuid.UUID `json:"sensorId"`
}

type GetSensorRadiosDTO struct {
	SensorID uuid.UUID `json:"sensorId"`
	Page     int
	Size     int
}
