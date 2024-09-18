package postgres

import "errors"

// ErrNotFound occurs when no desired entity in database
var ErrNotFound = errors.New("entity not found")
