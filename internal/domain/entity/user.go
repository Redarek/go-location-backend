package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID           uuid.UUID           `db:"id"`
	Username     string              `db:"username"`
	PasswordHash string              `db:"password"`
	CreatedAt    pgtype.Timestamptz  `db:"created_at"`
	UpdatedAt    pgtype.Timestamptz  `db:"updated_at"`
	DeletedAt    *pgtype.Timestamptz `db:"deleted_at"`
}

// type UserView struct {

// }
