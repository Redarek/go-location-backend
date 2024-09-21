package dto

type CreateUserDTO struct {
	Username     string `db:"username"`
	PasswordHash string `db:"password"`
}
