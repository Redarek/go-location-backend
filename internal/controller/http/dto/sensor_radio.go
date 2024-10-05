package dto

import (
	// "encoding/json"

	"github.com/google/uuid"
	// "github.com/jackc/pgx/v5/pgtype"
)

type SensorRadioDTO struct {
	ID uuid.UUID `json:"id"`
}
