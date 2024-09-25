package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type WallTypeDTO struct {
	ID uuid.UUID `db:"id"`
	CreateWallTypeDTO
	CreatedAt pgtype.Timestamptz  `db:"created_at"`
	UpdatedAt pgtype.Timestamptz  `db:"updated_at"`
	DeletedAt *pgtype.Timestamptz `db:"deleted_at"`
}

type CreateWallTypeDTO struct {
	Name          string    `db:"name"`
	Color         string    `db:"color"`
	Attenuation24 float64   `db:"attenuation_24"`
	Attenuation5  float64   `db:"attenuation_5"`
	Attenuation6  float64   `db:"attenuation_6"`
	Thickness     float64   `db:"thickness"`
	SiteID        uuid.UUID `db:"site_id"`
}

type PatchUpdateWallTypeDTO struct {
	ID            uuid.UUID `db:"id"`
	Name          *string   `db:"name"`
	Color         *string   `db:"color"`
	Attenuation24 *float64  `db:"attenuation_24"`
	Attenuation5  *float64  `db:"attenuation_5"`
	Attenuation6  *float64  `db:"attenuation_6"`
	Thickness     *float64  `db:"thickness"`
}

type GetWallTypesDTO struct {
	SiteID uuid.UUID `db:"id"`
	Limit  int
	Offset int
}
