package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Site struct {
	ID          uuid.UUID           `db:"id"`
	Name        string              `db:"name"`
	Description *string             `db:"description"`
	UserID      uuid.UUID           `db:"user_id"`
	CreatedAt   pgtype.Timestamptz  `db:"created_at"`
	UpdatedAt   pgtype.Timestamptz  `db:"updated_at"`
	DeletedAt   *pgtype.Timestamptz `db:"deleted_at"`
}

type SiteDetailed struct {
	Site
	Buildings        []*Building
	WallTypes        []*WallType
	AccessPointTypes []*AccessPointType
	SensorTypes      []*SensorType
}

// type SiteView struct {
// 	Site
// }
