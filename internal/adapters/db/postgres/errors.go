package postgres

import "errors"

var (
	// Occurs when no desired entity in database
	ErrNotFound = errors.New("entity not found")

	// Occurs when desired entity was not updated
	ErrNotUpdated = errors.New("entity not updated")
)
