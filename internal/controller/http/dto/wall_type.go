package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type WallTypeDTO struct {
	ID            uuid.UUID           `json:"id"`
	Name          string              `json:"name"`
	Color         string              `json:"color"`
	Attenuation24 float64             `json:"attenuation24"`
	Attenuation5  float64             `json:"attenuation5"`
	Attenuation6  float64             `json:"attenuation6"`
	Thickness     float64             `json:"thickness"`
	SiteID        uuid.UUID           `json:"siteId"`
	CreatedAt     pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt     pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt     *pgtype.Timestamptz `json:"deletedAt"`
}

type CreateWallTypeDTO struct {
	Name          string    `json:"name"`
	Color         string    `json:"color"`
	Attenuation24 float64   `json:"attenuation24"`
	Attenuation5  float64   `json:"attenuation5"`
	Attenuation6  float64   `json:"attenuation6"`
	Thickness     float64   `json:"thickness"`
	SiteID        uuid.UUID `json:"siteId"`
}

type PatchUpdateWallTypeDTO struct {
	ID            uuid.UUID `db:"id"`
	Name          *string   `json:"name"`
	Color         *string   `json:"color"`
	Attenuation24 *float64  `json:"attenuation24"`
	Attenuation5  *float64  `json:"attenuation5"`
	Attenuation6  *float64  `json:"attenuation6"`
	Thickness     *float64  `json:"thickness"`
}

type GetWallTypesDTO struct {
	SiteID uuid.UUID `json:"siteId"`
	Page   int       `json:"page"`
	Size   int       `json:"size"`
}
