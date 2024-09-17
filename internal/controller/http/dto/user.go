package dto

import "github.com/google/uuid"

type RegisterUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PathUpdateUserDTO struct {
	ID       uuid.UUID `json:"id"`
	Username *string   `json:"username"`
	Password *string   `json:"password"`
}
