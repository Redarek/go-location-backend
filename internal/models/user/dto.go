package user

type CreateUserDTO struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
