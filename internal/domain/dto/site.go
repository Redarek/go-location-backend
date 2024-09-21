package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SiteDTO struct {
	ID          uuid.UUID          `db:"id"`
	Name        string             `db:"name"`
	Description *string            `db:"description"`
	UserID      uuid.UUID          `db:"user_id"`
	CreatedAt   pgtype.Timestamptz `db:"created_at"`
	UpdatedAt   pgtype.Timestamptz `db:"updated_at"`
	DeletedAt   pgtype.Timestamptz `db:"deleted_at"`
}

type CreateSiteDTO struct {
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	UserID      uuid.UUID `db:"user_id"`
}

type GetSiteDTO struct {
	ID uuid.UUID `db:"id"`
}
