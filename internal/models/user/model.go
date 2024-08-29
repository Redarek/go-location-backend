package user

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        uuid.UUID          `json:"id" db:"id"`
	Username  string             `json:"username" db:"username"`
	Password  string             `json:"password" db:"password"`
	CreatedAt pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
}
