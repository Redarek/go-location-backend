package dto

import "github.com/google/uuid"

type CreateRoleDTO struct {
	Name string `db:"name"`
}

type GetRoleByNameDTO struct {
	Name string `db:"name"`
}

type GetRoleDTO struct {
	ID uuid.UUID `db:"id"`
}
