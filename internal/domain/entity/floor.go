package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Floor struct {
	ID                   uuid.UUID           `json:"id" db:"id"`
	Name                 string              `json:"name" db:"name"`
	Number               int                 `json:"number" db:"number"`
	Image                *string             `json:"image" db:"image"`
	Heatmap              *string             `json:"heatmap" db:"heatmap"`
	WidthInPixels        int                 `json:"widthInPixels" db:"width_in_pixels"`
	HeightInPixels       int                 `json:"heightInPixels" db:"height_in_pixels"`
	Scale                float64             `json:"scale" db:"scale"`
	CellSizeMeter        float64             `json:"cellSizeMeter" db:"cell_size_meter"`
	NorthAreaIndentMeter float64             `json:"northAreaIndentMeter" db:"north_area_indent_meter"`
	SouthAreaIndentMeter float64             `json:"southAreaIndentMeter" db:"south_area_indent_meter"`
	WestAreaIndentMeter  float64             `json:"westAreaIndentMeter" db:"west_area_indent_meter"`
	EastAreaIndentMeter  float64             `json:"eastAreaIndentMeter" db:"east_area_indent_meter"`
	BuildingID           uuid.UUID           `json:"buildingId" db:"building_id"`
	CreatedAt            pgtype.Timestamptz  `json:"createdAt" db:"created_at"`
	UpdatedAt            pgtype.Timestamptz  `json:"updatedAt" db:"updated_at"`
	DeletedAt            *pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
}

type FloorDetailed struct {
	Floor
	AccessPoints []*AccessPointDetailed `json:"accessPoints"`
	Walls        []*WallDetailed        `json:"walls"`
	Sensors      []*Sensor              `json:"sensors"`
}
