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

type sensorRadioRepo struct {
	pool *pgxpool.Pool
}

func NewSensorRadioRepo(pool *pgxpool.Pool) *sensorRadioRepo {
	return &sensorRadioRepo{pool: pool}
}

func (r *sensorRadioRepo) Create(ctx context.Context, createSensorRadioDTO *dto.CreateSensorRadioDTO) (sensorRadioID uuid.UUID, err error) {
	query := `INSERT INTO sensor_radios (
			number, 
			channel, channel2,
			channel_width,
			wifi,
			power,
			bandwidth,
			guard_interval,
			is_active,
			sensor_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		createSensorRadioDTO.Number,
		createSensorRadioDTO.Channel, createSensorRadioDTO.Channel2,
		createSensorRadioDTO.ChannelWidth,
		createSensorRadioDTO.WiFi,
		createSensorRadioDTO.Power,
		createSensorRadioDTO.Bandwidth,
		createSensorRadioDTO.GuardInterval,
		createSensorRadioDTO.IsActive,
		createSensorRadioDTO.SensorID,
	)
	err = row.Scan(&sensorRadioID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan sensorRadio")
		return
	}

	return
}

func (r *sensorRadioRepo) GetOne(ctx context.Context, sensorRadioID uuid.UUID) (sensorRadio *entity.SensorRadio, err error) {
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
			sensor_id,
			created_at, updated_at, deleted_at
		FROM sensor_radios WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, sensorRadioID)

	sensorRadio = &entity.SensorRadio{}
	err = row.Scan(
		&sensorRadio.ID,
		&sensorRadio.Number,
		&sensorRadio.Channel, &sensorRadio.Channel2,
		&sensorRadio.ChannelWidth,
		&sensorRadio.WiFi,
		&sensorRadio.Power,
		&sensorRadio.Bandwidth,
		&sensorRadio.GuardInterval,
		&sensorRadio.IsActive,
		&sensorRadio.SensorID,
		&sensorRadio.CreatedAt, &sensorRadio.UpdatedAt, &sensorRadio.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msgf("no sensor radio found with ID %v", sensorRadioID)
			return nil, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve sensor radio")
		return
	}
	log.Debug().Msgf("retrieved sensor radio: %v", sensorRadio)
	return
}

// func (r *sensorRadioRepo) GetOneDetailed(ctx context.Context, sensorRadioID uuid.UUID) (sensorRadio *entity.SensorRadioDetailed, err error) {

// }

func (r *sensorRadioRepo) GetAll(ctx context.Context, sensorID uuid.UUID, limit, offset int) (sensorRadios []*entity.SensorRadio, err error) {
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
			sensor_id,
			created_at, updated_at, deleted_at
		FROM sensor_radios 
		WHERE sensor_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, sensorID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve sensor radios")
		return
	}
	defer rows.Close()

	for rows.Next() {
		sensorRadio := &entity.SensorRadio{}
		err = rows.Scan(
			&sensorRadio.ID,
			&sensorRadio.Number,
			&sensorRadio.Channel, &sensorRadio.Channel2,
			&sensorRadio.ChannelWidth,
			&sensorRadio.WiFi,
			&sensorRadio.Power,
			&sensorRadio.Bandwidth,
			&sensorRadio.GuardInterval,
			&sensorRadio.IsActive,
			&sensorRadio.SensorID,
			&sensorRadio.CreatedAt, &sensorRadio.UpdatedAt, &sensorRadio.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan sensor radio")
			return
		}
		sensorRadios = append(sensorRadios, sensorRadio)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(sensorRadios)
	if length == 0 {
		log.Info().Msgf("sensor radios for sensor ID %v were not found", sensorID)
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved %d sensor radios", length)
	return
}

func (r *sensorRadioRepo) Update(ctx context.Context, updateSensorRadioDTO *dto.PatchUpdateSensorRadioDTO) (err error) {
	query := "UPDATE sensor_radios SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateSensorRadioDTO.Number != nil {
		updates = append(updates, fmt.Sprintf("number = $%d", paramID))
		params = append(params, updateSensorRadioDTO.Number)
		paramID++
	}
	if updateSensorRadioDTO.Channel != nil {
		updates = append(updates, fmt.Sprintf("channel = $%d", paramID))
		params = append(params, updateSensorRadioDTO.Channel)
		paramID++
	}
	if updateSensorRadioDTO.Channel2 != nil {
		updates = append(updates, fmt.Sprintf("channel2 = $%d", paramID))
		params = append(params, updateSensorRadioDTO.Channel2)
		paramID++
	}
	if updateSensorRadioDTO.ChannelWidth != nil {
		updates = append(updates, fmt.Sprintf("channel_width = $%d", paramID))
		params = append(params, updateSensorRadioDTO.ChannelWidth)
		paramID++
	}
	if updateSensorRadioDTO.WiFi != nil {
		updates = append(updates, fmt.Sprintf("wifi = $%d", paramID))
		params = append(params, updateSensorRadioDTO.WiFi)
		paramID++
	}
	if updateSensorRadioDTO.Power != nil {
		updates = append(updates, fmt.Sprintf("power = $%d", paramID))
		params = append(params, updateSensorRadioDTO.Power)
		paramID++
	}
	if updateSensorRadioDTO.Bandwidth != nil {
		updates = append(updates, fmt.Sprintf("bandwidth = $%d", paramID))
		params = append(params, updateSensorRadioDTO.Bandwidth)
		paramID++
	}
	if updateSensorRadioDTO.GuardInterval != nil {
		updates = append(updates, fmt.Sprintf("guard_interval = $%d", paramID))
		params = append(params, updateSensorRadioDTO.GuardInterval)
		paramID++
	}
	if updateSensorRadioDTO.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", paramID))
		params = append(params, updateSensorRadioDTO.IsActive)
		paramID++
	}
	if updateSensorRadioDTO.SensorID != nil {
		updates = append(updates, fmt.Sprintf("sensor_id = $%d", paramID))
		params = append(params, updateSensorRadioDTO.SensorID)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return service.ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateSensorRadioDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no sensor radio found with the ID: %v", updateSensorRadioDTO.ID)
		return service.ErrNotFound
	}

	return
}

// Checks if the sensorRadio has been soft deleted
func (r *sensorRadioRepo) IsSensorRadioSoftDeleted(ctx context.Context, sensorRadioID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM sensor_radios WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, sensorRadioID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msgf("no sensor radio found with ID %v", sensorRadioID)
			return false, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve sensor radio")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is sensor radio deleted: %v", isDeleted)
	return
}

func (r *sensorRadioRepo) SoftDelete(ctx context.Context, sensorRadioID uuid.UUID) (err error) {
	query := `UPDATE sensor_radios SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, sensorRadioID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete sensor radio")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no sensor radio found with the ID: %v", sensorRadioID)
		return service.ErrNotFound
	}

	log.Debug().Msg("sensor radio deleted_at timestamp updated successfully")
	return
}

func (r *sensorRadioRepo) Restore(ctx context.Context, sensorRadioID uuid.UUID) (err error) {
	query := `UPDATE sensor_radios SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, sensorRadioID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore sensor radio")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no sensor radio found with the ID: %v", sensorRadioID)
		return service.ErrNotFound
	}

	log.Debug().Msg("sensor radio deleted_at timestamp set NULL successfully")
	return
}
