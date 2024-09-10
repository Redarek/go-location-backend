package user_usecase

import "github.com/google/uuid"

// That DTO is for avoiding HTTP DTO conflicts

type CreateUserDTO struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type PathUpdateUserDTO struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username *string   `json:"username" db:"username"`
	Password *string   `json:"password" db:"password"`
}
