package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Device struct {
	ID              uuid.UUID          `db:"id"`
	MAC             string             `db:"mac"`
	SensorID        *uuid.UUID         `db:"sensor_id"` // TODO решить ссылка или нет
	RSSI            float64            `db:"rssi"`
	Band            *string            `db:"band"`
	ChannelWidth    *string            `db:"channel_width"`
	LastContactTime pgtype.Timestamptz `db:"last_contact_time"`
}

type DeviceDetailed struct {
	Device
	FloorID uuid.UUID
}

// TODO централизовать
type SearchParameters struct {
	FloorID        uuid.UUID
	Band           string
	SensorsBetween map[uuid.UUID]BetweenTuple // TODO придумать нормальное имя
	DetectCount    int
}

type BetweenTuple struct {
	From float64
	To   float64
}
