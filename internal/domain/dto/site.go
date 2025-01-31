package dto

import (
	"github.com/google/uuid"
)

type CreateSiteDTO struct {
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	UserID      uuid.UUID `json:"userId" db:"user_id"`
}

// ? нужно ли DTO из одного элемента?
type GetSiteDTO struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type GetSiteDetailedDTO struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Page int
	Size int
}

type GetSitesDTO struct {
	UserID uuid.UUID `json:"userId" db:"user_id"`
	Page   int
	Size   int
}

// ? Возможно удалить
// type GetSitesDetailedDTO struct {
// 	UserID uuid.UUID `db:"user_id"`
// 	Limit  int
// 	Offset int
// }

type PatchUpdateSiteDTO struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        *string   `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	// UserID      *uuid.UUID          `db:"user_id"` // TODO Возможно позже стоит добавить
}

// type SoftDeleteSiteDTO struct {
// 	ID uuid.UUID `db:"id"`
// }
