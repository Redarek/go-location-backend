package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SensorRadioTemplateDTO struct {
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
	SensorTypeID  uuid.UUID           `json:"sensorTypeId"`
	CreatedAt     pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt     pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt     *pgtype.Timestamptz `json:"deletedAt"`
}

type CreateSensorRadioTemplateDTO struct {
	Number        int       `json:"number"`
	Channel       int       `json:"channel"`
	Channel2      *int      `json:"channel2"`
	ChannelWidth  string    `json:"channelWidth"`
	WiFi          string    `json:"wifi"`
	Power         int       `json:"power"`
	Bandwidth     string    `json:"bandwidth"`
	GuardInterval int       `json:"guardInterval"`
	IsActive      bool      `json:"isActive"`
	SensorTypeID  uuid.UUID `json:"sensorTypeId"`
}

type PatchUpdateSensorRadioTemplateDTO struct {
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
	SensorTypeID  *uuid.UUID `json:"sensorTypeId"`
}

type GetSensorRadioTemplatesDTO struct {
	SensorTypeID uuid.UUID `json:"sensorTypeId"`
	Page         int
	Size         int
}
