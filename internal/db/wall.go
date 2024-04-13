package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"strings"
)

// CreateWall creates a wall
func (p *postgres) CreateWall(w *Wall) (id int, err error) {
	sql := `INSERT INTO walls (x1, y1, x2, y2, floor_id, wall_type_id)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), sql, w.X1, w.Y1, w.X2, w.Y2, w.FloorID, w.WallTypeID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create wall")
	}
	return
}

// GetWall retrieves a wall
func (p *postgres) GetWall(wallUUID uuid.UUID) (w *Wall, err error) {
	sql := `SELECT * FROM walls WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), sql, wallUUID)
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
	sql := `SELECT deleted_at FROM walls WHERE id = $1`
	row := p.Pool.QueryRow(context.Background(), sql, wallUUID)
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
	sql := `SELECT * FROM walls WHERE floor_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), sql, floorUUID)
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

// SoftDeleteWall soft delete a wall
func (p *postgres) SoftDeleteWall(wallUUID uuid.UUID) (err error) {
	sql := `UPDATE walls SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), sql, wallUUID)
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
	sql := `UPDATE walls SET deleted_at = NULL WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), sql, wallUUID)
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
