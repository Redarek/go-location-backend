package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"strings"
)

// CreateAccessPoint creates an access point
func (p *postgres) CreateAccessPoint(ap *AccessPoint) (id int, err error) {
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
	query := `SELECT * FROM access_points WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, accessPointUUID)
	ap = &AccessPoint{}
	err = row.Scan(&ap.Name, &ap.X, &ap.Y, &ap.Z, &ap.FloorID, &ap.AccessPointTypeID)
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
	query := `SELECT * FROM access_points WHERE floor_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve access points")
		return
	}
	defer rows.Close()

	var ap *AccessPoint
	for rows.Next() {
		ap = new(AccessPoint)
		err = rows.Scan(&ap.Name, &ap.X, &ap.Y, &ap.Z, &ap.FloorID, &ap.AccessPointTypeID)
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
