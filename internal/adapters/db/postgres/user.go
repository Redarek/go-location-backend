package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/entity"
)

type UserRepo interface {
	Create(userCreate entity.UserCreate) (userID uuid.UUID, err error)
	GetOneByName(username string) (user entity.User, err error)
	// GetOneByName(username string) entity.User
	// GetAll(limit, offset int) []entity.User
	// Delete(book entity.Book) error
}

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *userRepo {
	return &userRepo{pool: pool}
}

func (r *userRepo) Create(userCreate entity.UserCreate) (userID uuid.UUID, err error) {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	row := r.pool.QueryRow(context.Background(), query,
		userCreate.Username,
		userCreate.PasswordHash,
	)
	var user entity.User
	err = row.Scan(&user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to scan user")
		return uuid.UUID{}, err
	}

	return user.ID, nil
}

func (r *userRepo) GetOneByName(username string) (user entity.User, err error) {
	query := `SELECT 
		id, 
		username, 
		password, 
		created_at, 
		updated_at, 
		deleted_at 
	FROM users 
	WHERE username = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(context.Background(), query, username)
	err = row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msgf("user %v not found", username)
			return user, ErrNotFound
		}
		log.Error().Err(err).Msg("Failed when scanning user")
		return
	}
	log.Debug().Msgf("Retrieved user: %v", user)
	return
}

//	func (bs *userRepo) GetAll(limit, offset int) []*entity.User {
//		return nil
//	}

// func (bs *userRepo) Delete(user *entity.User) error {
// 	return nil
// }
