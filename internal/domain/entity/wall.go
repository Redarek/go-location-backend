package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Wall struct {
	ID         uuid.UUID           `db:"id"`
	X1         int                 `db:"x1"`
	Y1         int                 `db:"y1"`
	X2         int                 `db:"x2"`
	Y2         int                 `db:"y2"`
	WallTypeID uuid.UUID           `db:"wall_type_id"`
	FloorID    uuid.UUID           `db:"floor_id"`
	CreatedAt  pgtype.Timestamptz  `db:"created_at"`
	UpdatedAt  pgtype.Timestamptz  `db:"updated_at"`
	DeletedAt  *pgtype.Timestamptz `db:"deleted_at"`
}

type WallDetailed struct {
	Wall
	WallType WallType
}
