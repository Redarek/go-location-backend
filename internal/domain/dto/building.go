package dto

import (
	"github.com/google/uuid"
)

type CreateBuildingDTO struct {
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	Country     string    `json:"country" db:"country"`
	City        string    `json:"city" db:"city"`
	Address     string    `json:"address" db:"address"`
	SiteID      uuid.UUID `json:"siteId" db:"site_id"`
}

type GetBuildingDTO struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type GetBuildingsDTO struct {
	SiteID uuid.UUID `json:"siteId" db:"site_id"`
	Page   int
	Size   int
}

type PatchUpdateBuildingDTO struct {
	ID          uuid.UUID `db:"id"`
	Name        *string   `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	Country     *string   `json:"country" db:"country"`
	City        *string   `json:"city" db:"city"`
	Address     *string   `json:"address" db:"address"`
	// SiteID      *uuid.UUID          `db:"user_id"` // TODO Возможно позже стоит добавить
}
