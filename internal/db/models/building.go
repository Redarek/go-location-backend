package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Building struct {
	ID          uuid.UUID          `json:"id" db:"id"`
	Name        string             `json:"name" db:"name"`
	Description string             `json:"description" db:"description"`
	Country     string             `json:"country" db:"country"`
	City        string             `json:"city" db:"city"`
	Address     string             `json:"address" db:"address"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
	DeletedAt   pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
	SiteID      uuid.UUID          `json:"siteId" db:"site_id"`
	Floors      []*Floor           `json:"floors"`
}
