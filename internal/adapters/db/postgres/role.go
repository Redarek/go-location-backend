package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type roleRepo struct {
	pool *pgxpool.Pool
}

func NewRoleRepo(pool *pgxpool.Pool) *roleRepo {
	return &roleRepo{pool: pool}
}

func (r *roleRepo) Create(ctx context.Context, createRoleDTO *dto.CreateRoleDTO) (roleID uuid.UUID, err error) {
	query := `INSERT INTO roles (name) VALUES ($1) RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		createRoleDTO.Name,
	)

	err = row.Scan(&roleID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan role")
		return uuid.UUID{}, err
	}

	return
}

func (r *roleRepo) GetOne(ctx context.Context, roleID uuid.UUID) (role *entity.Role, err error) {
	query := `SELECT 
			id, 
			name,
			created_at, 
			updated_at, 
			deleted_at 
		FROM roles
		WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, roleID)

	role = &entity.Role{}
	err = row.Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt, &role.UpdatedAt, &role.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msgf("role with ID %v not found", roleID)
			return role, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to scan role")
		return
	}
	log.Debug().Msgf("retrieved role: %v", role)
	return
}

func (r *roleRepo) GetOneByName(ctx context.Context, name string) (role *entity.Role, err error) {
	query := `SELECT 
			id, 
			name,
			created_at, 
			updated_at, 
			deleted_at 
		FROM roles
		WHERE name = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, name)

	role = &entity.Role{}
	err = row.Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt, &role.UpdatedAt, &role.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msgf("role %v not found", name)
			return role, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to scan role")
		return
	}
	log.Debug().Msgf("retrieved role: %v", role)
	return
}
