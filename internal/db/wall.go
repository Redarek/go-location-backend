package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// CreateWall creates a wall
func (p *postgres) CreateWall(w *Wall) (id uuid.UUID, err error) {
	query := `INSERT INTO walls (x1, y1, x2, y2, floor_id, wall_type_id)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, w.X1, w.Y1, w.X2, w.Y2, w.FloorID, w.WallTypeID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create wall")
	}
	return
}

// GetWall retrieves a wall
func (p *postgres) GetWall(wallUUID uuid.UUID) (w *Wall, err error) {
	query := `SELECT * FROM walls WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, wallUUID)
	w = &Wall{}
	err = row.Scan(&w.ID, &w.X1, &w.Y1, &w.X2, &w.Y2, &w.CreatedAt, &w.UpdatedAt, &w.DeletedAt, &w.FloorID, &w.WallTypeID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msgf("No wall found with uuid %v", wallUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve wall")
		return
	}
	log.Debug().Msgf("Retrieved wall: %v", w)
	return
}

// IsWallSoftDeleted checks if the wall has been soft deleted
func (p *postgres) IsWallSoftDeleted(wallUUID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM walls WHERE id = $1`
	row := p.Pool.QueryRow(context.Background(), query, wallUUID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msgf("No wall found with uuid %v", wallUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve wall")
		return
	}
	log.Debug().Msgf("Is wall deleted: %v", deletedAt.Valid)
	isDeleted = deletedAt.Valid
	return
}

// GetWalls retrieves walls
func (p *postgres) GetWalls(floorUUID uuid.UUID) (ws []*Wall, err error) {
	query := `SELECT * FROM walls WHERE floor_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve walls")
		return
	}
	defer rows.Close()

	var w *Wall
	for rows.Next() {
		w = new(Wall)
		err = rows.Scan(&w.ID, &w.X1, &w.Y1, &w.X2, &w.Y2, &w.CreatedAt, &w.UpdatedAt, &w.DeletedAt, &w.FloorID, &w.WallTypeID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan walls")
			return
		}
		ws = append(ws, w)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d walls", len(ws))
	return
}

func (p *postgres) GetWallsDetailed(floorUUID uuid.UUID) (walls []*WallDetailed, err error) {
	query := `
SELECT w.id, w.x1, w.y1, w.x2, w.y2, w.created_at, w.updated_at, w.deleted_at, w.floor_id, w.wall_type_id, wt.id, wt.name, wt.color, wt.attenuation1, wt.attenuation2, wt.attenuation3, wt.thickness, wt.created_at, wt.updated_at, wt.deleted_at, wt.site_id
FROM walls w
LEFT JOIN wall_types wt ON w.wall_type_id = wt.id
WHERE w.floor_id = $1 AND w.deleted_at IS NULL AND wt.deleted_at IS NULL
`
	rows, err := p.Pool.Query(context.Background(), query, floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve walls detailed")
		return
	}
	defer rows.Close()

	wallsMap := make(map[uuid.UUID]*WallDetailed) // Map to track access points and avoid duplicates

	for rows.Next() {
		w := new(WallDetailed)
		wt := new(WallType)

		err = rows.Scan(
			&w.ID, &w.X1, &w.Y1, &w.X2, &w.Y2, &w.CreatedAt, &w.UpdatedAt, &w.DeletedAt, &w.FloorID, &w.WallTypeID,
			&wt.ID, &wt.Name, &wt.Color, &wt.Attenuation24, &wt.Attenuation5, &wt.Attenuation6, &wt.Thickness, &wt.CreatedAt, &wt.UpdatedAt, &wt.DeletedAt, &wt.SiteID,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan walls and related data")
			return
		}

		if _, exists := wallsMap[w.ID]; exists {
			// If wall is already in the map, continue
			continue
		} else {
			// If it's a wall, initialize and add to map
			w.WallType = wt
			wallsMap[w.ID] = w
		}
	}

	// Convert map to slice
	for _, w := range wallsMap {
		walls = append(walls, w)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d unique walls with detailed info", len(walls))
	return
}

// SoftDeleteWall soft delete a wall
func (p *postgres) SoftDeleteWall(wallUUID uuid.UUID) (err error) {
	query := `UPDATE walls SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, wallUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to soft delete wall")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No wall found with the uuid: %v", wallUUID)
		return
	}
	log.Debug().Msg("Wall deleted_at timestamp updated successfully")
	return
}

// RestoreWall restore a wall
func (p *postgres) RestoreWall(wallUUID uuid.UUID) (err error) {
	query := `UPDATE walls SET deleted_at = NULL WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, wallUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to restore wall")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No wall found with the uuid: %v", wallUUID)
		return
	}
	log.Debug().Msg("Wall deleted_at timestamp set null successfully")
	return
}

// PatchUpdateWall updates only the specified fields of a wall
func (p *postgres) PatchUpdateWall(w *Wall) (err error) {
	query := "UPDATE walls SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if w.X1 != nil {
		updates = append(updates, fmt.Sprintf("x1 = $%d", paramID))
		params = append(params, w.X1)
		paramID++
	}
	if w.Y1 != nil {
		updates = append(updates, fmt.Sprintf("y1 = $%d", paramID))
		params = append(params, w.Y1)
		paramID++
	}
	if w.X2 != nil {
		updates = append(updates, fmt.Sprintf("x2 = $%d", paramID))
		params = append(params, w.X2)
		paramID++
	}
	if w.Y2 != nil {
		updates = append(updates, fmt.Sprintf("y2 = $%d", paramID))
		params = append(params, w.Y2)
		paramID++
	}

	if len(updates) == 0 {
		log.Error().Msg("No fields provided for update")
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, w.ID)

	_, err = p.Pool.Exec(context.Background(), query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute update")
		return
	}

	return
}
