package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// TODO Отвязать зависимость от pgtype

// User без пароля
type UserDTO struct {
	ID           uuid.UUID           `json:"id"`
	Username     string              `json:"username"`
	PasswordHash string              `json:"-"`
	CreatedAt    pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt    *pgtype.Timestamptz `json:"deletedAt"`
}

type RegisterUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetUserByNameDTO struct {
	Username string `json:"username"`
}

type PatchUpdateUserDTO struct {
	ID       uuid.UUID `json:"id"`
	Username *string   `json:"username,omitempty"`
	Password *string   `json:"password,omitempty"`
}
