package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// TODO Отвязать зависимость от pgtype

// User без пароля
type UserDTO struct {
	ID        uuid.UUID           `db:"id"`
	Username  string              `db:"username"`
	CreatedAt pgtype.Timestamptz  `db:"created_at"`
	UpdatedAt pgtype.Timestamptz  `db:"updated_at"`
	DeletedAt *pgtype.Timestamptz `db:"deleted_at"`
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
	Username string `db:"username"`
}

type PathUpdateUserDTO struct {
	ID       uuid.UUID `json:"id"`
	Username *string   `json:"username"`
	Password *string   `json:"password"`
}
