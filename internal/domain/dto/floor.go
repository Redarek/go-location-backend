package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type FloorDTO struct {
	ID                   uuid.UUID           `db:"id"`
	Name                 string              `db:"name"`
	Number               int                 `db:"number"`
	Image                *string             `db:"image"`
	Heatmap              *string             `db:"heatmap"`
	WidthInPixels        int                 `db:"width_in_pixels"`
	HeightInPixels       int                 `db:"height_in_pixels"`
	Scale                float64             `db:"scale"`
	CellSizeMeter        float64             `db:"cell_size_meter"`
	NorthAreaIndentMeter float64             `db:"north_area_indent_meter"`
	SouthAreaIndentMeter float64             `db:"south_area_indent_meter"`
	WestAreaIndentMeter  float64             `db:"west_area_indent_meter"`
	EastAreaIndentMeter  float64             `db:"east_area_indent_meter"`
	BuildingID           uuid.UUID           `db:"building_id"`
	CreatedAt            pgtype.Timestamptz  `db:"createdAt"`
	UpdatedAt            pgtype.Timestamptz  `db:"updatedAt"`
	DeletedAt            *pgtype.Timestamptz `db:"deletedAt"`
}

type CreateFloorDTO struct {
	Name                 string    `db:"name"`
	Number               int       `db:"number"`
	Image                *string   `db:"image"`
	Heatmap              *string   `db:"heatmap"`
	WidthInPixels        int       `db:"width_in_pixels"`
	HeightInPixels       int       `db:"height_in_pixels"`
	Scale                float64   `db:"scale"`
	CellSizeMeter        float64   `db:"cell_size_meter"`
	NorthAreaIndentMeter float64   `db:"north_area_indent_meter"`
	SouthAreaIndentMeter float64   `db:"south_area_indent_meter"`
	WestAreaIndentMeter  float64   `db:"west_area_indent_meter"`
	EastAreaIndentMeter  float64   `db:"east_area_indent_meter"`
	BuildingID           uuid.UUID `db:"building_id"`
}

type PatchUpdateFloorDTO struct {
	ID                   uuid.UUID `db:"id"`
	Name                 *string   `db:"name"`
	Number               *int      `db:"number"`
	Image                *string   `db:"image"`
	Heatmap              *string   `db:"heatmap"`
	WidthInPixels        *int      `db:"width_in_pixels"`
	HeightInPixels       *int      `db:"height_in_pixels"`
	Scale                *float64  `db:"scale"`
	CellSizeMeter        *float64  `db:"cell_size_meter"`
	NorthAreaIndentMeter *float64  `db:"north_area_indent_meter"`
	SouthAreaIndentMeter *float64  `db:"south_area_indent_meter"`
	WestAreaIndentMeter  *float64  `db:"west_area_indent_meter"`
	EastAreaIndentMeter  *float64  `db:"east_area_indent_meter"`
	// BuildingID           uuid.UUID `db:"buildingId" db:"building_id"`
}

type GetFloorsDTO struct {
	BuildingID uuid.UUID `db:"id"`
	Limit      int
	Offset     int
}
