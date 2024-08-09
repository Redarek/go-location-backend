package db

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid"`
	Token     string    `gorm:"type:varchar(1500); not null"`
	Expiry    time.Time `gorm:"type:timestamp; not null"`
	UserID    uuid.UUID `gorm:"type:uuid; not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
