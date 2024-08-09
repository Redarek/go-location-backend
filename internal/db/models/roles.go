package db

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid"`
	Name      string    `gorm:"type:varchar(255);unique"` // A regular string field
	CreatedAt time.Time // Automatically managed by GORM for update time
	UpdatedAt time.Time // Automatically managed by GORM for update time
	DeletedAt time.Time // Automatically managed by GORM for update time
	Users     []User    `gorm:"many2many:user_roles;"`
}
