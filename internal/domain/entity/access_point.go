package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessPoint struct {
	ID                uuid.UUID           `json:"id" db:"id"`
	Name              string              `json:"name" db:"name"`
	Color             *string             `json:"color" db:"color"`
	X                 *int                `json:"x" db:"x"`
	Y                 *int                `json:"y" db:"y"`
	Z                 *float64            `json:"z" db:"z"`
	IsVirtual         bool                `json:"isVirtual" db:"is_virtual"`
	AccessPointTypeID uuid.UUID           `json:"accessPointTypeId" db:"access_point_type_id"`
	FloorID           uuid.UUID           `json:"floorId" db:"floor_id"`
	CreatedAt         pgtype.Timestamptz  `json:"createdAt" db:"created_at"`
	UpdatedAt         pgtype.Timestamptz  `json:"updatedAt" db:"updated_at"`
	DeletedAt         *pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
}

type AccessPointDetailed struct {
	AccessPoint
	AccessPointType AccessPointType     `json:"accessPointType"`
	Radios          []*AccessPointRadio `json:"radios"`
}
