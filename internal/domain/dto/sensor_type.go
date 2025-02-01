package dto

import (
	"github.com/google/uuid"
)

type CreateSensorTypeDTO struct {
	Name      string    `json:"name" db:"name"`
	Model     string    `json:"model" db:"model"`
	Color     string    `json:"color" db:"color"`
	Z         float64   `json:"z" db:"z"`
	IsVirtual bool      `json:"isVirtual" db:"is_virtual"`
	SiteID    uuid.UUID `json:"siteId" db:"site_id"`
}

type GetSensorTypeDTO struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type GetSensorTypesDTO struct {
	SiteID uuid.UUID `json:"siteId" db:"site_id"`
	Page   int
	Size   int
}

type GetSensorTypeDetailedDTO struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Page int
	Size int
}

type PatchUpdateSensorTypeDTO struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      *string   `json:"name" db:"name"`
	Model     *string   `json:"model" db:"model"`
	Color     *string   `json:"color" db:"color"`
	Z         *float64  `json:"z" db:"z"`
	IsVirtual *bool     `json:"isVirtual" db:"is_virtual"`
	// SiteID      *uuid.UUID          `db:"user_id"` // TODO Возможно позже стоит добавить
}
