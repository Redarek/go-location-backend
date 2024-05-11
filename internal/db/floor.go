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

// CreateFloor creates a floor
func (p *postgres) CreateFloor(f *Floor) (id uuid.UUID, err error) {
	query := `INSERT INTO floors (name, number, scale, building_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, f.Name, f.Number, f.Scale, f.BuildingID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create floor")
	}
	return
}

// GetFloor retrieves a floor
func (p *postgres) GetFloor(floorUUID uuid.UUID) (f *Floor, err error) {
	query := `SELECT * FROM floors WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, floorUUID)
	f = &Floor{}
	err = row.Scan(&f.ID, &f.Name, &f.Number, &f.Image, &f.WidthInPixels, &f.HeightInPixels, &f.Scale, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt, &f.BuildingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No floor found with uuid %v", floorUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve floor")
		return
	}
	log.Debug().Msgf("Retrieved floor: %v", f)
	return
}

// IsFloorSoftDeleted checks if the floor has been soft deleted
func (p *postgres) IsFloorSoftDeleted(floorUUID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM floors WHERE id = $1`
	row := p.Pool.QueryRow(context.Background(), query, floorUUID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msgf("No floor found with uuid %v", floorUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve floor")
		return
	}
	isDeleted = deletedAt.Valid
	log.Debug().Msgf("Is floor deleted: %v", isDeleted)
	return
}

// GetFloors retrieves floors
func (p *postgres) GetFloors(buildingUUID uuid.UUID) (fs []*Floor, err error) {
	query := `SELECT * FROM floors WHERE building_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, buildingUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve floors")
		return
	}
	defer rows.Close()

	var f *Floor
	for rows.Next() {
		f = new(Floor)
		err = rows.Scan(&f.ID, &f.Name, &f.Number, &f.Image, &f.WidthInPixels, &f.HeightInPixels, &f.Scale, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt, &f.BuildingID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan floor")
			return
		}
		fs = append(fs, f)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d floors", len(fs))
	return
}

// RestoreFloor restore a floor
func (p *postgres) SoftDeleteFloor(floorUUID uuid.UUID) (err error) {
	query := `UPDATE floors SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to soft delete floor")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No floor found with the uuid: %v", floorUUID)
		return
	}
	log.Debug().Msg("Floor deleted_at timestamp updated successfully")
	return
}

// SoftDeleteFloor soft delete a floor
func (p *postgres) RestoreFloor(floorUUID uuid.UUID) (err error) {
	query := `UPDATE floors SET deleted_at = NULL WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to restore floor")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No floor found with the uuid: %v", floorUUID)
		return
	}
	log.Debug().Msg("Floor deleted_at timestamp set null successfully")
	return
}

// PatchUpdateFloor updates only the specified fields of a floor
func (p *postgres) PatchUpdateFloor(f *Floor) (err error) {
	query := "UPDATE floors SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if f.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, f.Name)
		paramID++
	}
	if f.Number != nil {
		updates = append(updates, fmt.Sprintf("number = $%d", paramID))
		params = append(params, f.Number)
		paramID++
	}
	if f.Image != nil {
		updates = append(updates, fmt.Sprintf("image = $%d", paramID))
		params = append(params, f.Image)
		paramID++
	}
	if f.WidthInPixels != nil {
		updates = append(updates, fmt.Sprintf("width_in_pixels = $%d", paramID))
		params = append(params, f.WidthInPixels)
		paramID++
	}
	if f.HeightInPixels != nil {
		updates = append(updates, fmt.Sprintf("height_in_pixels = $%d", paramID))
		params = append(params, f.HeightInPixels)
		paramID++
	}
	if f.Scale != nil {
		updates = append(updates, fmt.Sprintf("scale = $%d", paramID))
		params = append(params, f.Scale)
		paramID++
	}

	if len(updates) == 0 {
		log.Error().Msg("No fields provided for update")
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, f.ID)

	_, err = p.Pool.Exec(context.Background(), query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute update")
		return
	}

	return
}
