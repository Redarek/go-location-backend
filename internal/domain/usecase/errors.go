package usecase

import "errors"

var (
	// Occurs when entity not found
	ErrNotFound = errors.New("not found")

	// Occurs when entity already exist
	ErrAlreadyExists = errors.New("already exists")
)
