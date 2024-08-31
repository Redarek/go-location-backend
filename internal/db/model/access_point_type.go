package model

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessPointType struct {
	ID        uuid.UUID          `json:"id" db:"id"`
	Name      string             `json:"name" db:"name"`
	Model     string             `json:"model" db:"model"`
	Color     string             `json:"color" db:"color"`
	Z         int                `json:"z" db:"z"`
	SiteID    uuid.UUID          `json:"siteId" db:"site_id"`
	CreatedAt pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
}
