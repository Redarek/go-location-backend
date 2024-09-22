package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Role struct {
	ID        uuid.UUID           `json:"id" db:"id"`
	Name      string              `json:"name" db:"name"`
	CreatedAt pgtype.Timestamptz  `json:"createdAt" db:"created_at"`
	UpdatedAt pgtype.Timestamptz  `json:"updatedAt" db:"updated_at"`
	DeletedAt *pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
}
