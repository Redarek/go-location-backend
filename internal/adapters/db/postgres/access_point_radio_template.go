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

type accessPointRadioTemplateRepo struct {
	pool *pgxpool.Pool
}

func NewAccessPointRadioTemplateRepo(pool *pgxpool.Pool) *accessPointRadioTemplateRepo {
	return &accessPointRadioTemplateRepo{pool: pool}
}

func (r *accessPointRadioTemplateRepo) Create(ctx context.Context, createAccessPointRadioTemplateDTO *dto.CreateAccessPointRadioTemplateDTO) (accessPointRadioTemplateID uuid.UUID, err error) {
	query := `INSERT INTO access_point_radio_templates (
			number, 
			channel,
			channel2,
			channel_width,
			wifi,
			power,
			bandwidth,
			guard_interval,
			is_active,
			access_point_type_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		createAccessPointRadioTemplateDTO.Number,
		createAccessPointRadioTemplateDTO.Channel,
		createAccessPointRadioTemplateDTO.Channel2,
		createAccessPointRadioTemplateDTO.ChannelWidth,
		createAccessPointRadioTemplateDTO.WiFi,
		createAccessPointRadioTemplateDTO.Power,
		createAccessPointRadioTemplateDTO.Bandwidth,
		createAccessPointRadioTemplateDTO.GuardInterval,
		createAccessPointRadioTemplateDTO.IsActive,
		createAccessPointRadioTemplateDTO.AccessPointTypeID,
	)
	err = row.Scan(&accessPointRadioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan accessPointRadioTemplate")
		return
	}

	return
}

func (r *accessPointRadioTemplateRepo) GetOne(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (accessPointRadioTemplate *entity.AccessPointRadioTemplate, err error) {
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
			access_point_type_id,
			created_at, updated_at, deleted_at
		FROM access_point_radio_templates WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, accessPointRadioTemplateID)

	accessPointRadioTemplate = &entity.AccessPointRadioTemplate{}
	err = row.Scan(
		&accessPointRadioTemplate.ID,
		&accessPointRadioTemplate.Number,
		&accessPointRadioTemplate.Channel,
		&accessPointRadioTemplate.Channel2,
		&accessPointRadioTemplate.ChannelWidth,
		&accessPointRadioTemplate.WiFi,
		&accessPointRadioTemplate.Power,
		&accessPointRadioTemplate.Bandwidth,
		&accessPointRadioTemplate.GuardInterval,
		&accessPointRadioTemplate.IsActive,
		&accessPointRadioTemplate.AccessPointTypeID,
		&accessPointRadioTemplate.CreatedAt, &accessPointRadioTemplate.UpdatedAt, &accessPointRadioTemplate.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msgf("no access point radio template found with ID %v", accessPointRadioTemplateID)
			return nil, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve access point radio template")
		return
	}
	log.Debug().Msgf("retrieved access point radio template: %v", accessPointRadioTemplate)
	return
}

// func (r *accessPointRadioTemplateRepo) GetOneDetailed(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (accessPointRadioTemplate *entity.AccessPointRadioTemplateDetailed, err error) {

// }

func (r *accessPointRadioTemplateRepo) GetAll(ctx context.Context, accessPointTypeID uuid.UUID, limit, offset int) (accessPointRadioTemplates []*entity.AccessPointRadioTemplate, err error) {
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
			access_point_type_id,
			created_at, updated_at, deleted_at
		FROM access_point_radio_templates 
		WHERE access_point_type_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, accessPointTypeID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve access point radio template")
		return
	}
	defer rows.Close()

	for rows.Next() {
		accessPointRadioTemplate := &entity.AccessPointRadioTemplate{}
		err = rows.Scan(
			&accessPointRadioTemplate.ID,
			&accessPointRadioTemplate.Number,
			&accessPointRadioTemplate.Channel,
			&accessPointRadioTemplate.Channel2,
			&accessPointRadioTemplate.ChannelWidth,
			&accessPointRadioTemplate.WiFi,
			&accessPointRadioTemplate.Power,
			&accessPointRadioTemplate.Bandwidth,
			&accessPointRadioTemplate.GuardInterval,
			&accessPointRadioTemplate.IsActive,
			&accessPointRadioTemplate.AccessPointTypeID,
			&accessPointRadioTemplate.CreatedAt, &accessPointRadioTemplate.UpdatedAt, &accessPointRadioTemplate.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan access point radio template")
			return
		}
		accessPointRadioTemplates = append(accessPointRadioTemplates, accessPointRadioTemplate)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(accessPointRadioTemplates)
	if length == 0 {
		log.Info().Msgf("access point radio templates for accessPointTypeID %v were not found", accessPointTypeID)
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved %d access point radio templates", length)
	return
}

func (r *accessPointRadioTemplateRepo) Update(ctx context.Context, updateAccessPointRadioTemplateDTO *dto.PatchUpdateAccessPointRadioTemplateDTO) (err error) {
	query := "UPDATE access_point_radio_templates SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateAccessPointRadioTemplateDTO.Number != nil {
		updates = append(updates, fmt.Sprintf("number = $%d", paramID))
		params = append(params, updateAccessPointRadioTemplateDTO.Number)
		paramID++
	}
	if updateAccessPointRadioTemplateDTO.Channel != nil {
		updates = append(updates, fmt.Sprintf("channel = $%d", paramID))
		params = append(params, updateAccessPointRadioTemplateDTO.Channel)
		paramID++
	}
	if updateAccessPointRadioTemplateDTO.Channel2 != nil {
		updates = append(updates, fmt.Sprintf("channel2 = $%d", paramID))
		params = append(params, updateAccessPointRadioTemplateDTO.Channel2)
		paramID++
	}
	if updateAccessPointRadioTemplateDTO.ChannelWidth != nil {
		updates = append(updates, fmt.Sprintf("channel_width = $%d", paramID))
		params = append(params, updateAccessPointRadioTemplateDTO.ChannelWidth)
		paramID++
	}
	if updateAccessPointRadioTemplateDTO.WiFi != nil {
		updates = append(updates, fmt.Sprintf("wifi = $%d", paramID))
		params = append(params, updateAccessPointRadioTemplateDTO.WiFi)
		paramID++
	}
	if updateAccessPointRadioTemplateDTO.Power != nil {
		updates = append(updates, fmt.Sprintf("power = $%d", paramID))
		params = append(params, updateAccessPointRadioTemplateDTO.Power)
		paramID++
	}
	if updateAccessPointRadioTemplateDTO.Bandwidth != nil {
		updates = append(updates, fmt.Sprintf("bandwidth = $%d", paramID))
		params = append(params, updateAccessPointRadioTemplateDTO.Bandwidth)
		paramID++
	}
	if updateAccessPointRadioTemplateDTO.GuardInterval != nil {
		updates = append(updates, fmt.Sprintf("guard_interval = $%d", paramID))
		params = append(params, updateAccessPointRadioTemplateDTO.GuardInterval)
		paramID++
	}
	if updateAccessPointRadioTemplateDTO.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", paramID))
		params = append(params, updateAccessPointRadioTemplateDTO.IsActive)
		paramID++
	}
	if updateAccessPointRadioTemplateDTO.AccessPointTypeID != nil {
		updates = append(updates, fmt.Sprintf("access_point_type_id = $%d", paramID))
		params = append(params, updateAccessPointRadioTemplateDTO.AccessPointTypeID)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return service.ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateAccessPointRadioTemplateDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no access point radio template found with the ID: %v", updateAccessPointRadioTemplateDTO.ID)
		return service.ErrNotFound
	}

	return
}

// Checks if the accessPointRadioTemplate has been soft deleted
func (r *accessPointRadioTemplateRepo) IsAccessPointRadioTemplateSoftDeleted(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM access_point_radio_templates WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, accessPointRadioTemplateID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msgf("no access point radio template found with ID %v", accessPointRadioTemplateID)
			return false, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve access point radio template")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is access point radio template deleted: %v", isDeleted)
	return
}

func (r *accessPointRadioTemplateRepo) SoftDelete(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (err error) {
	query := `UPDATE access_point_radio_templates SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, accessPointRadioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete access point radio template")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no access point radio template found with the ID: %v", accessPointRadioTemplateID)
		return service.ErrNotFound
	}

	log.Debug().Msg("access point radio template deleted_at timestamp updated successfully")
	return
}

func (r *accessPointRadioTemplateRepo) Restore(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (err error) {
	query := `UPDATE access_point_radio_templates SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, accessPointRadioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore access point radio template")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no access point radio template found with the ID: %v", accessPointRadioTemplateID)
		return service.ErrNotFound
	}

	log.Debug().Msg("access point radio template deleted_at timestamp set NULL successfully")
	return
}
