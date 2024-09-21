package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// TODO make case
type CreateRefreshTokenDTO struct {
	Token  string             `json:"token"`
	Expiry pgtype.Timestamptz `json:"expiry"`
	UserID uuid.UUID          `json:"userId"`
}
