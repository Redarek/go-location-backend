package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SensorTypeDTO struct {
	ID        uuid.UUID           `db:"id"`
	Name      string              `db:"name"`
	Model     string              `db:"model"`
	Color     string              `db:"color"`
	Z         float64             `db:"z"`
	IsVirtual bool                `db:"is_virtual"`
	SiteID    uuid.UUID           `db:"site_id"`
	CreatedAt pgtype.Timestamptz  `db:"created_at"`
	UpdatedAt pgtype.Timestamptz  `db:"updated_at"`
	DeletedAt *pgtype.Timestamptz `db:"deleted_at"`
}

// type SensorTypeDetailedDTO struct {
// 	SensorTypeDTO
// 	RadioTemplatesDTO []*SensorRadioTemplateDTO
// }

type CreateSensorTypeDTO struct {
	Name      string    `db:"name"`
	Model     string    `db:"model"`
	Color     string    `db:"color"`
	Z         float64   `db:"z"`
	IsVirtual bool      `db:"is_virtual"`
	SiteID    uuid.UUID `db:"site_id"`
}

type GetSensorTypesDTO struct {
	SiteID uuid.UUID `db:"site_id"`
	Limit  int
	Offset int
}

type GetSensorTypeDetailedDTO struct {
	ID     uuid.UUID `db:"id"`
	Limit  int
	Offset int
}

type PatchUpdateSensorTypeDTO struct {
	ID        uuid.UUID `db:"id"`
	Name      *string   `db:"name"`
	Model     *string   `db:"model"`
	Color     *string   `db:"color"`
	Z         *float64  `db:"z"`
	IsVirtual *bool     `db:"is_virtual"`
	// SiteID      *uuid.UUID          `db:"user_id"` // TODO Возможно позже стоит добавить
}
