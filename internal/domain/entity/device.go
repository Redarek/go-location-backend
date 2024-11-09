package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Device struct {
	ID              uuid.UUID          `db:"id"`
	MAC             string             `db:"mac"`
	SensorID        *uuid.UUID         `db:"sensor_id"`
	RSSI            float64            `db:"rssi"`
	Band            *string            `db:"band"`
	ChannelWidth    *string            `db:"channel_width"`
	LastContactTime pgtype.Timestamptz `db:"last_contact_time"`
}

// type DeviceDetailed struct {
// 	Device
// 	Sensor Sensor
// }
