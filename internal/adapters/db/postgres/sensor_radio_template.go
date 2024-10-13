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

type sensorRadioTemplateRepo struct {
	pool *pgxpool.Pool
}

func NewSensorRadioTemplateRepo(pool *pgxpool.Pool) *sensorRadioTemplateRepo {
	return &sensorRadioTemplateRepo{pool: pool}
}

func (r *sensorRadioTemplateRepo) Create(ctx context.Context, createSensorRadioTemplateDTO *dto.CreateSensorRadioTemplateDTO) (sensorRadioTemplateID uuid.UUID, err error) {
	query := `INSERT INTO sensor_radio_templates (
			number, 
			channel, channel2,
			channel_width,
			wifi,
			power,
			bandwidth,
			guard_interval,
			is_active,
			sensor_type_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		createSensorRadioTemplateDTO.Number,
		createSensorRadioTemplateDTO.Channel, createSensorRadioTemplateDTO.Channel2,
		createSensorRadioTemplateDTO.ChannelWidth,
		createSensorRadioTemplateDTO.WiFi,
		createSensorRadioTemplateDTO.Power,
		createSensorRadioTemplateDTO.Bandwidth,
		createSensorRadioTemplateDTO.GuardInterval,
		createSensorRadioTemplateDTO.IsActive,
		createSensorRadioTemplateDTO.SensorTypeID,
	)
	err = row.Scan(&sensorRadioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan sensor radio template")
		return
	}

	return
}

func (r *sensorRadioTemplateRepo) GetOne(ctx context.Context, sensorRadioTemplateID uuid.UUID) (sensorRadioTemplate *entity.SensorRadioTemplate, err error) {
	query := `SELECT 
			id, 
			number, 
			channel, channel2,
			channel_width,
			wifi,
			power,
			bandwidth,
			guard_interval,
			is_active,
			sensor_type_id,
			created_at, updated_at, deleted_at
		FROM sensor_radio_templates WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, sensorRadioTemplateID)

	sensorRadioTemplate = &entity.SensorRadioTemplate{}
	err = row.Scan(
		&sensorRadioTemplate.ID,
		&sensorRadioTemplate.Number,
		&sensorRadioTemplate.Channel, &sensorRadioTemplate.Channel2,
		&sensorRadioTemplate.ChannelWidth,
		&sensorRadioTemplate.WiFi,
		&sensorRadioTemplate.Power,
		&sensorRadioTemplate.Bandwidth,
		&sensorRadioTemplate.GuardInterval,
		&sensorRadioTemplate.IsActive,
		&sensorRadioTemplate.SensorTypeID,
		&sensorRadioTemplate.CreatedAt, &sensorRadioTemplate.UpdatedAt, &sensorRadioTemplate.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msgf("no sensor radio template found with ID %v", sensorRadioTemplateID)
			return nil, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve sensor radio template")
		return
	}
	log.Debug().Msgf("retrieved sensor radio template: %v", sensorRadioTemplate)
	return
}

// func (r *sensorRadioTemplateRepo) GetOneDetailed(ctx context.Context, sensorRadioTemplateID uuid.UUID) (sensorRadioTemplate *entity.SensorRadioTemplateDetailed, err error) {

// }

func (r *sensorRadioTemplateRepo) GetAll(ctx context.Context, sensorTypeID uuid.UUID, limit, offset int) (sensorRadioTemplates []*entity.SensorRadioTemplate, err error) {
	query := `SELECT 
			id, 
			number, 
			channel, channel2,
			channel_width,
			wifi,
			power,
			bandwidth,
			guard_interval,
			is_active,
			sensor_type_id,
			created_at, updated_at, deleted_at
		FROM sensor_radio_templates 
		WHERE sensor_type_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, sensorTypeID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve sensor radio template")
		return
	}
	defer rows.Close()

	for rows.Next() {
		sensorRadioTemplate := &entity.SensorRadioTemplate{}
		err = rows.Scan(
			&sensorRadioTemplate.ID,
			&sensorRadioTemplate.Number,
			&sensorRadioTemplate.Channel, &sensorRadioTemplate.Channel2,
			&sensorRadioTemplate.ChannelWidth,
			&sensorRadioTemplate.WiFi,
			&sensorRadioTemplate.Power,
			&sensorRadioTemplate.Bandwidth,
			&sensorRadioTemplate.GuardInterval,
			&sensorRadioTemplate.IsActive,
			&sensorRadioTemplate.SensorTypeID,
			&sensorRadioTemplate.CreatedAt, &sensorRadioTemplate.UpdatedAt, &sensorRadioTemplate.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan sensor radio template")
			return
		}
		sensorRadioTemplates = append(sensorRadioTemplates, sensorRadioTemplate)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(sensorRadioTemplates)
	if length == 0 {
		log.Info().Msgf("sensor radio templates for sensor type ID %v were not found", sensorTypeID)
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved %d sensor radio templates", length)
	return
}

func (r *sensorRadioTemplateRepo) Update(ctx context.Context, updateSensorRadioTemplateDTO *dto.PatchUpdateSensorRadioTemplateDTO) (err error) {
	query := "UPDATE sensor_radio_templates SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateSensorRadioTemplateDTO.Number != nil {
		updates = append(updates, fmt.Sprintf("number = $%d", paramID))
		params = append(params, updateSensorRadioTemplateDTO.Number)
		paramID++
	}
	if updateSensorRadioTemplateDTO.Channel != nil {
		updates = append(updates, fmt.Sprintf("channel = $%d", paramID))
		params = append(params, updateSensorRadioTemplateDTO.Channel)
		paramID++
	}
	if updateSensorRadioTemplateDTO.Channel2 != nil {
		updates = append(updates, fmt.Sprintf("channel2 = $%d", paramID))
		params = append(params, updateSensorRadioTemplateDTO.Channel2)
		paramID++
	}
	if updateSensorRadioTemplateDTO.ChannelWidth != nil {
		updates = append(updates, fmt.Sprintf("channel_width = $%d", paramID))
		params = append(params, updateSensorRadioTemplateDTO.ChannelWidth)
		paramID++
	}
	if updateSensorRadioTemplateDTO.WiFi != nil {
		updates = append(updates, fmt.Sprintf("wifi = $%d", paramID))
		params = append(params, updateSensorRadioTemplateDTO.WiFi)
		paramID++
	}
	if updateSensorRadioTemplateDTO.Power != nil {
		updates = append(updates, fmt.Sprintf("power = $%d", paramID))
		params = append(params, updateSensorRadioTemplateDTO.Power)
		paramID++
	}
	if updateSensorRadioTemplateDTO.Bandwidth != nil {
		updates = append(updates, fmt.Sprintf("bandwidth = $%d", paramID))
		params = append(params, updateSensorRadioTemplateDTO.Bandwidth)
		paramID++
	}
	if updateSensorRadioTemplateDTO.GuardInterval != nil {
		updates = append(updates, fmt.Sprintf("guard_interval = $%d", paramID))
		params = append(params, updateSensorRadioTemplateDTO.GuardInterval)
		paramID++
	}
	if updateSensorRadioTemplateDTO.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", paramID))
		params = append(params, updateSensorRadioTemplateDTO.IsActive)
		paramID++
	}
	if updateSensorRadioTemplateDTO.SensorTypeID != nil {
		updates = append(updates, fmt.Sprintf("sensor_type_id = $%d", paramID))
		params = append(params, updateSensorRadioTemplateDTO.SensorTypeID)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return service.ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateSensorRadioTemplateDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no sensor radio template found with the ID: %v", updateSensorRadioTemplateDTO.ID)
		return service.ErrNotFound
	}

	return
}

// Checks if the sensorRadioTemplate has been soft deleted
func (r *sensorRadioTemplateRepo) IsSensorRadioTemplateSoftDeleted(ctx context.Context, sensorRadioTemplateID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM sensor_radio_templates WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, sensorRadioTemplateID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msgf("no sensor radio template found with ID %v", sensorRadioTemplateID)
			return false, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve sensor radio template")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is sensor radio template deleted: %v", isDeleted)
	return
}

func (r *sensorRadioTemplateRepo) SoftDelete(ctx context.Context, sensorRadioTemplateID uuid.UUID) (err error) {
	query := `UPDATE sensor_radio_templates SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, sensorRadioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete sensor radio template")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no sensor radio template found with the ID: %v", sensorRadioTemplateID)
		return service.ErrNotFound
	}

	log.Debug().Msg("sensor radio template deleted_at timestamp updated successfully")
	return
}

func (r *sensorRadioTemplateRepo) Restore(ctx context.Context, sensorRadioTemplateID uuid.UUID) (err error) {
	query := `UPDATE sensor_radio_templates SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, sensorRadioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore sensor radio template")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no sensor radio template found with the ID: %v", sensorRadioTemplateID)
		return service.ErrNotFound
	}

	log.Debug().Msg("sensor radio template deleted_at timestamp set NULL successfully")
	return
}
