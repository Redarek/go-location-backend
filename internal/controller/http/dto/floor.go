package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type FloorDTO struct {
	ID uuid.UUID `json:"id"`
	CreateFloorDTO
	CreatedAt pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt *pgtype.Timestamptz `json:"deletedAt"`
}

type CreateFloorDTO struct {
	Name                 string    `json:"name"`
	Number               int       `json:"number"`
	Image                *string   `json:"image"`
	Heatmap              *string   `json:"heatmap"`
	WidthInPixels        int       `json:"widthInPixels"`
	HeightInPixels       int       `json:"heightInPixels"`
	Scale                float64   `json:"scale"`
	CellSizeMeter        float64   `json:"cellSizeMeter"`
	NorthAreaIndentMeter float64   `json:"northAreaIndentMeter"`
	SouthAreaIndentMeter float64   `json:"southAreaIndentMeter"`
	WestAreaIndentMeter  float64   `json:"westAreaIndentMeter"`
	EastAreaIndentMeter  float64   `json:"eastAreaIndentMeter"`
	BuildingID           uuid.UUID `json:"buildingId"`
}

type PatchUpdateFloorDTO struct {
	ID                   uuid.UUID `json:"id"`
	Name                 *string   `json:"name"`
	Number               *int      `json:"number"`
	Image                *string   `json:"image"`
	Heatmap              *string   `json:"heatmap"`
	WidthInPixels        *int      `json:"widthInPixels"`
	HeightInPixels       *int      `json:"heightInPixels"`
	Scale                *float64  `json:"scale"`
	CellSizeMeter        *float64  `json:"cellSizeMeter"`
	NorthAreaIndentMeter *float64  `json:"northAreaIndentMeter"`
	SouthAreaIndentMeter *float64  `json:"southAreaIndentMeter"`
	WestAreaIndentMeter  *float64  `json:"westAreaIndentMeter"`
	EastAreaIndentMeter  *float64  `json:"eastAreaIndentMeter"`
	// BuildingID           uuid.UUID `json:"buildingId" db:"building_id"`
}

type GetFloorsDTO struct {
	BuildingID uuid.UUID `json:"buildingId"`
	Page       int       `json:"page"`
	Size       int       `json:"size"`
}
