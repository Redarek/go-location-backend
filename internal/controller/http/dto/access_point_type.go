package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessPointTypeDTO struct {
	ID        uuid.UUID           `json:"id"`
	Name      string              `json:"name"`
	Model     string              `json:"model"`
	Color     string              `json:"color"`
	Z         float64             `json:"z"`
	IsVirtual bool                `json:"isVirtual"`
	SiteID    uuid.UUID           `json:"site_id"`
	CreatedAt pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt *pgtype.Timestamptz `json:"deletedAt"`
}

type CreateAccessPointTypeDTO struct {
	Name      string    `json:"name"`
	Model     string    `json:"model"`
	Color     string    `json:"color"`
	Z         float64   `json:"z"`
	IsVirtual bool      `json:"isVirtual"`
	SiteID    uuid.UUID `json:"siteId"`
}

type GetAccessPointTypesDTO struct {
	SiteID uuid.UUID `json:"site_id"`
	Page   int
	Size   int
}

type PatchUpdateAccessPointTypeDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   `json:"name"`
	Model     *string   `json:"model"`
	Color     *string   `json:"color"`
	Z         *float64  `json:"z"`
	IsVirtual *bool     `json:"isVirtual"`
	// SiteID      *uuid.UUID          `json:"user_id"` // TODO Возможно позже стоит добавить
}
