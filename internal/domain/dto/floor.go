package dto

import (
	"github.com/google/uuid"
)

type CreateFloorDTO struct {
	Name                 string    `json:"name" db:"name"`
	Number               int       `json:"number" db:"number"`
	Image                *string   `json:"image" db:"image"`
	WidthInPixels        int       `json:"widthInPixels" db:"width_in_pixels"`
	HeightInPixels       int       `json:"heightInPixels" db:"height_in_pixels"`
	Scale                float64   `json:"scale" db:"scale"`
	CellSizeMeter        float64   `json:"cellSizeMeter" db:"cell_size_meter"`
	NorthAreaIndentMeter float64   `json:"northAreaIndentMeter" db:"north_area_indent_meter"`
	SouthAreaIndentMeter float64   `json:"southAreaIndentMeter" db:"south_area_indent_meter"`
	WestAreaIndentMeter  float64   `json:"westAreaIndentMeter" db:"west_area_indent_meter"`
	EastAreaIndentMeter  float64   `json:"eastAreaIndentMeter" db:"east_area_indent_meter"`
	BuildingID           uuid.UUID `json:"buildingId" db:"building_id"`
}

type GetFloorDTO struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type PatchUpdateFloorDTO struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	Name                 *string   `json:"name" db:"name"`
	Number               *int      `json:"number" db:"number"`
	Image                *string   `json:"image" db:"image"`
	WidthInPixels        *int      `json:"widthInPixels" db:"width_in_pixels"`
	HeightInPixels       *int      `json:"heightInPixels" db:"height_in_pixels"`
	Scale                *float64  `json:"scale" db:"scale"`
	CellSizeMeter        *float64  `json:"cellSizeMeter" db:"cell_size_meter"`
	NorthAreaIndentMeter *float64  `json:"northAreaIndentMeter" db:"north_area_indent_meter"`
	SouthAreaIndentMeter *float64  `json:"southAreaIndentMeter" db:"south_area_indent_meter"`
	WestAreaIndentMeter  *float64  `json:"westAreaIndentMeter" db:"west_area_indent_meter"`
	EastAreaIndentMeter  *float64  `json:"eastAreaIndentMeter" db:"east_area_indent_meter"`
}

type GetFloorsDTO struct {
	BuildingID uuid.UUID `json:"buildingId" db:"building_id"`
	Page       int
	Size       int
}
