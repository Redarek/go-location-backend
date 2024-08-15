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
)

// CreateAccessPointType creates an access point type
func (p *postgres) CreateAccessPointType(apt *AccessPointType) (id uuid.UUID, err error) {
	query := `INSERT INTO access_point_types (name, color, site_id)
			VALUES ($1, $2, $3)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, apt.Name, apt.Color, apt.SiteID)
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
	err = row.Scan(&apt.ID, &apt.Name, &apt.Color, &apt.CreatedAt, &apt.UpdatedAt, &apt.DeletedAt, &apt.SiteID)
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

// GetAccessPointTypeDetailed retrieves an access point type
func (p *postgres) GetAccessPointTypeDetailed(accessPointTypeUUID uuid.UUID) (apt *AccessPointTypeDetailed, err error) {
	query := `
	SELECT apt.id, apt.name, apt.color, apt.created_at, apt.updated_at, apt.deleted_at, apt.site_id,
	       rt.id, rt.number, rt.channel, rt.wifi, rt.power, rt.bandwidth, rt.guard_interval, rt.created_at, rt.updated_at, rt.deleted_at, rt.access_point_type_id
	FROM access_point_types apt
	LEFT JOIN radio_templates rt ON rt.access_point_type_id = apt.id AND rt.deleted_at IS NULL
	WHERE apt.id = $1 AND apt.deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, accessPointTypeUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve access point type")
		return
	}
	defer rows.Close()

	apt = new(AccessPointTypeDetailed)

	for rows.Next() {
		rt := new(RadioTemplate)

		err = rows.Scan(
			&apt.ID, &apt.Name, &apt.Color, &apt.CreatedAt, &apt.UpdatedAt, &apt.DeletedAt, &apt.SiteID,
			&rt.ID, &rt.Number, &rt.Channel, &rt.WiFi, &rt.Power, &rt.Bandwidth, &rt.GuardInterval, &rt.CreatedAt, &rt.UpdatedAt, &rt.DeletedAt, &rt.AccessPointTypeID,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan access point type and related data")
			return
		}

		apt.RadioTemplates = append(apt.RadioTemplates, rt)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved access point type with detailed info: %v", apt)
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
		err = rows.Scan(&apt.ID, &apt.Name, &apt.Color, &apt.CreatedAt, &apt.UpdatedAt, &apt.DeletedAt, &apt.SiteID)
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

func (p *postgres) GetAccessPointTypesDetailed(siteUUID uuid.UUID) (aps []*AccessPointTypeDetailed, err error) {
	query := `
SELECT apt.id, apt.name, apt.color, apt.created_at, apt.updated_at, apt.deleted_at, apt.site_id, r.id, r.number, r.channel, r.wifi, r.power, r.bandwidth, r.guard_interval, r.created_at, r.updated_at, r.deleted_at, r.access_point_type_id
FROM access_point_types apt
LEFT JOIN radio_templates r ON apt.id = r.access_point_type_id AND r.deleted_at IS NULL
WHERE apt.site_id = $1 AND apt.deleted_at IS NULL
`
	rows, err := p.Pool.Query(context.Background(), query, siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve access point types")
		return
	}
	defer rows.Close()

	aptMap := make(map[uuid.UUID]*AccessPointTypeDetailed) // Map to track access points and avoid duplicates

	for rows.Next() {
		apt := new(AccessPointTypeDetailed)
		r := new(RadioTemplate)

		err = rows.Scan(
			&apt.ID, &apt.Name, &apt.Color, &apt.CreatedAt, &apt.UpdatedAt, &apt.DeletedAt, &apt.SiteID,
			&r.ID, &r.Number, &r.Channel, &r.WiFi, &r.Power, &r.Bandwidth, &r.GuardInterval, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt, &r.AccessPointTypeID,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan access points and related data")
			return
		}

		if existingAP, exists := aptMap[apt.ID]; exists {
			// If access point is already in the map, append the new radio to its list
			existingAP.RadioTemplates = append(existingAP.RadioTemplates, r)
		} else {
			// If it's a new access point type, initialize and add to map
			apt.RadioTemplates = append(apt.RadioTemplates, r)
			aptMap[apt.ID] = apt
		}
	}

	// Convert map to slice
	for _, ap := range aptMap {
		aps = append(aps, ap)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d unique access point types with detailed info", len(aps))
	return
}

// Updates AccessPointType
func (p *postgres) PatchUpdateAccessPointType(apt *AccessPointType) (err error) {
	query := "UPDATE access_points SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if apt.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, apt.Name)
		paramID++
	}

	if apt.Color != "" {
		updates = append(updates, fmt.Sprintf("color = $%d", paramID))
		params = append(params, apt.Color)
		paramID++
	}

	if len(updates) == 0 {
		log.Error().Msg("No fields provided for update")
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, apt.ID)

	_, err = p.Pool.Exec(context.Background(), query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute update")
		return
	}

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
