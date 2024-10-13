package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SiteDTO struct {
	ID          uuid.UUID           `db:"id"`
	Name        string              `db:"name"`
	Description *string             `db:"description"`
	UserID      uuid.UUID           `db:"user_id"`
	CreatedAt   pgtype.Timestamptz  `db:"created_at"`
	UpdatedAt   pgtype.Timestamptz  `db:"updated_at"`
	DeletedAt   *pgtype.Timestamptz `db:"deleted_at"`
}

type CreateSiteDTO struct {
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	UserID      uuid.UUID `db:"user_id"`
}

type GetSiteDetailedDTO struct {
	ID     uuid.UUID `db:"id"`
	Limit  int
	Offset int
}

type GetSitesDTO struct {
	UserID uuid.UUID `db:"user_id"`
	Limit  int
	Offset int
}

// ? Возможно удалить
// type GetSitesDetailedDTO struct {
// 	UserID uuid.UUID `db:"user_id"`
// 	Limit  int
// 	Offset int
// }

type PatchUpdateSiteDTO struct {
	ID          uuid.UUID `db:"id"`
	Name        *string   `db:"name"`
	Description *string   `db:"description"`
	// UserID      *uuid.UUID          `db:"user_id"` // TODO Возможно позже стоит добавить
}

// type SoftDeleteSiteDTO struct {
// 	ID uuid.UUID `db:"id"`
// }
