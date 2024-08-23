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

// CreateAccessPoint creates an access point
func (p *postgres) CreateAccessPoint(ap *AccessPoint) (id uuid.UUID, err error) {
	query := `INSERT INTO access_points (name, x, y, z, floor_id, access_point_type_id)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, ap.Name, ap.X, ap.Y, ap.Z, ap.FloorID, ap.AccessPointTypeID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create access point")
	}
	return
}

// GetAccessPoint retrieves an access point
func (p *postgres) GetAccessPoint(accessPointUUID uuid.UUID) (ap *AccessPoint, err error) {
	query := `SELECT 
			id, 
			name,
			x, y, z,
			access_point_type_id,
			floor_id,
			created_at, updated_at, deleted_at
		FROM access_points WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, accessPointUUID)
	ap = &AccessPoint{}
	err = row.Scan(&ap.ID, &ap.Name, &ap.X, &ap.Y, &ap.Z, &ap.AccessPointTypeID, &ap.FloorID, &ap.CreatedAt, &ap.UpdatedAt, &ap.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No access point found with uuid %v", accessPointUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve access point")
		return
	}
	log.Debug().Msgf("Retrieved access point: %v", ap)
	return
}

// GetAccessPointDetailed retrieves an access point detailed
func (p *postgres) GetAccessPointDetailed(accessPointUUID uuid.UUID) (ap *AccessPointDetailed, err error) {
	query := `
	SELECT ap.id, ap.name, ap.x, ap.y, ap.z, ap.created_at, ap.updated_at, ap.deleted_at, ap.floor_id, ap.access_point_type_id, apt.id, apt.name, apt.color, apt.created_at, apt.updated_at, apt.deleted_at, apt.site_id, r.id, r.number, r.channel, r.wifi, r.power, r.bandwidth, r.guard_interval, r.is_active, r.created_at, r.updated_at, r.deleted_at, r.access_point_id
	FROM access_points ap
	LEFT JOIN access_point_types apt ON ap.access_point_type_id = apt.id AND apt.deleted_at IS NULL
	LEFT JOIN radios r ON ap.id = r.access_point_id AND r.deleted_at IS NULL
	WHERE ap.id = $1 AND ap.deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, accessPointUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve access point")
		return
	}
	defer rows.Close()

	ap = new(AccessPointDetailed)
	apt := new(AccessPointType)

	for rows.Next() {
		r := new(Radio)

		err = rows.Scan(
			&ap.ID, &ap.Name, &ap.X, &ap.Y, &ap.Z, &ap.CreatedAt, &ap.UpdatedAt, &ap.DeletedAt, &ap.FloorID, &ap.AccessPointTypeID,
			&apt.ID, &apt.Name, &apt.Color, &apt.CreatedAt, &apt.UpdatedAt, &apt.DeletedAt, &apt.SiteID,
			&r.ID, &r.Number, &r.Channel, &r.WiFi, &r.Power, &r.Bandwidth, &r.GuardInterval, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt, &r.AccessPointID,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan access point and related data")
			return
		}
		ap.AccessPointType = apt
		ap.Radios = append(ap.Radios, r)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved access point with detailed info: %v", ap)
	return
}

// IsAccessPointSoftDeleted checks if the access point has been soft deleted
func (p *postgres) IsAccessPointSoftDeleted(accessPointUUID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM access_points WHERE id = $1`
	row := p.Pool.QueryRow(context.Background(), query, accessPointUUID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No access point found with uuid %v", accessPointUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve access point")
		return
	}
	isDeleted = deletedAt.Valid
	log.Debug().Msgf("Is access point deleted: %v", isDeleted)
	return
}

// GetAccessPoints retrieves access points
func (p *postgres) GetAccessPoints(floorUUID uuid.UUID) (aps []*AccessPoint, err error) {
	query := `SELECT
			id, 
			name,
			x, y, z,
			access_point_type_id,
			floor_id,
			created_at, updated_at, deleted_at 
		FROM access_points WHERE floor_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve access points")
		return
	}
	defer rows.Close()

	var ap *AccessPoint
	for rows.Next() {
		ap = new(AccessPoint)
		err = rows.Scan(&ap.ID, &ap.Name, &ap.X, &ap.Y, &ap.Z, &ap.AccessPointTypeID, &ap.FloorID, &ap.CreatedAt, &ap.UpdatedAt, &ap.DeletedAt)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan access point")
			return
		}
		aps = append(aps, ap)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d access points", len(aps))
	return
}

func (p *postgres) GetAccessPointsDetailed(floorUUID uuid.UUID) (aps []*AccessPointDetailed, err error) {
	query := `SELECT 
			ap.id, 
			ap.name, 
			ap.x, ap.y, ap.z, 
			ap.created_at, ap.updated_at, ap.deleted_at, 
			ap.floor_id, 
			ap.access_point_type_id, 
			
			apt.id, 
			apt.name, 
			apt.color, 
			apt.created_at, apt.updated_at, apt.deleted_at, 
			apt.site_id, 
			r.id, r.number, 
			r.channel, 
			r.wifi, 
			r.power, 
			r.bandwidth, 
			r.guard_interval, 
			r.is_active, 
			r.created_at, 
			r.updated_at, 
			r.deleted_at, 
			r.access_point_id
		FROM access_points ap
		LEFT JOIN access_point_types apt ON ap.access_point_type_id = apt.id AND apt.deleted_at IS NULL
		LEFT JOIN radios r ON ap.id = r.access_point_id AND r.deleted_at IS NULL
		WHERE ap.floor_id = $1 AND ap.deleted_at IS NULL
		GROUP BY ap.id, ap.name, ap.x, ap.y, ap.z, ap.created_at, ap.updated_at, ap.deleted_at, ap.floor_id, ap.access_point_type_id, apt.id, r.id`
	rows, err := p.Pool.Query(context.Background(), query, floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve access points")
		return
	}
	defer rows.Close()

	apMap := make(map[uuid.UUID]*AccessPointDetailed) // Map to track access points and avoid duplicates

	for rows.Next() {
		ap := new(AccessPointDetailed)
		apt := new(AccessPointType)
		r := new(Radio)

		err = rows.Scan(
			&ap.ID, &ap.Name, &ap.X, &ap.Y, &ap.Z, &ap.CreatedAt, &ap.UpdatedAt, &ap.DeletedAt, &ap.FloorID, &ap.AccessPointTypeID,
			&apt.ID, &apt.Name, &apt.Color, &apt.CreatedAt, &apt.UpdatedAt, &apt.DeletedAt, &apt.SiteID,
			&r.ID, &r.Number, &r.Channel, &r.WiFi, &r.Power, &r.Bandwidth, &r.GuardInterval, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt, &r.AccessPointID,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan access points and related data")
			return
		}

		if existingAP, exists := apMap[ap.ID]; exists {
			// If access point is already in the map, append the new radio to its list
			existingAP.Radios = append(existingAP.Radios, r)
		} else {
			// If it's a new access point, initialize and add to map
			ap.AccessPointType = apt
			ap.Radios = append(ap.Radios, r)
			apMap[ap.ID] = ap
		}
	}

	// Convert map to slice
	for _, ap := range apMap {
		aps = append(aps, ap)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d unique access points with detailed info", len(aps))
	return
}

// SoftDeleteAccessPoint soft delete an access point
func (p *postgres) SoftDeleteAccessPoint(accessPointUUID uuid.UUID) (err error) {
	query := `UPDATE access_points SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, accessPointUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to soft delete access point")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No access point found with the uuid: %v", accessPointUUID)
		return
	}
	log.Debug().Msg("Access point deleted_at timestamp updated successfully")
	return
}

// RestoreAccessPoint restore an access point
func (p *postgres) RestoreAccessPoint(accessPointUUID uuid.UUID) (err error) {
	query := `UPDATE access_points SET deleted_at = NULL WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, accessPointUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to restore access point")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No access point found with the uuid: %v", accessPointUUID)
		return
	}
	log.Debug().Msg("Access point deleted_at timestamp set null successfully")
	return
}

// PatchUpdateAccessPoint updates only the specified fields of an access point
func (p *postgres) PatchUpdateAccessPoint(ap *AccessPoint) (err error) {
	query := "UPDATE access_points SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if ap.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, ap.Name)
		paramID++
	}
	if ap.X != nil {
		updates = append(updates, fmt.Sprintf("x = $%d", paramID))
		params = append(params, ap.X)
		paramID++
	}
	if ap.Y != nil {
		updates = append(updates, fmt.Sprintf("y = $%d", paramID))
		params = append(params, ap.Y)
		paramID++
	}
	if ap.Z != nil {
		updates = append(updates, fmt.Sprintf("z = $%d", paramID))
		params = append(params, ap.Z)
		paramID++
	}

	if len(updates) == 0 {
		log.Error().Msg("No fields provided for update")
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, ap.ID)

	_, err = p.Pool.Exec(context.Background(), query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute update")
		return
	}

	return
}

//func (p *postgres) SetRadioState(rs *RadioState) (id uuid.UUID, err error) {
//	query := `
//    INSERT INTO radio_states (access_point_id, radio_id, is_active)
//    VALUES ($1, $2, $3)
//    ON CONFLICT (access_point_id, radio_id)
//    DO UPDATE SET is_active = EXCLUDED.is_active;
//    `
//	row := p.Pool.QueryRow(context.Background(), query, rs.AccessPointID, rs.RadioID, rs.IsActive)
//	err = row.Scan(&id)
//	if err != nil {
//		log.Error().Err(err).Msg("Failed to set radio state")
//	}
//	return
//}

//
//func (p *postgres) GetRadioStates(accessPointID uuid.UUID) (radioStates []RadioState, err error) {
//	query := `
//    SELECT r.radio_id, r.is_active
//    FROM radios r
//    INNER JOIN radio_states rs ON rs.radio_id = r.id
//    WHERE rs.access_point_id = $1;
//    `
//
//	rows, err := p.Pool.Query(context.Background(), query, accessPointID)
//	if err != nil {
//		log.Error().Err(err).Msg("Unable to get radio states")
//		return
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var state RadioState
//		if err = rows.Scan(&state.RadioID, &state.IsActive); err != nil {
//			return
//		}
//		radioStates = append(radioStates, state)
//	}
//	return
//
//}
