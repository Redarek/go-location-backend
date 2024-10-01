package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessPoint struct {
	ID                uuid.UUID           `db:"id"`
	Name              string              `db:"name"`
	X                 *int                `db:"x"`
	Y                 *int                `db:"y"`
	Z                 *float64            `db:"z"`
	IsVirtual         bool                `db:"is_virtual"`
	AccessPointTypeID uuid.UUID           `db:"access_point_type_id"`
	FloorID           uuid.UUID           `db:"floor_id"`
	CreatedAt         pgtype.Timestamptz  `db:"created_at"`
	UpdatedAt         pgtype.Timestamptz  `db:"updated_at"`
	DeletedAt         *pgtype.Timestamptz `db:"deleted_at"`
}

type AccessPointDetailed struct {
	AccessPoint
	AccessPointType AccessPointType
	Radios          []*AccessPointRadio
}
