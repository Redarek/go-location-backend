package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoleDTO struct {
	ID        uuid.UUID           `json:"id"`
	Name      string              `json:"name"`
	CreatedAt pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt *pgtype.Timestamptz `json:"deletedAt"`
}

type CreateRoleDTO struct {
	Name string `json:"name"`
}

type GetRoleByNameDTO struct {
	Name string `json:"name"`
}

type GetRoleDTO struct {
	ID uuid.UUID `json:"id"`
}
