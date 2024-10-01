package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccessPointType struct {
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

type AccessPointTypeDetailed struct {
	AccessPointType
	RadioTemplates []*AccessPointRadioTemplate
}
