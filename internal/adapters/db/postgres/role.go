package postgres

// import (
// 	"context"
// 	"errors"

// 	"github.com/google/uuid"
// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/rs/zerolog/log"

// 	"location-backend/internal/domain/entity"
// )

// type RoleRepo interface {
// 	Create(userCreate entity.UserCreate) (userID uuid.UUID, err error)
// }

// type roleRepo struct {
// 	pool *pgxpool.Pool
// }

// func NewRoleRepo(pool *pgxpool.Pool) *roleRepo {
// 	return &roleRepo{pool: pool}
// }
