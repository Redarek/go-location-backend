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
	"location-backend/internal/domain/service"
)

type sensorTypeRepo struct {
	pool *pgxpool.Pool
}

func NewSensorTypeRepo(pool *pgxpool.Pool) *sensorTypeRepo {
	return &sensorTypeRepo{pool: pool}
}

func (r *sensorTypeRepo) Create(ctx context.Context, createSensorTypeDTO *dto.CreateSensorTypeDTO) (sensorTypeID uuid.UUID, err error) {
	query := `INSERT INTO sensor_types (
			name, 
			model,
			color, 
			z,
			is_virtual,
			site_id
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		createSensorTypeDTO.Name,
		createSensorTypeDTO.Model,
		createSensorTypeDTO.Color,
		createSensorTypeDTO.Z,
		createSensorTypeDTO.IsVirtual,
		createSensorTypeDTO.SiteID,
	)
	err = row.Scan(&sensorTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan sensor type")
		return
	}

	return
}

func (r *sensorTypeRepo) GetOne(ctx context.Context, sensorTypeID uuid.UUID) (sensorType *entity.SensorType, err error) {
	query := `SELECT 
			id, 
			name,
			model,
			color,
			z,
			is_virtual,
			site_id,
			created_at, updated_at, deleted_at
		FROM sensor_types WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, sensorTypeID)

	sensorType = &entity.SensorType{}
	err = row.Scan(
		&sensorType.ID,
		&sensorType.Name,
		&sensorType.Model,
		&sensorType.Color,
		&sensorType.Z,
		&sensorType.IsVirtual,
		&sensorType.SiteID,
		&sensorType.CreatedAt, &sensorType.UpdatedAt, &sensorType.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no sensorType found with ID %v", sensorTypeID)
			return nil, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve sensor type")
		return
	}
	log.Debug().Msgf("retrieved sensor type: %v", sensorType)
	return
}

// func (r *sensorTypeRepo) GetOneDetailed(ctx context.Context, sensorTypeID uuid.UUID) (sensorType *entity.SensorTypeDetailed, err error) {

// }

func (r *sensorTypeRepo) GetAll(ctx context.Context, siteID uuid.UUID, limit, offset int) (sensorTypes []*entity.SensorType, err error) {
	query := `SELECT 
			id, 
			name,
			model,
			color,
			z,
			is_virtual,
			site_id,
			created_at, updated_at, deleted_at
		FROM sensor_types 
		WHERE site_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, siteID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve sensor types")
		return
	}
	defer rows.Close()

	for rows.Next() {
		sensorType := &entity.SensorType{}
		err = rows.Scan(
			&sensorType.ID,
			&sensorType.Name,
			&sensorType.Model,
			&sensorType.Color,
			&sensorType.Z,
			&sensorType.IsVirtual,
			&sensorType.SiteID,
			&sensorType.CreatedAt, &sensorType.UpdatedAt, &sensorType.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan sensor type")
			return
		}
		sensorTypes = append(sensorTypes, sensorType)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(sensorTypes)
	if length == 0 {
		log.Info().Msgf("sensor types for site ID %v were not found", siteID)
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved %d sensor types", length)
	return
}

func (r *sensorTypeRepo) Update(ctx context.Context, updateSensorTypeDTO *dto.PatchUpdateSensorTypeDTO) (err error) {
	query := "UPDATE sensor_types SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateSensorTypeDTO.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, updateSensorTypeDTO.Name)
		paramID++
	}
	if updateSensorTypeDTO.Model != nil {
		updates = append(updates, fmt.Sprintf("model = $%d", paramID))
		params = append(params, updateSensorTypeDTO.Model)
		paramID++
	}

	if updateSensorTypeDTO.Color != nil {
		updates = append(updates, fmt.Sprintf("color = $%d", paramID))
		params = append(params, updateSensorTypeDTO.Color)
		paramID++
	}
	if updateSensorTypeDTO.Z != nil {
		updates = append(updates, fmt.Sprintf("z = $%d", paramID))
		params = append(params, updateSensorTypeDTO.Z)
		paramID++
	}
	if updateSensorTypeDTO.IsVirtual != nil {
		updates = append(updates, fmt.Sprintf("is_virtual = $%d", paramID))
		params = append(params, updateSensorTypeDTO.IsVirtual)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return service.ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateSensorTypeDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no sensor type found with the UUID: %v", updateSensorTypeDTO.ID)
		return service.ErrNotFound
	}

	return
}

// Checks if the sensorType has been soft deleted
func (r *sensorTypeRepo) IsSensorTypeSoftDeleted(ctx context.Context, sensorTypeID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM sensor_types WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, sensorTypeID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no sensor type found with UUID %v", sensorTypeID)
			return false, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve sensor type")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is sensor type deleted: %v", isDeleted)
	return
}

func (r *sensorTypeRepo) SoftDelete(ctx context.Context, sensorTypeID uuid.UUID) (err error) {
	query := `UPDATE sensor_types SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, sensorTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete sensor type")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no sensor type found with the ID: %v", sensorTypeID)
		return service.ErrNotFound
	}

	log.Debug().Msg("sensor type deleted_at timestamp updated successfully")
	return
}

func (r *sensorTypeRepo) Restore(ctx context.Context, sensorTypeID uuid.UUID) (err error) {
	query := `UPDATE sensor_types SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, sensorTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore sensor type")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no sensor type found with the UUID: %v", sensorTypeID)
		return service.ErrNotFound
	}

	log.Debug().Msg("sensor type deleted_at timestamp set NULL successfully")
	return
}
