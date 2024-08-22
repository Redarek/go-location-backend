package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RefreshToken struct {
	ID     uuid.UUID          `json:"id" db:"id"`
	Token  string             `json:"token" db:"token"`
	Expiry pgtype.Timestamptz `json:"expiry" db:"expiry"`
	UserID uuid.UUID          `json:"userId" db:"user_id"`
}
