package dto

import (
	"github.com/google/uuid"
)

type CreateWallDTO struct {
	X1         int       `json:"x1" db:"x1"`
	Y1         int       `json:"y1" db:"y1"`
	X2         int       `json:"x2" db:"x2"`
	Y2         int       `json:"y2" db:"y2"`
	WallTypeID uuid.UUID `json:"wallTypeId" db:"wall_type_id"`
	FloorID    uuid.UUID `json:"floorId" db:"floor_id"`
}

type GetWallDTO struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type PatchUpdateWallDTO struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	X1         *int       `json:"x1" db:"x1"`
	Y1         *int       `json:"y1" db:"y1"`
	X2         *int       `json:"x2" db:"x2"`
	Y2         *int       `json:"y2" db:"y2"`
	WallTypeID *uuid.UUID `json:"wallTypeId" db:"wall_type_id"`
}

type GetWallsDTO struct {
	FloorID uuid.UUID `json:"floorId" db:"floor_id"`
	Page    int
	Size    int
}
