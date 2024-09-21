package dto

import "github.com/google/uuid"

type RegisterUserDTO struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

type CreateUserDTO struct {
	Username     string `db:"username"`
	PasswordHash string `db:"password"`
}

type LoginUserDTO struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

type GetUserByNameDTO struct {
	Username string `db:"username"`
}

type PathUpdateUserDTO struct {
	ID       uuid.UUID `db:"id"`
	Username *string   `db:"username"`
	Password *string   `db:"password"`
}
