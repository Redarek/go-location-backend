package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	. "location-backend/internal/db/models"
)

// CreateWallType creates a wall type
func (p *postgres) CreateWallType(wt *WallType) (id uuid.UUID, err error) {
	sql := `INSERT INTO wall_types (name, color, attenuation_24, attenuation_5, attenuation_6, thickness, site_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), sql, wt.Name, wt.Color, wt.Attenuation24, wt.Attenuation5, wt.Attenuation6, wt.Thickness, wt.SiteID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create wall type")
	}
	return
}

// GetWallType retrieves a wall type
func (p *postgres) GetWallType(wallTypeUUID uuid.UUID) (wt *WallType, err error) {
	sql := `SELECT * FROM wall_types WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), sql, wallTypeUUID)
	wt = &WallType{}
	err = row.Scan(&wt.ID, &wt.Name, &wt.Color, &wt.Attenuation24, &wt.Attenuation5, &wt.Attenuation6, &wt.Thickness, &wt.CreatedAt, &wt.UpdatedAt, &wt.DeletedAt, &wt.SiteID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msgf("No wall type found with uuid %v", wallTypeUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve wall type")
		return
	}
	log.Debug().Msgf("Retrieved wall type: %v", wt)
	return
}

// IsWallTypeSoftDeleted checks if the wall type has been soft deleted
func (p *postgres) IsWallTypeSoftDeleted(wallTypeUUID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	sql := `SELECT deleted_at FROM wall_types WHERE id = $1`
	row := p.Pool.QueryRow(context.Background(), sql, wallTypeUUID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msgf("No wall type found with uuid %v", wallTypeUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve wall type")
		return
	}
	log.Debug().Msgf("Is wall type deleted: %v", deletedAt.Valid)
	isDeleted = deletedAt.Valid
	return
}

// GetWallTypes retrieves wall types
func (p *postgres) GetWallTypes(siteUUID uuid.UUID) (wts []*WallType, err error) {
	sql := `SELECT * FROM wall_types WHERE site_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), sql, siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve wall types")
		return
	}
	defer rows.Close()

	var wt *WallType
	for rows.Next() {
		wt = new(WallType)
		err = rows.Scan(&wt.ID, &wt.Name, &wt.Color, &wt.Attenuation24, &wt.Attenuation5, &wt.Attenuation6, &wt.Thickness, &wt.CreatedAt, &wt.UpdatedAt, &wt.DeletedAt, &wt.SiteID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan wall type")
			return
		}
		wts = append(wts, wt)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d wall types", len(wts))
	return
}

// SoftDeleteWallType soft delete a wall type
func (p *postgres) SoftDeleteWallType(wallTypeUUID uuid.UUID) (err error) {
	sql := `UPDATE wall_types SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), sql, wallTypeUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to soft delete wall type")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No wall type found with the uuid: %v", wallTypeUUID)
		return
	}
	log.Debug().Msg("Wall type deleted_at timestamp updated successfully")
	return
}

// RestoreWallType restore a wall type
func (p *postgres) RestoreWallType(wallTypeUUID uuid.UUID) (err error) {
	sql := `UPDATE wall_types SET deleted_at = NULL WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), sql, wallTypeUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to restore wall type")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No wall type found with the uuid: %v", wallTypeUUID)
		return
	}
	log.Debug().Msg("Wall type deleted_at timestamp set null successfully")
	return
}

// PatchUpdateWallType updates only the specified fields of a wall type
func (p *postgres) PatchUpdateWallType(wt *WallType) (err error) {
	log.Debug().Msgf("Patching wall type: %v", wt)
	query := "UPDATE wall_types SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if wt.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, wt.Name)
		paramID++
	}
	if wt.Color != "" {
		updates = append(updates, fmt.Sprintf("color = $%d", paramID))
		params = append(params, wt.Color)
		paramID++
	}
	if wt.Attenuation24 != nil {
		updates = append(updates, fmt.Sprintf("attenuation_24 = $%d", paramID))
		params = append(params, wt.Attenuation24)
		paramID++
	}
	if wt.Attenuation5 != nil {
		updates = append(updates, fmt.Sprintf("attenuation_5 = $%d", paramID))
		params = append(params, wt.Attenuation5)
		paramID++
	}
	if wt.Attenuation6 != nil {
		updates = append(updates, fmt.Sprintf("attenuation_6 = $%d", paramID))
		params = append(params, wt.Attenuation6)
		paramID++
	}
	if wt.Thickness != nil {
		updates = append(updates, fmt.Sprintf("thickness = $%d", paramID))
		params = append(params, wt.Thickness)
		paramID++
	}

	if len(updates) == 0 {
		log.Error().Msg("No fields provided for update")
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, wt.ID)

	_, err = p.Pool.Exec(context.Background(), query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute update")
		return
	}

	return
}
