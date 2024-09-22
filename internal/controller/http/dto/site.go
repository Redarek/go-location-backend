package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SiteDTO struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description *string            `json:"description"`
	UserID      uuid.UUID          `json:"userId"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
	DeletedAt   pgtype.Timestamptz `json:"deletedAt"`
}

type CreateSiteDTO struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type GetSiteDTO struct {
	ID uuid.UUID `json:"id"`
}

type GetSitesDTO struct {
	UserID uuid.UUID `json:"id"`
	Page   int       `json:"page"`
	Size   int       `json:"size"`
}

type PatchUpdateSiteDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description"`
	// UserID      *uuid.UUID          `db:"user_id"` // TODO Возможно позже стоит добавить
}

// type SoftDeleteSiteDTO struct {
// 	ID uuid.UUID `json:"id"`
// }
