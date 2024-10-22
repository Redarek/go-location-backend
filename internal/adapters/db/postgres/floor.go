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
	"location-backend/internal/domain/service"
)

type floorRepo struct {
	pool *pgxpool.Pool
}

func NewFloorRepo(pool *pgxpool.Pool) *floorRepo {
	return &floorRepo{pool: pool}
}

func (r *floorRepo) Create(ctx context.Context, createFloorDTO *dto.CreateFloorDTO) (floorID uuid.UUID, err error) {
	// inserts := []string{}
	// params := []interface{}{}
	// numbers := []int{}
	// paramID := 4 // Количество параметров перед динамическим добавлением

	query := `INSERT INTO floors (
			name,
			number, 
			image,
			width_in_pixels,
			height_in_pixels,
			scale,
			building_id,
			cell_size_meter,
			north_area_indent_meter, south_area_indent_meter, west_area_indent_meter, east_area_indent_meter
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		&createFloorDTO.Name,
		&createFloorDTO.Number,
		&createFloorDTO.Image,
		&createFloorDTO.WestAreaIndentMeter,
		&createFloorDTO.HeightInPixels,
		&createFloorDTO.Scale,
		&createFloorDTO.BuildingID,
		&createFloorDTO.CellSizeMeter,
		&createFloorDTO.NorthAreaIndentMeter, &createFloorDTO.SouthAreaIndentMeter, &createFloorDTO.WestAreaIndentMeter, &createFloorDTO.EastAreaIndentMeter,
	)
	err = row.Scan(&floorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan floor")
		return
	}

	return
}

func (r *floorRepo) GetOne(ctx context.Context, floorID uuid.UUID) (floor *entity.Floor, err error) {
	query := `SELECT
			id,
			name,
			number,
			image,
			heatmap,
			width_in_pixels,
			height_in_pixels,
			scale,
			building_id,
			cell_size_meter,
			north_area_indent_meter, south_area_indent_meter, west_area_indent_meter, east_area_indent_meter,
			created_at, updated_at, deleted_at
		FROM floors WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, floorID)

	floor = &entity.Floor{}
	err = row.Scan(
		&floor.ID,
		&floor.Name,
		&floor.Number,
		&floor.Image,
		&floor.Heatmap,
		&floor.WidthInPixels, &floor.HeightInPixels,
		&floor.Scale,
		&floor.BuildingID,
		&floor.CellSizeMeter,
		&floor.NorthAreaIndentMeter, &floor.SouthAreaIndentMeter, &floor.WestAreaIndentMeter, &floor.EastAreaIndentMeter,
		&floor.CreatedAt, &floor.UpdatedAt, &floor.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no floor found with ID %v", floorID)
			return nil, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve floor")
		return
	}

	log.Debug().Msgf("retrieved floor: %v", floor)
	return
}

func (r *floorRepo) GetAll(ctx context.Context, buildingID uuid.UUID, limit, offset int) (floors []*entity.Floor, err error) {
	query := `SELECT
			id,
			name,
			number,
			image,
			heatmap,
			width_in_pixels,
			height_in_pixels,
			scale,
			building_id,
			cell_size_meter,
			north_area_indent_meter, south_area_indent_meter, west_area_indent_meter, east_area_indent_meter,
			created_at, updated_at, deleted_at
		FROM floors 
		WHERE building_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, buildingID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve floors")
		return
	}
	defer rows.Close()

	for rows.Next() {
		floor := &entity.Floor{}
		err = rows.Scan(
			&floor.ID,
			&floor.Name,
			&floor.Number,
			&floor.Image,
			&floor.Heatmap,
			&floor.WidthInPixels, &floor.HeightInPixels,
			&floor.Scale,
			&floor.BuildingID,
			&floor.CellSizeMeter,
			&floor.NorthAreaIndentMeter, &floor.SouthAreaIndentMeter, &floor.WestAreaIndentMeter, &floor.EastAreaIndentMeter,
			&floor.CreatedAt, &floor.UpdatedAt, &floor.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan floor")
			return
		}
		floors = append(floors, floor)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(floors)
	if length == 0 {
		log.Info().Msgf("floors for building ID %v were not found", buildingID)
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved %d floors", length)
	return
}

func (r *floorRepo) Update(ctx context.Context, patchUpdateFloorDTO *dto.PatchUpdateFloorDTO) (err error) {
	query := "UPDATE floors SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if patchUpdateFloorDTO.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, patchUpdateFloorDTO.Name)
		paramID++
	}
	if patchUpdateFloorDTO.Number != nil {
		updates = append(updates, fmt.Sprintf("number = $%d", paramID))
		params = append(params, patchUpdateFloorDTO.Number)
		paramID++
	}
	if patchUpdateFloorDTO.Image != nil {
		updates = append(updates, fmt.Sprintf("image = $%d", paramID))
		params = append(params, patchUpdateFloorDTO.Image)
		paramID++
	}
	// if patchUpdateFloorDTO.Heatmap != nil {
	// 	updates = append(updates, fmt.Sprintf("heatmap = $%d", paramID))
	// 	params = append(params, patchUpdateFloorDTO.Heatmap)
	// 	paramID++
	// }
	if patchUpdateFloorDTO.WidthInPixels != nil {
		updates = append(updates, fmt.Sprintf("width_in_pixels = $%d", paramID))
		params = append(params, patchUpdateFloorDTO.WidthInPixels)
		paramID++
	}
	if patchUpdateFloorDTO.HeightInPixels != nil {
		updates = append(updates, fmt.Sprintf("height_in_pixels = $%d", paramID))
		params = append(params, patchUpdateFloorDTO.HeightInPixels)
		paramID++
	}
	if patchUpdateFloorDTO.Scale != nil {
		updates = append(updates, fmt.Sprintf("scale = $%d", paramID))
		params = append(params, patchUpdateFloorDTO.Scale)
		paramID++
	}

	if patchUpdateFloorDTO.CellSizeMeter != nil {
		updates = append(updates, fmt.Sprintf("cell_size_meter = $%d", paramID))
		params = append(params, patchUpdateFloorDTO.CellSizeMeter)
		paramID++
	}
	if patchUpdateFloorDTO.NorthAreaIndentMeter != nil {
		updates = append(updates, fmt.Sprintf("north_area_indent_meter = $%d", paramID))
		params = append(params, patchUpdateFloorDTO.NorthAreaIndentMeter)
		paramID++
	}
	if patchUpdateFloorDTO.SouthAreaIndentMeter != nil {
		updates = append(updates, fmt.Sprintf("south_area_indent_meter = $%d", paramID))
		params = append(params, patchUpdateFloorDTO.SouthAreaIndentMeter)
		paramID++
	}
	if patchUpdateFloorDTO.WestAreaIndentMeter != nil {
		updates = append(updates, fmt.Sprintf("west_area_indent_meter = $%d", paramID))
		params = append(params, patchUpdateFloorDTO.WestAreaIndentMeter)
		paramID++
	}
	if patchUpdateFloorDTO.EastAreaIndentMeter != nil {
		updates = append(updates, fmt.Sprintf("east_area_indent_meter = $%d", paramID))
		params = append(params, patchUpdateFloorDTO.EastAreaIndentMeter)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return service.ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, patchUpdateFloorDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no floor found with the ID: %v", patchUpdateFloorDTO.ID)
		return service.ErrNotFound
	}

	return
}

func (r *floorRepo) UpdateHeatmap(ctx context.Context, floorID uuid.UUID, heatmap string) (err error) {
	query := `UPDATE floors SET 
		updated_at = NOW(), 
		heatmap = $1
	WHERE id = $2 AND deleted_at IS NULL`

	commandTag, err := r.pool.Exec(ctx, query, heatmap, floorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no floor found with the ID: %v", floorID)
		return service.ErrNotFound
	}

	return
}

// Checks if the floor has been soft deleted
func (r *floorRepo) IsFloorSoftDeleted(ctx context.Context, floorID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM floors WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, floorID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no floor found with ID %v", floorID)
			return false, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve floor")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is floor deleted: %v", isDeleted)
	return
}

func (r *floorRepo) SoftDelete(ctx context.Context, floorID uuid.UUID) (err error) {
	query := `UPDATE floors SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, floorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete floor")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no floor found with the UUID: %v", floorID)
		return service.ErrNotFound
	}

	log.Debug().Msg("floor deleted_at timestamp updated successfully")
	return
}

func (r *floorRepo) Restore(ctx context.Context, floorID uuid.UUID) (err error) {
	query := `UPDATE floors SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, floorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore floor")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no floor found with the UUID: %v", floorID)
		return service.ErrNotFound
	}

	log.Debug().Msg("floor deleted_at timestamp set NULL successfully")
	return
}
