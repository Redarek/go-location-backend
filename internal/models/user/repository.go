package user

type Repository interface {
	Create(user *User) error
	GetOne(username string) (*User, error)
	// GetAll(limit int, offset int) (*[]User, error)
	// Delete(user *User) (error)
}
