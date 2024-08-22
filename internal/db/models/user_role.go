package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRole struct {
	ID        uuid.UUID          `json:"id" db:"id"`
	UserID    uuid.UUID          `json:"userId" db:"user_id"`
	RoleID    uuid.UUID          `json:"roleId" db:"role_id"`
	CreatedAt pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
}
