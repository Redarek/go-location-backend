package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Role struct {
	ID        uuid.UUID           `db:"id"`
	Name      string              `db:"name"`
	CreatedAt pgtype.Timestamptz  `db:"created_at"`
	UpdatedAt pgtype.Timestamptz  `db:"updated_at"`
	DeletedAt *pgtype.Timestamptz `db:"deleted_at"`
}
