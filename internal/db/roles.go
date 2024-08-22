package db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	. "location-backend/internal/db/models"
)

// Creates a role
func (p *postgres) CreateRole(name string) (id uuid.UUID, err error) {
	query := `INSERT INTO roles (name) VALUES ($1) RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, name)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create role")
	}
	return
}

// Retrieves a role
func (p *postgres) GetRoleByName(name string) (u User, err error) {
	query := `SELECT * FROM roles WHERE name = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, name)
	err = row.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No role found with name %v", name)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve role")
		return
	}
	log.Debug().Msgf("Retrieved role: %v", u)
	return
}
