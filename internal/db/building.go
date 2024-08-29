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

// CreateBuilding creates a building
func (p *postgres) CreateBuilding(b *Building) (id uuid.UUID, err error) {
	query := `INSERT INTO buildings (name, description, country, city, address, site_id)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, b.Name, b.Description, b.Country, b.City, b.Address, b.SiteID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create building")
	}
	return
}

// GetBuilding retrieves a building
func (p *postgres) GetBuilding(buildingUUID uuid.UUID) (b *Building, err error) {
	query := `SELECT 
			id, 
			name, 
			description, 
			country,
			city,
			address,
			site_id,
			created_at, updated_at, deleted_at
		FROM buildings WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, buildingUUID)
	b = &Building{}
	err = row.Scan(&b.ID, &b.Name, &b.Description, &b.Country, &b.City, &b.Address, &b.SiteID, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No building found with uuid %v", buildingUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve building")
		return
	}
	log.Debug().Msgf("Retrieved building: %v", b)
	return
}

// IsBuildingSoftDeleted checks if the building has been soft deleted
func (p *postgres) IsBuildingSoftDeleted(buildingUUID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM buildings WHERE id = $1`
	row := p.Pool.QueryRow(context.Background(), query, buildingUUID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msgf("No building found with uuid %v", buildingUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve building")
		return
	}
	isDeleted = deletedAt.Valid
	log.Debug().Msgf("Is building deleted: %v", isDeleted)
	return
}

// GetBuildings retrieves buildings
func (p *postgres) GetBuildings(siteUUID uuid.UUID) (bs []*Building, err error) {
	query := `SELECT
			id, 
			name, 
			description, 
			country,
			city,
			address,
			site_id,
			created_at, updated_at, deleted_at
		FROM buildings WHERE site_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve buildings")
		return
	}
	defer rows.Close()

	var b *Building
	for rows.Next() {
		b = new(Building)
		err = rows.Scan(&b.ID, &b.Name, &b.Description, &b.Country, &b.City, &b.Address, &b.SiteID, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan building")
			return
		}
		bs = append(bs, b)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d buildings", len(bs))
	return
}

// SoftDeleteBuilding soft delete a building
func (p *postgres) SoftDeleteBuilding(buildingUUID uuid.UUID) (err error) {
	query := `UPDATE buildings SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, buildingUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to soft delete building")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No building found with the uuid: %v", buildingUUID)
		return
	}
	log.Debug().Msg("Building deleted_at timestamp updated successfully")
	return
}

// RestoreBuilding restore a building
func (p *postgres) RestoreBuilding(buildingUUID uuid.UUID) (err error) {
	query := `UPDATE buildings SET deleted_at = NULL WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, buildingUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to restore building")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No building found with the uuid: %v", buildingUUID)
		return
	}
	log.Debug().Msg("Building deleted_at timestamp set null successfully")
	return
}

// PatchUpdateBuilding updates only the specified fields of a building
func (p *postgres) PatchUpdateBuilding(b *Building) (err error) {
	query := "UPDATE buildings SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if b.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, b.Name)
		paramID++
	}
	if b.Description != "" {
		updates = append(updates, fmt.Sprintf("description = $%d", paramID))
		params = append(params, b.Description)
		paramID++
	}
	if b.Country != "" {
		updates = append(updates, fmt.Sprintf("country = $%d", paramID))
		params = append(params, b.Country)
		paramID++
	}
	if b.City != "" {
		updates = append(updates, fmt.Sprintf("city = $%d", paramID))
		params = append(params, b.City)
		paramID++
	}
	if b.Address != "" {
		updates = append(updates, fmt.Sprintf("address = $%d", paramID))
		params = append(params, b.Address)
		paramID++
	}

	if len(updates) == 0 {
		log.Error().Msg("No fields provided for update")
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, b.ID)

	_, err = p.Pool.Exec(context.Background(), query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute update")
		return
	}

	return
}
