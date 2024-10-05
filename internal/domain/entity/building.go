package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Building struct {
	ID          uuid.UUID           `db:"id"`
	Name        string              `db:"name"`
	Description *string             `db:"description"`
	Country     string              `db:"country"`
	City        string              `db:"city"`
	Address     string              `db:"address"`
	SiteID      uuid.UUID           `db:"site_id"`
	CreatedAt   pgtype.Timestamptz  `db:"created_at"`
	UpdatedAt   pgtype.Timestamptz  `db:"updated_at"`
	DeletedAt   *pgtype.Timestamptz `db:"deleted_at"`
}

type BuildingDetailed struct {
	Building
	Floors []*Floor
}
