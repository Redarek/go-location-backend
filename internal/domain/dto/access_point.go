package dto

import (
	"github.com/google/uuid"
)

type CreateAccessPointDTO struct {
	Name              string    `json:"name" db:"name"`
	Color             *string   `json:"color" db:"color"`
	X                 *int      `json:"x" db:"x"`
	Y                 *int      `json:"y" db:"y"`
	Z                 *float64  `json:"z" db:"z"`
	IsVirtual         bool      `json:"isVirtual" db:"is_virtual"`
	AccessPointTypeID uuid.UUID `json:"accessPointTypeId" db:"access_point_type_id"`
	FloorID           uuid.UUID `json:"floorId" db:"floor_id"`
}

type GetAccessPointDTO struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type GetAccessPointDetailedDTO struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Page int
	Size int
}

type GetAccessPointsDTO struct {
	FloorID uuid.UUID `json:"floorId" db:"floor_id"`
	Page    int
	Size    int
}

type GetAccessPointsDetailedDTO struct {
	FloorID uuid.UUID `json:"floorId" db:"floor_id"`
	Page    int
	Size    int
	//? Нужно ли добавить возможность выбирать лимит для каждой составляющей?
}

type PatchUpdateAccessPointDTO struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	Name              *string    `json:"name" db:"name"`
	Color             *string    `json:"color" db:"color"`
	X                 *int       `json:"x" db:"x"`
	Y                 *int       `json:"y" db:"y"`
	Z                 *float64   `json:"z" db:"z"`
	IsVirtual         *bool      `json:"isVirtual" db:"is_virtual"`
	AccessPointTypeID *uuid.UUID `json:"accessPointTypeId" db:"access_point_type_id"`
}
