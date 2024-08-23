package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	. "location-backend/internal/db/model"
)

// CreateRadioTemplate creates a radio template
func (p *postgres) CreateRadioTemplate(rt *RadioTemplate) (id uuid.UUID, err error) {
	query := `INSERT INTO radio_templates (number, channel, wifi, power, bandwidth, guard_interval, access_point_type_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, rt.Number, rt.Channel, rt.WiFi, rt.Power, rt.Bandwidth, rt.GuardInterval, rt.AccessPointTypeID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create radio template")
	}
	return
}

// GetRadioTemplate retrieves a radio template
func (p *postgres) GetRadioTemplate(radioUUID uuid.UUID) (rt RadioTemplate, err error) {
	query := `SELECT 
			id,
			number,
			channel,
			wifi,
			power,
			bandwidth,
			guard_interval,
			access_point_type_id,
			created_at, updated_at, deleted_at
		FROM radio_templates WHERE id=$1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, radioUUID)
	err = row.Scan(&rt.ID, &rt.Number, &rt.Channel, &rt.WiFi, &rt.Power, &rt.Bandwidth, &rt.GuardInterval, &rt.AccessPointTypeID, &rt.CreatedAt, &rt.UpdatedAt, &rt.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No radio template found with ID %v", radioUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve radio template")
		return
	}
	log.Debug().Msgf("Retrieved radio template: %v", rt)
	return
}

// IsRadioTemplateSoftDeleted checks if the radio template has been soft deleted
func (p *postgres) IsRadioTemplateSoftDeleted(radioUUID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM radio_templates WHERE id = $1`
	row := p.Pool.QueryRow(context.Background(), query, radioUUID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No radio template found with uuid %v", radioUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve radio template")
		return
	}
	isDeleted = deletedAt.Valid
	log.Debug().Msgf("Is radio template deleted: %v", isDeleted)
	return
}

// GetRadioTemplates retrieves radio templates
func (p *postgres) GetRadioTemplates(accessPointTypeID uuid.UUID) (rs []*RadioTemplate, err error) {
	query := `SELECT 
			id,
			number,
			channel,
			wifi,
			power,
			bandwidth,
			guard_interval,
			access_point_type_id,
			created_at, updated_at, deleted_at
		FROM radio_templates WHERE access_point_type_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, accessPointTypeID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve radio templates")
		return
	}
	defer rows.Close()

	var r *RadioTemplate
	for rows.Next() {
		r = new(RadioTemplate)
		err = rows.Scan(&r.ID, &r.Number, &r.Channel, &r.WiFi, &r.Power, &r.Bandwidth, &r.GuardInterval, &r.AccessPointTypeID, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan radio templates")
			return
		}
		rs = append(rs, r)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d radio templates", len(rs))
	return
}

//// GetRadios retrieves radio templates
//func (p *postgres) GetRadios(accessPointTypeUUID uuid.UUID) (rs []*Radio, err error) {
//	//query := `SELECT * FROM radios WHERE access_point_type_id = $1 AND deleted_at IS NULL`
//	query := `
//		SELECT r.id, r.number, r.channel, r.wifi, r.power, r.bandwidth, r.guard_interval, r.created_at, r.updated_at, r.deleted_at, r.access_point_type_id, rs.is_active
//		FROM radios r
//		LEFT JOIN radio_states rs ON rs.radio_id = r.id AND rs.access_point_id = $1
//		WHERE r.access_point_type_id = (SELECT access_point_type_id FROM access_points WHERE id = $1) AND r.deleted_at IS NULL
//		`
//	rows, err := p.Pool.Query(context.Background(), query, accessPointTypeUUID)
//	if err != nil {
//		log.Error().Err(err).Msg("Failed to retrieve radios")
//		return
//	}
//	defer rows.Close()
//
//	var r *Radio
//	for rows.Next() {
//		r = new(Radio)
//		var isActive sql.NullBool
//		err = rows.Scan(&r.ID, &r.Number, &r.Channel, &r.WiFi, &r.Power, &r.Bandwidth, &r.GuardInterval, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt, &r.AccessPointTypeID, &isActive)
//		if err != nil {
//			log.Error().Err(err).Msg("Failed to scan radios")
//			return
//		}
//		if isActive.Valid {
//			r.IsActive = isActive.Bool
//		} else {
//			r.IsActive = false // Set default state if not specified
//		}
//		rs = append(rs, r)
//	}
//
//	if err = rows.Err(); err != nil {
//		log.Error().Err(err).Msg("Rows iteration error")
//		return
//	}
//
//	log.Debug().Msgf("Retrieved %d radios", len(rs))
//	return
//}

// SoftDeleteRadioTemplate soft delete a radio template
func (p *postgres) SoftDeleteRadioTemplate(radioUUID uuid.UUID) (err error) {
	query := `UPDATE radio_templates SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, radioUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to soft delete radio template")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No radio template found with the uuid: %v", radioUUID)
		return
	}
	log.Debug().Msg("Radio template deleted_at timestamp updated successfully")
	return
}

// RestoreRadioTemplate restore a radio template
func (p *postgres) RestoreRadioTemplate(radioUUID uuid.UUID) (err error) {
	query := `UPDATE radio_templates SET deleted_at = NULL WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, radioUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to restore radio template")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No radio template found with the uuid: %v", radioUUID)
		return
	}
	log.Debug().Msg("Radio template deleted_at timestamp set null successfully")
	return
}

// PatchUpdateRadioTemplate updates only the specified fields of a radio template
func (p *postgres) PatchUpdateRadioTemplate(r *RadioTemplate) (err error) {
	query := "UPDATE radio_templates SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if r.Number != nil {
		updates = append(updates, fmt.Sprintf("number = $%d", paramID))
		params = append(params, r.Number)
		paramID++
	}
	if r.Channel != nil {
		updates = append(updates, fmt.Sprintf("channel = $%d", paramID))
		params = append(params, r.Channel)
		paramID++
	}
	if r.WiFi != nil {
		updates = append(updates, fmt.Sprintf("wifi = $%d", paramID))
		params = append(params, r.WiFi)
		paramID++
	}
	if r.Power != nil {
		updates = append(updates, fmt.Sprintf("power = $%d", paramID))
		params = append(params, r.Power)
		paramID++
	}
	if r.Bandwidth != nil {
		updates = append(updates, fmt.Sprintf("bandwidth = $%d", paramID))
		params = append(params, r.Bandwidth)
		paramID++
	}
	if r.GuardInterval != nil {
		updates = append(updates, fmt.Sprintf("guard_interval = $%d", paramID))
		params = append(params, r.GuardInterval)
		paramID++
	}

	if len(updates) == 0 {
		log.Error().Msg("No fields provided for update")
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, r.ID)

	_, err = p.Pool.Exec(context.Background(), query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute update")
		return
	}

	return
}
