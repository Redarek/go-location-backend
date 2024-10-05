package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type WallDTO struct {
	ID         uuid.UUID           `json:"id"`
	X1         int                 `json:"x1"`
	Y1         int                 `json:"y1"`
	X2         int                 `json:"x2"`
	Y2         int                 `json:"y2"`
	WallTypeID uuid.UUID           `json:"wallTypeId"`
	FloorID    uuid.UUID           `json:"floorId"`
	CreatedAt  pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt  pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt  *pgtype.Timestamptz `json:"deletedAt"`
}

type WallDetailedDTO struct {
	WallDTO
	WallTypeDTO WallTypeDTO `json:"wallType"`
}

type CreateWallDTO struct {
	X1         int       `json:"x1"`
	Y1         int       `json:"y1"`
	X2         int       `json:"x2"`
	Y2         int       `json:"y2"`
	WallTypeID uuid.UUID `json:"wallTypeId"`
	FloorID    uuid.UUID `json:"floorId"`
}

type PatchUpdateWallDTO struct {
	ID         uuid.UUID  `json:"id"`
	X1         *int       `json:"x1"`
	Y1         *int       `json:"y1"`
	X2         *int       `json:"x2"`
	Y2         *int       `json:"y2"`
	WallTypeID *uuid.UUID `json:"wallTypeId"`
	// FloorID    uuid.UUID           `json:"floor_id"`
}

type GetWallsDTO struct {
	FloorID uuid.UUID `json:"floorId"`
	Page    int
	Size    int
}
