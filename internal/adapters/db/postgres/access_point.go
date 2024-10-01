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

type accessPointRepo struct {
	pool *pgxpool.Pool
}

func NewAccessPointRepo(pool *pgxpool.Pool) *accessPointRepo {
	return &accessPointRepo{pool: pool}
}

func (r *accessPointRepo) Create(ctx context.Context, createAccessPointDTO *dto.CreateAccessPointDTO) (accessPointID uuid.UUID, err error) {
	query := `INSERT INTO access_points (
			name, 
			x, 
			y,
			z,
			is_virtual,
			access_point_type_id,
			floor_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		createAccessPointDTO.Name,
		createAccessPointDTO.X,
		createAccessPointDTO.Y,
		createAccessPointDTO.Z,
		createAccessPointDTO.IsVirtual,
		createAccessPointDTO.AccessPointTypeID,
		createAccessPointDTO.FloorID,
	)
	err = row.Scan(&accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan access point")
		return
	}

	return
}

func (r *accessPointRepo) GetOne(ctx context.Context, accessPointID uuid.UUID) (accessPoint *entity.AccessPoint, err error) {
	query := `SELECT 
			id, 
			name, 
			x, 
			y,
			z,
			is_virtual,
			access_point_type_id,
			floor_id,
			created_at, updated_at, deleted_at
		FROM access_points WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, accessPointID)

	accessPoint = &entity.AccessPoint{}
	err = row.Scan(
		&accessPoint.ID,
		&accessPoint.Name,
		&accessPoint.X,
		&accessPoint.Y,
		&accessPoint.Z,
		&accessPoint.IsVirtual,
		&accessPoint.AccessPointTypeID,
		&accessPoint.FloorID,
		&accessPoint.CreatedAt, &accessPoint.UpdatedAt, &accessPoint.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no access point found with ID %v", accessPointID)
			return nil, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve access point")
		return
	}
	log.Debug().Msgf("retrieved access point: %v", accessPoint)
	return
}

// func (r *accessPointRepo) GetOneDetailed(ctx context.Context, accessPointID uuid.UUID) (accessPoint *entity.AccessPointDetailed, err error) {

// }

func (r *accessPointRepo) GetAll(ctx context.Context, floorID uuid.UUID, limit, offset int) (accessPoints []*entity.AccessPoint, err error) {
	query := `SELECT 
			id, 
			name, 
			x, 
			y,
			z,
			is_virtual,
			access_point_type_id,
			floor_id,
			created_at, updated_at, deleted_at
		FROM access_points 
		WHERE floor_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, floorID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve accessPoints")
		return
	}
	defer rows.Close()

	for rows.Next() {
		accessPoint := &entity.AccessPoint{}
		err = rows.Scan(
			&accessPoint.ID,
			&accessPoint.Name,
			&accessPoint.X,
			&accessPoint.Y,
			&accessPoint.Z,
			&accessPoint.IsVirtual,
			&accessPoint.AccessPointTypeID,
			&accessPoint.FloorID,
			&accessPoint.CreatedAt, &accessPoint.UpdatedAt, &accessPoint.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan access point")
			return
		}
		accessPoints = append(accessPoints, accessPoint)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(accessPoints)
	if length == 0 {
		log.Info().Msgf("access points for floor ID %v were not found", floorID)
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved %d access points", length)
	return
}

func (r *accessPointRepo) Update(ctx context.Context, updateAccessPointDTO *dto.PatchUpdateAccessPointDTO) (err error) {
	query := "UPDATE access_points SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateAccessPointDTO.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, updateAccessPointDTO.Name)
		paramID++
	}
	if updateAccessPointDTO.X != nil {
		updates = append(updates, fmt.Sprintf("x = $%d", paramID))
		params = append(params, updateAccessPointDTO.X)
		paramID++
	}
	if updateAccessPointDTO.Y != nil {
		updates = append(updates, fmt.Sprintf("y = $%d", paramID))
		params = append(params, updateAccessPointDTO.Y)
		paramID++
	}
	if updateAccessPointDTO.Z != nil {
		updates = append(updates, fmt.Sprintf("z = $%d", paramID))
		params = append(params, updateAccessPointDTO.Z)
		paramID++
	}
	if updateAccessPointDTO.IsVirtual != nil {
		updates = append(updates, fmt.Sprintf("is_virtual = $%d", paramID))
		params = append(params, updateAccessPointDTO.IsVirtual)
		paramID++
	}
	if updateAccessPointDTO.AccessPointTypeID != nil {
		updates = append(updates, fmt.Sprintf("access_point_type_id = $%d", paramID))
		params = append(params, updateAccessPointDTO.AccessPointTypeID)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return service.ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateAccessPointDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no access point found with the ID: %v", updateAccessPointDTO.ID)
		return service.ErrNotFound
	}

	return
}

// Checks if the accessPoint has been soft deleted
func (r *accessPointRepo) IsAccessPointSoftDeleted(ctx context.Context, accessPointID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM access_points WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, accessPointID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no access point found with UUID %v", accessPointID)
			return false, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve access point")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is accessPoint deleted: %v", isDeleted)
	return
}

func (r *accessPointRepo) SoftDelete(ctx context.Context, accessPointID uuid.UUID) (err error) {
	query := `UPDATE access_points SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete access point")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no access point found with the ID: %v", accessPointID)
		return service.ErrNotFound
	}

	log.Debug().Msg("access point deleted_at timestamp updated successfully")
	return
}

func (r *accessPointRepo) Restore(ctx context.Context, accessPointID uuid.UUID) (err error) {
	query := `UPDATE access_points SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore access point")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no access point found with the UUID: %v", accessPointID)
		return service.ErrNotFound
	}

	log.Debug().Msg("access point deleted_at timestamp set NULL successfully")
	return
}
