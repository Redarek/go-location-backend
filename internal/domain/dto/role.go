package dto

import (
	"github.com/google/uuid"
)

type CreateRoleDTO struct {
	Name string `json:"name" db:"name"`
}

type GetRoleDTO struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type GetRoleByNameDTO struct {
	Name string `json:"name"`
}
