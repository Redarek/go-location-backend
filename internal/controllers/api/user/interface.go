package user

import (
	"context"

	"github.com/google/uuid"

	"location-backend/internal/models/user"
)

type Service interface {
	// GetUserByUsername(username string) (u User, err error)
	GetUserByUsername(ctx context.Context, username string) (user *user.User, err error)
	CreateUser(ctx context.Context, username, password string) (id uuid.UUID, err error)
}
