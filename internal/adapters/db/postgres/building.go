package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type BuildingRepo interface {
	Create(ctx context.Context, createBuildingDTO *dto.CreateBuildingDTO) (buildingID uuid.UUID, err error)
	GetOne(ctx context.Context, buildingID uuid.UUID) (building *entity.Building, err error)
	GetAll(ctx context.Context, siteID uuid.UUID, limit, offset int) (buildings []*entity.Building, err error)

	Update(ctx context.Context, updateBuildingDTO *dto.PatchUpdateBuildingDTO) (err error)

	IsBuildingSoftDeleted(ctx context.Context, buildingID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, buildingID uuid.UUID) (err error)
	Restore(ctx context.Context, buildingID uuid.UUID) (err error)
}

type buildingRepo struct {
	pool *pgxpool.Pool
}

func NewBuildingRepo(pool *pgxpool.Pool) *buildingRepo {
	return &buildingRepo{pool: pool}
}

func (r *buildingRepo) Create(ctx context.Context, createBuildingDTO *dto.CreateBuildingDTO) (buildingID uuid.UUID, err error) {
	query := `INSERT INTO buildings (
			name, 
			description, 
			country, 
			city, 
			address, 
			site_id
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		createBuildingDTO.Name,
		createBuildingDTO.Description,
		createBuildingDTO.Country,
		createBuildingDTO.City,
		createBuildingDTO.Address,
		createBuildingDTO.SiteID,
	)
	err = row.Scan(&buildingID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan building")
		return
	}

	return
}

func (r *buildingRepo) GetOne(ctx context.Context, buildingID uuid.UUID) (building *entity.Building, err error) {
	query := `SELECT 
			id, 
			name, 
			description, 
			country,
			city,
			address,
			site_id,
			created_at, updated_at, deleted_at
		FROM buildings WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, buildingID)

	building = &entity.Building{}
	err = row.Scan(
		&building.ID,
		&building.Name,
		&building.Description,
		&building.Country,
		&building.City,
		&building.Address,
		&building.SiteID,
		&building.CreatedAt, &building.UpdatedAt, &building.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("no building found with ID %v", buildingID)
			return nil, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve building")
		return
	}
	log.Debug().Msgf("retrieved building: %v", building)
	return
}

func (r *buildingRepo) GetAll(ctx context.Context, siteID uuid.UUID, limit, offset int) (buildings []*entity.Building, err error) {
	query := `SELECT
			id, 
			name, 
			description, 
			country,
			city,
			address,
			site_id,
			created_at, updated_at, deleted_at
		FROM buildings WHERE site_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, siteID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve buildings")
		return
	}
	defer rows.Close()

	for rows.Next() {
		building := &entity.Building{}
		err = rows.Scan(
			&building.ID,
			&building.Name,
			&building.Description,
			&building.Country,
			&building.City,
			&building.Address,
			&building.SiteID,
			&building.CreatedAt, &building.UpdatedAt, &building.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan building")
			return
		}
		buildings = append(buildings, building)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(buildings)
	if length == 0 {
		log.Info().Msgf("buildings for site ID %v were not found", siteID)
		return nil, ErrNotFound
	}

	log.Debug().Msgf("retrieved %d buildings", length)
	return
}

func (r *buildingRepo) Update(ctx context.Context, updateBuildingDTO *dto.PatchUpdateBuildingDTO) (err error) {
	query := "UPDATE buildings SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateBuildingDTO.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, updateBuildingDTO.Name)
		paramID++
	}
	if updateBuildingDTO.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", paramID))
		params = append(params, updateBuildingDTO.Description)
		paramID++
	}
	if updateBuildingDTO.Country != nil {
		updates = append(updates, fmt.Sprintf("country = $%d", paramID))
		params = append(params, updateBuildingDTO.Country)
		paramID++
	}
	if updateBuildingDTO.City != nil {
		updates = append(updates, fmt.Sprintf("city = $%d", paramID))
		params = append(params, updateBuildingDTO.City)
		paramID++
	}
	if updateBuildingDTO.Address != nil {
		updates = append(updates, fmt.Sprintf("address = $%d", paramID))
		params = append(params, updateBuildingDTO.Address)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateBuildingDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no building found with the UUID: %v", updateBuildingDTO.ID)
		return ErrNotFound
	}

	return
}

// Checks if the building has been soft deleted
func (r *buildingRepo) IsBuildingSoftDeleted(ctx context.Context, buildingID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM buildings WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, buildingID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no building found with UUID %v", buildingID)
			return false, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve building")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is building deleted: %v", isDeleted)
	return
}

func (r *buildingRepo) SoftDelete(ctx context.Context, buildingID uuid.UUID) (err error) {
	query := `UPDATE buildings SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, buildingID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete building")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no building found with the UUID: %v", buildingID)
		return ErrNotFound
	}

	log.Debug().Msg("building deleted_at timestamp updated successfully")
	return
}

func (r *buildingRepo) Restore(ctx context.Context, buildingID uuid.UUID) (err error) {
	query := `UPDATE buildings SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, buildingID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore building")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no building found with the UUID: %v", buildingID)
		return ErrNotFound
	}

	log.Debug().Msg("building deleted_at timestamp set NULL successfully")
	return
}
