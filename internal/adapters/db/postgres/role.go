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

type RoleRepo interface {
	Create(createRoleDTO dto.CreateRoleDTO) (roleID uuid.UUID, err error)
	GetOne(roleID uuid.UUID) (role entity.Role, err error)
	GetOneByName(name string) (user entity.Role, err error)
}

type roleRepo struct {
	pool *pgxpool.Pool
}

func NewRoleRepo(pool *pgxpool.Pool) *roleRepo {
	return &roleRepo{pool: pool}
}

func (r *roleRepo) Create(createRoleDTO dto.CreateRoleDTO) (roleID uuid.UUID, err error) {
	query := `INSERT INTO roles (name) VALUES ($1) RETURNING id`
	row := r.pool.QueryRow(context.Background(), query,
		createRoleDTO.Name,
	)
	var role entity.Role
	err = row.Scan(&role.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan role")
		return uuid.UUID{}, err
	}

	return role.ID, nil
}

func (r *roleRepo) GetOne(roleID uuid.UUID) (role entity.Role, err error) {
	query := `SELECT 
			id, 
			name,
			created_at, 
			updated_at, 
			deleted_at 
		FROM roles
		WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(context.Background(), query, roleID)
	err = row.Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
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

func (r *roleRepo) GetOneByName(name string) (role entity.Role, err error) {
	query := `SELECT 
			id, 
			name,
			created_at, 
			updated_at, 
			deleted_at 
		FROM roles
		WHERE name = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(context.Background(), query, name)
	err = row.Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
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
