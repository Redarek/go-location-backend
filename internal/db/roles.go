package db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type Role struct {
	ID        uuid.UUID          `json:"id" db:"id"`
	Name      string             `json:"name" db:"name"`
	CreatedAt pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
}

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
