package dto

import (
	"github.com/google/uuid"
)

type CreateAccessPointDTO struct {
	Name              string    `db:"name"`
	X                 *int      `db:"x"`
	Y                 *int      `db:"y"`
	Z                 *float64  `db:"z"`
	IsVirtual         bool      `db:"is_virtual"`
	AccessPointTypeID uuid.UUID `db:"access_point_type_id"`
	FloorID           uuid.UUID `db:"floor_id"`
}

type GetAccessPointDetailedDTO struct {
	ID     uuid.UUID `db:"id"`
	Limit  int
	Offset int
}

type GetAccessPointsDTO struct {
	FloorID uuid.UUID `db:"floor_id"`
	Limit   int
	Offset  int
}

type GetAccessPointsDetailedDTO struct {
	FloorID uuid.UUID `db:"floor_id"`
	Limit   int
	Offset  int
}

type PatchUpdateAccessPointDTO struct {
	ID                uuid.UUID  `db:"id"`
	Name              *string    `db:"name"`
	X                 *int       `db:"x"`
	Y                 *int       `db:"y"`
	Z                 *float64   `db:"z"`
	IsVirtual         *bool      `db:"is_virtual"`
	AccessPointTypeID *uuid.UUID `db:"access_point_type_id"`
}
