package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Site struct {
	ID               uuid.UUID          `json:"id" db:"id"`
	Name             string             `json:"name" db:"name"`
	Description      *string            `json:"description" db:"description"`
	CreatedAt        pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
	DeletedAt        pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
	UserID           uuid.UUID          `json:"userId" db:"user_id"`
	Buildings        []*Building        `json:"buildings"`
	AccessPointTypes []*AccessPointType `json:"accessPointTypes"`
	WallTypes        []*WallType        `json:"wallTypes"`
	SensorTypes      []*SensorType      `json:"sensorTypes"`
}
