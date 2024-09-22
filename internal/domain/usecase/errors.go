package usecase

import "errors"

var (
	// Occurs when entity not found
	ErrNotFound = errors.New("not found")

	// Occurs when entity already exist
	ErrAlreadyExists = errors.New("already exists")

	// Occurs when entity already soft deleted
	ErrAlreadySoftDeleted = errors.New("already soft deleted")

	// Occurs when desired entity was not updated
	ErrNotUpdated = errors.New("entity not updated")
)
