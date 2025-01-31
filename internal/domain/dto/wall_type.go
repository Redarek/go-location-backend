package dto

import (
	"github.com/google/uuid"
)

type CreateWallTypeDTO struct {
	Name          string    `json:"name" db:"name"`
	Color         string    `json:"color" db:"color"`
	Attenuation24 float64   `json:"attenuation24" db:"attenuation_24"`
	Attenuation5  float64   `json:"attenuation5" db:"attenuation_5"`
	Attenuation6  float64   `json:"attenuation6" db:"attenuation_6"`
	Thickness     float64   `json:"thickness" db:"thickness"`
	SiteID        uuid.UUID `json:"siteId" db:"site_id"`
}

type GetWallTypeDTO struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type PatchUpdateWallTypeDTO struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Name          *string   `json:"name" db:"name"`
	Color         *string   `json:"color" db:"color"`
	Attenuation24 *float64  `json:"attenuation24" db:"attenuation_24"`
	Attenuation5  *float64  `json:"attenuation5" db:"attenuation_5"`
	Attenuation6  *float64  `json:"attenuation6" db:"attenuation_6"`
	Thickness     *float64  `json:"thickness" db:"thickness"`
}

type GetWallTypesDTO struct {
	SiteID uuid.UUID `json:"siteId" db:"site_id"`
	Page   int
	Size   int
}
