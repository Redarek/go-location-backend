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

type wallTypeRepo struct {
	pool *pgxpool.Pool
}

func NewWallTypeRepo(pool *pgxpool.Pool) *wallTypeRepo {
	return &wallTypeRepo{pool: pool}
}

func (r *wallTypeRepo) Create(ctx context.Context, dto *dto.CreateWallTypeDTO) (wallTypeID uuid.UUID, err error) {
	query := `INSERT INTO wall_types (
		name, 
		color, 
		attenuation_24, attenuation_5, attenuation_6, 
		thickness, 
		site_id
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		dto.Name,
		dto.Color,
		dto.Attenuation24, dto.Attenuation5, dto.Attenuation6,
		dto.Thickness,
		dto.SiteID,
	)
	err = row.Scan(&wallTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan wallType")
		return
	}

	return
}

func (r *wallTypeRepo) GetOne(ctx context.Context, wallTypeID uuid.UUID) (wallType *entity.WallType, err error) {
	query := `SELECT 
		id, 
		name, 
		color, 
		attenuation_24, attenuation_5, attenuation_6, 
		thickness, 
		site_id,
		created_at, updated_at, deleted_at
	FROM wall_types 
	WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, wallTypeID)

	wallType = &entity.WallType{}
	err = row.Scan(
		&wallType.ID,
		&wallType.Name,
		&wallType.Color,
		&wallType.Attenuation24, &wallType.Attenuation5, &wallType.Attenuation6,
		&wallType.Thickness,
		&wallType.SiteID,
		&wallType.CreatedAt, &wallType.UpdatedAt, &wallType.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("no wallType found with ID %v", wallTypeID)
			return nil, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve wallType")
		return
	}
	log.Debug().Msgf("retrieved wallType: %v", wallType)
	return
}

// func (r *wallTypeRepo) GetOneDetailed(ctx context.Context, wallTypeID uuid.UUID) (wallType *entity.WallTypeDetailed, err error) {

// }

func (r *wallTypeRepo) GetAll(ctx context.Context, siteID uuid.UUID, limit, offset int) (wallTypes []*entity.WallType, err error) {
	query := `SELECT 
			id, 
			name, 
			color, 
			attenuation_24, attenuation_5, attenuation_6, 
			thickness, 
			site_id,
			created_at, updated_at, deleted_at
		FROM wall_types 
		WHERE site_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, siteID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve wallTypes")
		return
	}
	defer rows.Close()

	for rows.Next() {
		wallType := &entity.WallType{}
		err = rows.Scan(
			&wallType.ID,
			&wallType.Name,
			&wallType.Color,
			&wallType.Attenuation24, &wallType.Attenuation5, &wallType.Attenuation6,
			&wallType.Thickness,
			&wallType.SiteID,
			&wallType.CreatedAt, &wallType.UpdatedAt, &wallType.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan wallType")
			return
		}
		wallTypes = append(wallTypes, wallType)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(wallTypes)
	if length == 0 {
		log.Info().Msgf("wallTypes for site ID %v were not found", siteID)
		return nil, ErrNotFound
	}

	log.Debug().Msgf("retrieved %d wallTypes", length)
	return
}

func (r *wallTypeRepo) Update(ctx context.Context, updateWallTypeDTO *dto.PatchUpdateWallTypeDTO) (err error) {
	query := "UPDATE wall_types SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateWallTypeDTO.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, updateWallTypeDTO.Name)
		paramID++
	}
	if updateWallTypeDTO.Color != nil {
		updates = append(updates, fmt.Sprintf("color = $%d", paramID))
		params = append(params, updateWallTypeDTO.Color)
		paramID++
	}
	if updateWallTypeDTO.Attenuation24 != nil {
		updates = append(updates, fmt.Sprintf("attenuation_24 = $%d", paramID))
		params = append(params, updateWallTypeDTO.Attenuation24)
		paramID++
	}
	if updateWallTypeDTO.Attenuation5 != nil {
		updates = append(updates, fmt.Sprintf("attenuation_5 = $%d", paramID))
		params = append(params, updateWallTypeDTO.Attenuation5)
		paramID++
	}
	if updateWallTypeDTO.Attenuation6 != nil {
		updates = append(updates, fmt.Sprintf("attenuation_6 = $%d", paramID))
		params = append(params, updateWallTypeDTO.Attenuation6)
		paramID++
	}
	if updateWallTypeDTO.Thickness != nil {
		updates = append(updates, fmt.Sprintf("thickness = $%d", paramID))
		params = append(params, updateWallTypeDTO.Thickness)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateWallTypeDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no wallType found with the ID: %v", updateWallTypeDTO.ID)
		return ErrNotFound
	}

	return
}

// Checks if the wallType has been soft deleted
func (r *wallTypeRepo) IsWallTypeSoftDeleted(ctx context.Context, wallTypeID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM wall_types WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, wallTypeID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no wallType found with UUID %v", wallTypeID)
			return false, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve wallType")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is wallType deleted: %v", isDeleted)
	return
}

func (r *wallTypeRepo) SoftDelete(ctx context.Context, wallTypeID uuid.UUID) (err error) {
	query := `UPDATE wall_types SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, wallTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete wallType")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no wallType found with the UUID: %v", wallTypeID)
		return ErrNotFound
	}

	log.Debug().Msg("wallType deleted_at timestamp updated successfully")
	return
}

func (r *wallTypeRepo) Restore(ctx context.Context, wallTypeID uuid.UUID) (err error) {
	query := `UPDATE wall_types SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, wallTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore wallType")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no wallType found with the UUID: %v", wallTypeID)
		return ErrNotFound
	}

	log.Debug().Msg("wallType deleted_at timestamp set NULL successfully")
	return
}
