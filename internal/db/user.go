package db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// CreateUser creates a user
func (p *postgres) CreateUser(username, password string) (id uuid.UUID, err error) {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, username, password)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
	}
	return
}

// GetUserByUsername retrieves a user
func (p *postgres) GetUserByUsername(username string) (user User, err error) {
	query := `SELECT * FROM users WHERE username = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, username)
	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No user found with username %v", username)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve user")
		return
	}
	log.Debug().Msgf("Retrieved user: %v", user)
	return
}

// CreateRefreshToken creates a refresh token
func (p *postgres) CreateRefreshToken(rt RefreshToken) (id uuid.UUID, err error) {
	query := `INSERT INTO refresh_tokens (token, user_id) VALUES ($1, $2) RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, rt.Token, rt.UserID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create token")
	}
	return
}
