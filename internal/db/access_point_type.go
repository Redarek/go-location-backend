package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// CreateAccessPointType creates a access point type
func (p *postgres) CreateAccessPointType(apt AccessPointType) (id int, err error) {
	query := `INSERT INTO access_point_types (site_id)
			VALUES ($1)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, apt.SiteID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create access point type")
	}
	return
}

// GetAccessPointType retrieves an access point type
func (p *postgres) GetAccessPointType(accessPointTypeUUID uuid.UUID) (apt *AccessPointType, err error) {
	query := `SELECT * FROM access_point_types WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, accessPointTypeUUID)
	apt = &AccessPointType{}
	err = row.Scan(&apt.ID, &apt.CreatedAt, &apt.UpdatedAt, &apt.DeletedAt, &apt.SiteID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No access point type found with uuid %v", accessPointTypeUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve access point type")
		return
	}
	log.Debug().Msgf("Retrieved access point type: %v", apt)
	return
}

// IsAccessPointTypeSoftDeleted checks if the access point type has been soft deleted
func (p *postgres) IsAccessPointTypeSoftDeleted(accessPointTypeUUID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM access_point_types WHERE id = $1`
	row := p.Pool.QueryRow(context.Background(), query, accessPointTypeUUID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No access point type found with uuid %v", accessPointTypeUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve access point type")
		return
	}
	log.Debug().Msgf("Is access point type deleted: %v", deletedAt.Valid)
	isDeleted = deletedAt.Valid
	return
}

// GetAccessPointTypes retrieves access point types
func (p *postgres) GetAccessPointTypes(siteUUID uuid.UUID) (apts []*AccessPointType, err error) {
	query := `SELECT * FROM access_point_types WHERE site_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve access point types")
		return
	}
	defer rows.Close()

	var apt *AccessPointType
	for rows.Next() {
		apt = new(AccessPointType)
		err = rows.Scan(&apt.ID, &apt.CreatedAt, &apt.UpdatedAt, &apt.DeletedAt, &apt.SiteID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan access point types")
			return
		}
		apts = append(apts, apt)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d access point types", len(apts))
	return
}

// SoftDeleteAccessPointType soft delete an access point type
func (p *postgres) SoftDeleteAccessPointType(accessPointTypeUUID uuid.UUID) (err error) {
	query := `UPDATE access_point_types SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, accessPointTypeUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to soft delete access point type")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No access point type found with the uuid: %v", accessPointTypeUUID)
		return
	}
	log.Debug().Msg("Access point type deleted_at timestamp updated successfully")
	return
}

// RestoreAccessPointType restore an access point type
func (p *postgres) RestoreAccessPointType(accessPointTypeUUID uuid.UUID) (err error) {
	query := `UPDATE access_point_types SET deleted_at = NULL WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, accessPointTypeUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to restore access point type")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No access point type found with the uuid: %v", accessPointTypeUUID)
		return
	}
	log.Debug().Msg("Access point type deleted_at timestamp set null successfully")
	return
}
