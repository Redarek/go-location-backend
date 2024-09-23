package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type BuildingDTO struct {
	ID          uuid.UUID           `json:"id" db:"id"`
	Name        string              `json:"name" db:"name"`
	Description *string             `json:"description"`
	Country     string              `json:"country"`
	City        string              `json:"city"`
	Address     string              `json:"address"`
	SiteID      uuid.UUID           `json:"siteId"`
	CreatedAt   pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt   *pgtype.Timestamptz `json:"deletedAt"`
}

type CreateBuildingDTO struct {
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	Address     string    `json:"address"`
	SiteID      uuid.UUID `json:"siteId"`
}

type PatchUpdateBuildingDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Country     *string   `json:"country"`
	City        *string   `json:"city"`
	Address     *string   `json:"address"`
	// SiteID      *uuid.UUID          `db:"user_id"` // TODO Возможно позже стоит добавить
}

type GetBuildingDTO struct {
	ID uuid.UUID `json:"id"`
}

type GetBuildingsDTO struct {
	SiteID uuid.UUID `json:"siteId"`
	Page   int       `json:"page"`
	Size   int       `json:"size"`
}
