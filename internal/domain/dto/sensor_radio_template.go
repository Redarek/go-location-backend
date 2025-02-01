package dto

import (
	"github.com/google/uuid"
)

type CreateSensorRadioTemplateDTO struct {
	Number        int       `json:"number" db:"number"`
	Channel       int       `json:"channel" db:"channel"`
	Channel2      *int      `json:"channel2" db:"channel2"`
	ChannelWidth  string    `json:"channelWidth" db:"channel_width"`
	WiFi          string    `json:"wifi" db:"wifi"`
	Power         int       `json:"power" db:"power"`
	Bandwidth     string    `json:"bandwidth" db:"bandwidth"`
	GuardInterval int       `json:"guardInterval" db:"guard_interval"`
	IsActive      bool      `json:"isActive" db:"is_active"`
	SensorTypeID  uuid.UUID `json:"sensorTypeId" db:"sensor_type_id"`
}

type GetSensorRadioTemplateDTO struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type PatchUpdateSensorRadioTemplateDTO struct {
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
	SensorTypeID  *uuid.UUID `json:"sensorTypeId" db:"sensor_type_id"`
}

type GetSensorRadioTemplatesDTO struct {
	SensorTypeID uuid.UUID `json:"sensorTypeId" db:"sensor_type_id"`
	Page         int
	Size         int
}
