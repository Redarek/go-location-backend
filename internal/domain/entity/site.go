package entity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Site struct {
	ID          uuid.UUID           `json:"id" db:"id"`
	Name        string              `json:"name" db:"name"`
	Description *string             `json:"description" db:"description"`
	UserID      uuid.UUID           `db:"user_id"`
	CreatedAt   pgtype.Timestamptz  `json:"createdAt" db:"created_at"`
	UpdatedAt   pgtype.Timestamptz  `json:"updatedAt" db:"updated_at"`
	DeletedAt   *pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
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
