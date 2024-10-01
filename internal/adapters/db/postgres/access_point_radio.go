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

type accessPointRadioRepo struct {
	pool *pgxpool.Pool
}

func NewAccessPointRadioRepo(pool *pgxpool.Pool) *accessPointRadioRepo {
	return &accessPointRadioRepo{pool: pool}
}

func (r *accessPointRadioRepo) Create(ctx context.Context, createAccessPointRadioDTO *dto.CreateAccessPointRadioDTO) (accessPointRadioID uuid.UUID, err error) {
	query := `INSERT INTO access_point_radios (
			number, 
			channel,
			channel2,
			channel_width,
			wifi,
			power,
			bandwidth,
			guard_interval,
			is_active,
			access_point_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		createAccessPointRadioDTO.Number,
		createAccessPointRadioDTO.Channel,
		createAccessPointRadioDTO.Channel2,
		createAccessPointRadioDTO.ChannelWidth,
		createAccessPointRadioDTO.WiFi,
		createAccessPointRadioDTO.Power,
		createAccessPointRadioDTO.Bandwidth,
		createAccessPointRadioDTO.GuardInterval,
		createAccessPointRadioDTO.IsActive,
		createAccessPointRadioDTO.AccessPointID,
	)
	err = row.Scan(&accessPointRadioID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan accessPointRadio")
		return
	}

	return
}

func (r *accessPointRadioRepo) GetOne(ctx context.Context, accessPointRadioID uuid.UUID) (accessPointRadio *entity.AccessPointRadio, err error) {
	query := `SELECT 
			id, 
			number, 
			channel,
			channel2,
			channel_width,
			wifi,
			power,
			bandwidth,
			guard_interval,
			is_active,
			access_point_id,
			created_at, updated_at, deleted_at
		FROM access_point_radios WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, accessPointRadioID)

	accessPointRadio = &entity.AccessPointRadio{}
	err = row.Scan(
		&accessPointRadio.ID,
		&accessPointRadio.Number,
		&accessPointRadio.Channel,
		&accessPointRadio.Channel2,
		&accessPointRadio.ChannelWidth,
		&accessPointRadio.WiFi,
		&accessPointRadio.Power,
		&accessPointRadio.Bandwidth,
		&accessPointRadio.GuardInterval,
		&accessPointRadio.IsActive,
		&accessPointRadio.AccessPointID,
		&accessPointRadio.CreatedAt, &accessPointRadio.UpdatedAt, &accessPointRadio.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msgf("no access point radio found with ID %v", accessPointRadioID)
			return nil, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve access point radio")
		return
	}
	log.Debug().Msgf("retrieved access point radio: %v", accessPointRadio)
	return
}

// func (r *accessPointRadioRepo) GetOneDetailed(ctx context.Context, accessPointRadioID uuid.UUID) (accessPointRadio *entity.AccessPointRadioDetailed, err error) {

// }

func (r *accessPointRadioRepo) GetAll(ctx context.Context, accessPointID uuid.UUID, limit, offset int) (accessPointRadios []*entity.AccessPointRadio, err error) {
	query := `SELECT 
			id, 
			number, 
			channel,
			channel2,
			channel_width,
			wifi,
			power,
			bandwidth,
			guard_interval,
			is_active,
			access_point_id,
			created_at, updated_at, deleted_at
		FROM access_point_radios 
		WHERE access_point_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, accessPointID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve access point radios")
		return
	}
	defer rows.Close()

	for rows.Next() {
		accessPointRadio := &entity.AccessPointRadio{}
		err = rows.Scan(
			&accessPointRadio.ID,
			&accessPointRadio.Number,
			&accessPointRadio.Channel,
			&accessPointRadio.Channel2,
			&accessPointRadio.ChannelWidth,
			&accessPointRadio.WiFi,
			&accessPointRadio.Power,
			&accessPointRadio.Bandwidth,
			&accessPointRadio.GuardInterval,
			&accessPointRadio.IsActive,
			&accessPointRadio.AccessPointID,
			&accessPointRadio.CreatedAt, &accessPointRadio.UpdatedAt, &accessPointRadio.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan access point radio")
			return
		}
		accessPointRadios = append(accessPointRadios, accessPointRadio)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(accessPointRadios)
	if length == 0 {
		log.Info().Msgf("access point radios for access point ID %v were not found", accessPointID)
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved %d access point radios", length)
	return
}

func (r *accessPointRadioRepo) Update(ctx context.Context, updateAccessPointRadioDTO *dto.PatchUpdateAccessPointRadioDTO) (err error) {
	query := "UPDATE access_point_radios SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateAccessPointRadioDTO.Number != nil {
		updates = append(updates, fmt.Sprintf("number = $%d", paramID))
		params = append(params, updateAccessPointRadioDTO.Number)
		paramID++
	}
	if updateAccessPointRadioDTO.Channel != nil {
		updates = append(updates, fmt.Sprintf("channel = $%d", paramID))
		params = append(params, updateAccessPointRadioDTO.Channel)
		paramID++
	}
	if updateAccessPointRadioDTO.Channel2 != nil {
		updates = append(updates, fmt.Sprintf("channel2 = $%d", paramID))
		params = append(params, updateAccessPointRadioDTO.Channel2)
		paramID++
	}
	if updateAccessPointRadioDTO.ChannelWidth != nil {
		updates = append(updates, fmt.Sprintf("channel_width = $%d", paramID))
		params = append(params, updateAccessPointRadioDTO.ChannelWidth)
		paramID++
	}
	if updateAccessPointRadioDTO.WiFi != nil {
		updates = append(updates, fmt.Sprintf("wifi = $%d", paramID))
		params = append(params, updateAccessPointRadioDTO.WiFi)
		paramID++
	}
	if updateAccessPointRadioDTO.Power != nil {
		updates = append(updates, fmt.Sprintf("power = $%d", paramID))
		params = append(params, updateAccessPointRadioDTO.Power)
		paramID++
	}
	if updateAccessPointRadioDTO.Bandwidth != nil {
		updates = append(updates, fmt.Sprintf("bandwidth = $%d", paramID))
		params = append(params, updateAccessPointRadioDTO.Bandwidth)
		paramID++
	}
	if updateAccessPointRadioDTO.GuardInterval != nil {
		updates = append(updates, fmt.Sprintf("guard_interval = $%d", paramID))
		params = append(params, updateAccessPointRadioDTO.GuardInterval)
		paramID++
	}
	if updateAccessPointRadioDTO.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", paramID))
		params = append(params, updateAccessPointRadioDTO.IsActive)
		paramID++
	}
	if updateAccessPointRadioDTO.AccessPointID != nil {
		updates = append(updates, fmt.Sprintf("access_point_id = $%d", paramID))
		params = append(params, updateAccessPointRadioDTO.AccessPointID)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return service.ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateAccessPointRadioDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no access point radio found with the ID: %v", updateAccessPointRadioDTO.ID)
		return service.ErrNotFound
	}

	return
}

// Checks if the accessPointRadio has been soft deleted
func (r *accessPointRadioRepo) IsAccessPointRadioSoftDeleted(ctx context.Context, accessPointRadioID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM access_point_radios WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, accessPointRadioID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msgf("no access point radio found with ID %v", accessPointRadioID)
			return false, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve access point radio")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is access point radio deleted: %v", isDeleted)
	return
}

func (r *accessPointRadioRepo) SoftDelete(ctx context.Context, accessPointRadioID uuid.UUID) (err error) {
	query := `UPDATE access_point_radios SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, accessPointRadioID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete access point radio")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no access point radio found with the ID: %v", accessPointRadioID)
		return service.ErrNotFound
	}

	log.Debug().Msg("access point radio deleted_at timestamp updated successfully")
	return
}

func (r *accessPointRadioRepo) Restore(ctx context.Context, accessPointRadioID uuid.UUID) (err error) {
	query := `UPDATE access_point_radios SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, accessPointRadioID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore access point radio")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no access point radio found with the ID: %v", accessPointRadioID)
		return service.ErrNotFound
	}

	log.Debug().Msg("access point radio deleted_at timestamp set NULL successfully")
	return
}
