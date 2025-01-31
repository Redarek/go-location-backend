package dto

import (
	"github.com/google/uuid"
)

type RegisterUserDTO struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type CreateUserDTO struct {
	Username     string `db:"username"`
	PasswordHash string `db:"password"`
}

type LoginUserDTO struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type GetUserByNameDTO struct {
	Username string `json:"username" db:"username"`
}

type PatchUpdateUserDTO struct {
	ID       uuid.UUID `json:"id"`
	Username *string   `json:"username,omitempty" db:"username"`
	Password *string   `json:"password,omitempty" db:"password"`
}
