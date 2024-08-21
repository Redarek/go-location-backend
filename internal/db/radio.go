package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type Radio struct {
	ID            uuid.UUID          `json:"id" db:"id"`
	Number        *int               `json:"number" db:"number"`
	Channel       *int               `json:"channel" db:"channel"`
	WiFi          *string            `json:"wifi" db:"wifi"`
	Power         *int               `json:"power" db:"power"`
	Bandwidth     *string            `json:"bandwidth" db:"bandwidth"`
	GuardInterval *int               `json:"guardInterval" db:"guard_interval"`
	IsActive      *bool              `json:"isActive" db:"is_active"`
	CreatedAt     pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt     pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
	DeletedAt     pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
	AccessPointID uuid.UUID          `json:"accessPointId" db:"access_point_id"`
}

//type RadioState struct {
//	AccessPointID uuid.UUID `json:"accessPointId" db:"access_point_id"`
//	RadioID       uuid.UUID `json:"radioId" db:"radio_id"`
//	IsActive      bool      `json:"isActive" db:"is_active"`
//}

// CreateRadio creates a radio
func (p *postgres) CreateRadio(r *Radio) (id uuid.UUID, err error) {
	query := `INSERT INTO radios (number, channel, wifi, power, bandwidth, guard_interval, is_active, access_point_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, r.Number, r.Channel, r.WiFi, r.Power, r.Bandwidth, r.GuardInterval, r.IsActive, r.AccessPointID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create radio")
	}
	return
}

// GetRadio retrieves a radio
func (p *postgres) GetRadio(radioUUID uuid.UUID) (r Radio, err error) {
	query := `SELECT * FROM radios WHERE id=$1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, radioUUID)
	err = row.Scan(&r.ID, &r.Number, &r.Channel, &r.WiFi, &r.Power, &r.Bandwidth, &r.GuardInterval, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt, &r.AccessPointID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No radio found with ID %v", radioUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve radio")
		return
	}
	log.Debug().Msgf("Retrieved radio: %v", r)
	return
}

// IsRadioSoftDeleted checks if the radio has been soft deleted
func (p *postgres) IsRadioSoftDeleted(radioUUID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM radios WHERE id = $1`
	row := p.Pool.QueryRow(context.Background(), query, radioUUID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No radio found with uuid %v", radioUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve radio")
		return
	}
	isDeleted = deletedAt.Valid
	log.Debug().Msgf("Is radio deleted: %v", isDeleted)
	return
}

// GetRadios retrieves radios
func (p *postgres) GetRadios(accessPointID uuid.UUID) (rs []*Radio, err error) {
	query := `SELECT * FROM radios WHERE access_point_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve radios")
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := new(Radio)
		err = rows.Scan(&r.ID, &r.Number, &r.Channel, &r.WiFi, &r.Power, &r.Bandwidth, &r.GuardInterval, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt, &r.AccessPointID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan radios")
			return
		}
		rs = append(rs, r)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d radios", len(rs))
	return
}

//// GetRadios retrieves radios
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

// SoftDeleteRadio soft delete a radio
func (p *postgres) SoftDeleteRadio(radioUUID uuid.UUID) (err error) {
	query := `UPDATE radios SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, radioUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to soft delete radio")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No radio found with the uuid: %v", radioUUID)
		return
	}
	log.Debug().Msg("Access point deleted_at timestamp updated successfully")
	return
}

// RestoreRadio restore a radio
func (p *postgres) RestoreRadio(radioUUID uuid.UUID) (err error) {
	query := `UPDATE radios SET deleted_at = NULL WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, radioUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to restore radio")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No radio found with the uuid: %v", radioUUID)
		return
	}
	log.Debug().Msg("Radio deleted_at timestamp set null successfully")
	return
}

// PatchUpdateRadio updates only the specified fields of a radio
func (p *postgres) PatchUpdateRadio(r *Radio) (err error) {
	query := "UPDATE radios SET updated_at = NOW(), "
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
	if r.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", paramID))
		params = append(params, r.IsActive)
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
