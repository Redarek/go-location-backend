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

// CreateSite creates a site
func (p *postgres) CreateSite(userUUID uuid.UUID, s *Site) (id uuid.UUID, err error) {
	query := `INSERT INTO sites (name, description, user_id)
			VALUES ($1, $2, $3)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, s.Name, s.Description, userUUID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create site")
	}
	return
}

// GetSite retrieves a site
func (p *postgres) GetSite(siteUUID uuid.UUID) (s *Site, err error) {
	query := `SELECT * FROM sites WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, siteUUID)
	s = &Site{}
	err = row.Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No site found with uuid %v", siteUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve site")
		return
	}
	log.Debug().Msgf("Retrieved site: %v", s)
	return
}

// IsSiteSoftDeleted checks if the site has been soft deleted
func (p *postgres) IsSiteSoftDeleted(siteUUID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM sites WHERE id = $1`
	row := p.Pool.QueryRow(context.Background(), query, siteUUID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No site found with uuid %v", siteUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve site")
		return
	}
	isDeleted = deletedAt.Valid
	log.Debug().Msgf("Is site deleted: %v", isDeleted)
	return
}

// GetSites retrieves sites
func (p *postgres) GetSites(userUUID uuid.UUID) (sites []*Site, err error) {
	query := `SELECT * FROM sites WHERE user_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, userUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve sites")
		return
	}
	defer rows.Close()

	var s *Site
	for rows.Next() {
		s = new(Site)
		err = rows.Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.UserID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan site")
			return
		}
		sites = append(sites, s)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d sites", len(sites))
	return
}

// SoftDeleteSite soft delete a site
func (p *postgres) SoftDeleteSite(siteUUID uuid.UUID) (err error) {
	query := `UPDATE sites SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to soft delete site")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No site found with the uuid: %v", siteUUID)
		return
	}
	log.Debug().Msg("Site deleted_at timestamp updated successfully")
	return
}

// RestoreSite restore a site
func (p *postgres) RestoreSite(siteUUID uuid.UUID) (err error) {
	query := `UPDATE sites SET deleted_at = NULL WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to restore site")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No site found with the uuid: %v", siteUUID)
		return
	}
	log.Debug().Msg("Site deleted_at timestamp set null successfully")
	return
}

// PatchUpdateSite updates only the specified fields of a site
func (p *postgres) PatchUpdateSite(site *Site) (err error) {
	query := "UPDATE sites SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if site.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, site.Name)
		paramID++
	}
	if site.Description != nil && *site.Description != "" {
		updates = append(updates, fmt.Sprintf("description = $%d", paramID))
		params = append(params, site.Description)
		paramID++
	}

	if len(updates) == 0 {
		log.Error().Msg("No fields provided for update")
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, site.ID)

	_, err = p.Pool.Exec(context.Background(), query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute update")
		return
	}

	return
}
