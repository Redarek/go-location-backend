package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type wallRepo struct {
	pool *pgxpool.Pool
}

func NewWallRepo(pool *pgxpool.Pool) *wallRepo {
	return &wallRepo{pool: pool}
}

func (r *wallRepo) Create(ctx context.Context, dto *dto.CreateWallDTO) (wallID uuid.UUID, err error) {
	query := `INSERT INTO walls (
			x1, y1, 
			x2, y2, 
			wall_type_id,
			floor_id 
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		dto.X1, dto.Y1,
		dto.X2, dto.Y2,
		dto.WallTypeID,
		dto.FloorID,
	)
	err = row.Scan(&wallID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create wall")
	}
	return
}

func (r *wallRepo) GetOne(ctx context.Context, wallID uuid.UUID) (wall *entity.Wall, err error) {
	query := `SELECT 
		id, 
		x1, y1, 
		x2, y2, 
		wall_type_id,
		floor_id,
		created_at, updated_at, deleted_at
	FROM walls 
	WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, wallID)

	wall = &entity.Wall{}
	err = row.Scan(
		&wall.ID,
		&wall.X1, &wall.Y1,
		&wall.X2, &wall.Y2,
		&wall.WallTypeID,
		&wall.FloorID,
		&wall.CreatedAt, &wall.UpdatedAt, &wall.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no wall found with ID %v", wallID)
			return nil, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve wall")
		return
	}
	log.Debug().Msgf("retrieved wall: %v", wall)
	return
}

// func (r *wallRepo) GetOneDetailed(ctx context.Context, wallID uuid.UUID) (wall *entity.WallDetailed, err error) {

// }

func (r *wallRepo) GetAll(ctx context.Context, floorID uuid.UUID, limit, offset int) (walls []*entity.Wall, err error) {
	query := `SELECT 
			id, 
			x1, y1, 
			x2, y2, 
			wall_type_id,
			floor_id,
			created_at, updated_at, deleted_at
		FROM walls 
		WHERE floor_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, floorID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve walls")
		return
	}
	defer rows.Close()

	for rows.Next() {
		wall := &entity.Wall{}
		err = rows.Scan(
			&wall.ID,
			&wall.X1, &wall.Y1,
			&wall.X2, &wall.Y2,
			&wall.WallTypeID,
			&wall.FloorID,
			&wall.CreatedAt, &wall.UpdatedAt, &wall.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan wall")
			return
		}
		walls = append(walls, wall)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(walls)
	if length == 0 {
		log.Info().Msgf("walls for floor ID %v were not found", floorID)
		return nil, ErrNotFound
	}

	log.Debug().Msgf("retrieved %d walls", length)
	return
}

func (r *wallRepo) Update(ctx context.Context, updateWallDTO *dto.PatchUpdateWallDTO) (err error) {
	query := "UPDATE walls SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateWallDTO.X1 != nil {
		updates = append(updates, fmt.Sprintf("x1 = $%d", paramID))
		params = append(params, updateWallDTO.X1)
		paramID++
	}
	if updateWallDTO.Y1 != nil {
		updates = append(updates, fmt.Sprintf("y1 = $%d", paramID))
		params = append(params, updateWallDTO.Y1)
		paramID++
	}
	if updateWallDTO.X2 != nil {
		updates = append(updates, fmt.Sprintf("x2 = $%d", paramID))
		params = append(params, updateWallDTO.X2)
		paramID++
	}
	if updateWallDTO.Y2 != nil {
		updates = append(updates, fmt.Sprintf("y2 = $%d", paramID))
		params = append(params, updateWallDTO.Y2)
		paramID++
	}
	if updateWallDTO.WallTypeID != nil {
		updates = append(updates, fmt.Sprintf("wall_type_id = $%d", paramID))
		params = append(params, updateWallDTO.WallTypeID)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateWallDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no wall found with the ID: %v", updateWallDTO.ID)
		return ErrNotFound
	}

	return
}

// Checks if the wall has been soft deleted
func (r *wallRepo) IsWallSoftDeleted(ctx context.Context, wallID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM walls WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, wallID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no wall found with UUID %v", wallID)
			return false, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve wall")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is wall deleted: %v", isDeleted)
	return
}

func (r *wallRepo) SoftDelete(ctx context.Context, wallID uuid.UUID) (err error) {
	query := `UPDATE walls SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, wallID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete wall")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no wall found with the UUID: %v", wallID)
		return ErrNotFound
	}

	log.Debug().Msg("wall deleted_at timestamp updated successfully")
	return
}

func (r *wallRepo) Restore(ctx context.Context, wallID uuid.UUID) (err error) {
	query := `UPDATE walls SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, wallID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore wall")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no wall found with the UUID: %v", wallID)
		return ErrNotFound
	}

	log.Debug().Msg("wall deleted_at timestamp set NULL successfully")
	return
}
