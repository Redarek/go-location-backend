package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SiteDTO struct {
	ID          uuid.UUID           `json:"id"`
	Name        string              `json:"name"`
	Description *string             `json:"description"`
	UserID      uuid.UUID           `json:"userId"`
	CreatedAt   pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt   *pgtype.Timestamptz `json:"deletedAt"`
}

type SiteDetailedDTO struct {
	SiteDTO
	Buildings        []*BuildingDTO
	WallTypes        []*WallTypeDTO
	AccessPointTypes []*AccessPointTypeDTO
	SensorTypes      []*SensorTypeDTO
}

type CreateSiteDTO struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type GetSiteDTO struct {
	ID uuid.UUID `json:"id"`
}

// ? Возможно удалить
type GetSiteDetailedDTO struct {
	ID   uuid.UUID `json:"id"`
	Page int       `json:"page"`
	Size int       `json:"size"`
}

type GetSitesDTO struct {
	UserID uuid.UUID `json:"userId"`
	Page   int       `json:"page"`
	Size   int       `json:"size"`
}

// ? Возможно удалить
type GetSitesDetailedDTO struct {
	UserID uuid.UUID `json:"userId"`
	Page   int       `json:"page"`
	Size   int       `json:"size"`
}

type PatchUpdateSiteDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	// UserID      *uuid.UUID          `db:"user_id"` // TODO Возможно позже стоит добавить
}

// type SoftDeleteSiteDTO struct {
// 	ID uuid.UUID `json:"id"`
// }
