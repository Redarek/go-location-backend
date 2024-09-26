package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type WallDTO struct {
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

type WallDetailedDTO struct {
	WallDTO
	WallTypeDTO WallTypeDTO
}

type CreateWallDTO struct {
	X1         int       `db:"x1"`
	Y1         int       `db:"y1"`
	X2         int       `db:"x2"`
	Y2         int       `db:"y2"`
	WallTypeID uuid.UUID `db:"wall_type_id"`
	FloorID    uuid.UUID `db:"floor_id"`
}

type PatchUpdateWallDTO struct {
	ID         uuid.UUID  `db:"id"`
	X1         *int       `db:"x1"`
	Y1         *int       `db:"y1"`
	X2         *int       `db:"x2"`
	Y2         *int       `db:"y2"`
	WallTypeID *uuid.UUID `db:"wall_type_id"`
	// FloorID    uuid.UUID           `db:"floor_id"`
}

type GetWallsDTO struct {
	FloorID uuid.UUID `db:"floor_id"`
	Limit   int
	Offset  int
}
