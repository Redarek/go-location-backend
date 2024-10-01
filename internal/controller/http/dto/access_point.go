package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessPointDTO struct {
	ID                uuid.UUID           `json:"id"`
	Name              string              `json:"name"`
	X                 *int                `json:"x"`
	Y                 *int                `json:"y"`
	Z                 *float64            `json:"z"`
	IsVirtual         bool                `json:"isVirtual"`
	AccessPointTypeID uuid.UUID           `json:"accessPointTypeId"`
	FloorID           uuid.UUID           `json:"floorId"`
	CreatedAt         pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt         pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt         *pgtype.Timestamptz `json:"deletedAt"`
}

type CreateAccessPointDTO struct {
	Name              string    `json:"name"`
	X                 *int      `json:"x"`
	Y                 *int      `json:"y"`
	Z                 *float64  `json:"z"`
	IsVirtual         bool      `json:"isVirtual"`
	AccessPointTypeID uuid.UUID `json:"accessPointTypeId"`
	FloorID           uuid.UUID `json:"floorId"`
}

type GetAccessPointsDTO struct {
	FloorID uuid.UUID `json:"floorId"`
	Page    int
	Size    int
}

type PatchUpdateAccessPointDTO struct {
	ID                uuid.UUID  `json:"id"`
	Name              *string    `json:"name"`
	X                 *int       `json:"x"`
	Y                 *int       `json:"y"`
	Z                 *float64   `json:"z"`
	IsVirtual         *bool      `json:"isVirtual"`
	AccessPointTypeID *uuid.UUID `json:"accessPointTypeId"`
}
