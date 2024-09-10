package dto

import "github.com/google/uuid"

type CreateUserDTO struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type PathUpdateUserDTO struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username *string   `json:"username" db:"username"`
	Password *string   `json:"password" db:"password"`
}
