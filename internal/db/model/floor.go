package model

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Floor struct {
	ID             uuid.UUID          `json:"id" db:"id"`
	Name           *string            `json:"name" db:"name"`
	Number         *int               `json:"number" db:"number"`
	Image          *string            `json:"image" db:"image"`
	Heatmap        *string            `json:"heatmap" db:"heatmap"`
	WidthInPixels  *int               `json:"widthInPixels" db:"width_in_pixels"`
	HeightInPixels *int               `json:"heightInPixels" db:"height_in_pixels"`
	Scale          *float64           `json:"scale" db:"scale"`
	BuildingID     uuid.UUID          `json:"buildingId" db:"building_id"`
	CreatedAt      pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
	DeletedAt      pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
}
