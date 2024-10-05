package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SensorTypeDTO struct {
	ID        uuid.UUID           `json:"id"`
	Name      string              `json:"name"`
	Model     string              `json:"model"`
	Color     string              `json:"color"`
	Z         float64             `json:"z"`
	IsVirtual bool                `json:"isVirtual"`
	SiteID    uuid.UUID           `json:"siteId"`
	CreatedAt pgtype.Timestamptz  `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz  `json:"updatedAt"`
	DeletedAt *pgtype.Timestamptz `json:"deletedAt"`
}

// type SensorTypeDetailedDTO struct {
// 	SensorTypeDTO
// 	RadioTemplatesDTO []*SensorRadioTemplateDTO `json:"radioTemplates"`
// }

type CreateSensorTypeDTO struct {
	Name      string    `json:"name"`
	Model     string    `json:"model"`
	Color     string    `json:"color"`
	Z         float64   `json:"z"`
	IsVirtual bool      `json:"isVirtual"`
	SiteID    uuid.UUID `json:"siteId"`
}

type GetSensorTypesDTO struct {
	SiteID uuid.UUID `json:"siteId"`
	Page   int
	Size   int
}

type GetSensorTypeDetailedDTO struct {
	ID   uuid.UUID `json:"id"`
	Page int
	Size int
}

type PatchUpdateSensorTypeDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   `json:"name"`
	Model     *string   `json:"model"`
	Color     *string   `json:"color"`
	Z         *float64  `json:"z"`
	IsVirtual *bool     `json:"isVirtual"`
	// SiteID      *uuid.UUID          `json:"user_id"` // TODO Возможно позже стоит добавить
}
