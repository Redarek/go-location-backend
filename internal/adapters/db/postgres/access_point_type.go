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

type AccessPointTypeRepo interface {
	Create(ctx context.Context, createAccessPointTypeDTO *dto.CreateAccessPointTypeDTO) (accessPointTypeID uuid.UUID, err error)
	GetOne(ctx context.Context, accessPointTypeID uuid.UUID) (accessPointType *entity.AccessPointType, err error)
	// GetOneDetailed(ctx context.Context, accessPointTypeID uuid.UUID) (accessPointType *entity.AccessPointTypeDetailed, err error) // TODO

	GetAll(ctx context.Context, siteID uuid.UUID, limit, offset int) (accessPointTypes []*entity.AccessPointType, err error)

	Update(ctx context.Context, updateAccessPointTypeDTO *dto.PatchUpdateAccessPointTypeDTO) (err error)

	IsAccessPointTypeSoftDeleted(ctx context.Context, accessPointTypeID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, accessPointTypeID uuid.UUID) (err error)
	Restore(ctx context.Context, accessPointTypeID uuid.UUID) (err error)
}

type accessPointTypeRepo struct {
	pool *pgxpool.Pool
}

func NewAccessPointTypeRepo(pool *pgxpool.Pool) *accessPointTypeRepo {
	return &accessPointTypeRepo{pool: pool}
}

func (r *accessPointTypeRepo) Create(ctx context.Context, createAccessPointTypeDTO *dto.CreateAccessPointTypeDTO) (accessPointTypeID uuid.UUID, err error) {
	query := `INSERT INTO access_point_types (
			name, 
			model,
			color, 
			z,
			is_virtual,
			site_id
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		createAccessPointTypeDTO.Name,
		createAccessPointTypeDTO.Model,
		createAccessPointTypeDTO.Color,
		createAccessPointTypeDTO.Z,
		createAccessPointTypeDTO.IsVirtual,
		createAccessPointTypeDTO.SiteID,
	)
	err = row.Scan(&accessPointTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan accessPointType")
		return
	}

	return
}

func (r *accessPointTypeRepo) GetOne(ctx context.Context, accessPointTypeID uuid.UUID) (accessPointType *entity.AccessPointType, err error) {
	query := `SELECT 
			id, 
			name,
			model,
			color,
			z,
			is_virtual,
			site_id,
			created_at, updated_at, deleted_at
		FROM access_point_types WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, accessPointTypeID)

	accessPointType = &entity.AccessPointType{}
	err = row.Scan(
		&accessPointType.ID,
		&accessPointType.Name,
		&accessPointType.Model,
		&accessPointType.Color,
		&accessPointType.Z,
		&accessPointType.IsVirtual,
		&accessPointType.SiteID,
		&accessPointType.CreatedAt, &accessPointType.UpdatedAt, &accessPointType.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no accessPointType found with ID %v", accessPointTypeID)
			return nil, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve accessPointType")
		return
	}
	log.Debug().Msgf("retrieved accessPointType: %v", accessPointType)
	return
}

// func (r *accessPointTypeRepo) GetOneDetailed(ctx context.Context, accessPointTypeID uuid.UUID) (accessPointType *entity.AccessPointTypeDetailed, err error) {

// }

func (r *accessPointTypeRepo) GetAll(ctx context.Context, siteID uuid.UUID, limit, offset int) (accessPointTypes []*entity.AccessPointType, err error) {
	query := `SELECT 
			id, 
			name,
			model,
			color,
			z,
			is_virtual,
			site_id,
			created_at, updated_at, deleted_at
		FROM access_point_types 
		WHERE site_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, siteID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve accessPointTypes")
		return
	}
	defer rows.Close()

	for rows.Next() {
		accessPointType := &entity.AccessPointType{}
		err = rows.Scan(
			&accessPointType.ID,
			&accessPointType.Name,
			&accessPointType.Model,
			&accessPointType.Color,
			&accessPointType.Z,
			&accessPointType.IsVirtual,
			&accessPointType.SiteID,
			&accessPointType.CreatedAt, &accessPointType.UpdatedAt, &accessPointType.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan accessPointType")
			return
		}
		accessPointTypes = append(accessPointTypes, accessPointType)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(accessPointTypes)
	if length == 0 {
		log.Info().Msgf("accessPointTypes for site ID %v were not found", siteID)
		return nil, ErrNotFound
	}

	log.Debug().Msgf("retrieved %d accessPointTypes", length)
	return
}

func (r *accessPointTypeRepo) Update(ctx context.Context, updateAccessPointTypeDTO *dto.PatchUpdateAccessPointTypeDTO) (err error) {
	query := "UPDATE access_point_types SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateAccessPointTypeDTO.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, updateAccessPointTypeDTO.Name)
		paramID++
	}
	if updateAccessPointTypeDTO.Model != nil {
		updates = append(updates, fmt.Sprintf("model = $%d", paramID))
		params = append(params, updateAccessPointTypeDTO.Model)
		paramID++
	}

	if updateAccessPointTypeDTO.Color != nil {
		updates = append(updates, fmt.Sprintf("color = $%d", paramID))
		params = append(params, updateAccessPointTypeDTO.Color)
		paramID++
	}
	if updateAccessPointTypeDTO.Z != nil {
		updates = append(updates, fmt.Sprintf("z = $%d", paramID))
		params = append(params, updateAccessPointTypeDTO.Z)
		paramID++
	}
	if updateAccessPointTypeDTO.IsVirtual != nil {
		updates = append(updates, fmt.Sprintf("is_virtual = $%d", paramID))
		params = append(params, updateAccessPointTypeDTO.IsVirtual)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateAccessPointTypeDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no accessPointType found with the UUID: %v", updateAccessPointTypeDTO.ID)
		return ErrNotFound
	}

	return
}

// Checks if the accessPointType has been soft deleted
func (r *accessPointTypeRepo) IsAccessPointTypeSoftDeleted(ctx context.Context, accessPointTypeID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM access_point_types WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, accessPointTypeID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no accessPointType found with UUID %v", accessPointTypeID)
			return false, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve accessPointType")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is accessPointType deleted: %v", isDeleted)
	return
}

func (r *accessPointTypeRepo) SoftDelete(ctx context.Context, accessPointTypeID uuid.UUID) (err error) {
	query := `UPDATE access_point_types SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, accessPointTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete accessPointType")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no accessPointType found with the UUID: %v", accessPointTypeID)
		return ErrNotFound
	}

	log.Debug().Msg("accessPointType deleted_at timestamp updated successfully")
	return
}

func (r *accessPointTypeRepo) Restore(ctx context.Context, accessPointTypeID uuid.UUID) (err error) {
	query := `UPDATE access_point_types SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, accessPointTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore accessPointType")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no accessPointType found with the UUID: %v", accessPointTypeID)
		return ErrNotFound
	}

	log.Debug().Msg("accessPointType deleted_at timestamp set NULL successfully")
	return
}
