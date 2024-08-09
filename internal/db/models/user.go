package db

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	// Pointers to types like *string and *time.Time indicate nullable fields
	// Conventions
	// Primary Key: GORM uses a field named ID as the default primary key for each model.
	// Table Names: By default, GORM converts struct names to snake_case and pluralizes them for table names.
	// 		For instance, a User struct becomes users in the database.
	// Column Names: GORM automatically converts struct field names to snake_case for column names in the database.
	// Timestamp Fields: GORM uses fields named CreatedAt and UpdatedAt to automatically track the creation and update times of records.

	ID        uuid.UUID `gorm:"primaryKey; type:uuid"`
	Username  string    `gorm:"type:varchar(255); unique; not null"` // A regular string field
	Password  string    `gorm:"type:varchar(255); not null"`         // A regular string field
	CreatedAt time.Time // Automatically managed by GORM for update time
	UpdatedAt time.Time // Automatically managed by GORM for update time
	DeletedAt time.Time // Automatically managed by GORM for update time
	Roles     []Role    `gorm:"many2many:user_roles;"`
	Tokens    []RefreshToken
}
