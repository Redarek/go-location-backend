package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Wall struct {
	ID         uuid.UUID          `json:"id" db:"id"`
	X1         int                `json:"x1" db:"x1"`
	Y1         int                `json:"y1" db:"y1"`
	X2         int                `json:"x2" db:"x2"`
	Y2         int                `json:"y2" db:"y2"`
	WallTypeID uuid.UUID          `json:"wallTypeId" db:"wall_type_id"`
	FloorID    uuid.UUID          `json:"floorId" db:"floor_id"`
	CreatedAt  pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt  pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
	DeletedAt  pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
}

type WallDetailed struct {
	Wall
	WallType WallType `json:"wallType"`
}
