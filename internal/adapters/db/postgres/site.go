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

type SiteRepo interface {
	Create(ctx context.Context, createSiteDTO dto.CreateSiteDTO) (siteID uuid.UUID, err error)
	GetOne(ctx context.Context, siteID uuid.UUID) (site *entity.Site, err error)
	GetAll(ctx context.Context, userID uuid.UUID, limit, offset int) (sites []*entity.Site, err error)

	Update(ctx context.Context, updateSiteDTO dto.PatchUpdateSiteDTO) (err error)

	IsSiteSoftDeleted(ctx context.Context, siteID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, siteID uuid.UUID) (err error)
	Restore(ctx context.Context, siteID uuid.UUID) (err error)
}

type siteRepo struct {
	pool *pgxpool.Pool
}

func NewSiteRepo(pool *pgxpool.Pool) *siteRepo {
	return &siteRepo{pool: pool}
}

func (r *siteRepo) Create(ctx context.Context, createSiteDTO dto.CreateSiteDTO) (siteID uuid.UUID, err error) {
	query := `INSERT INTO sites (
			name, 
			description, 
			user_id
		)
		VALUES ($1, $2, $3)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		createSiteDTO.Name,
		createSiteDTO.Description,
		createSiteDTO.UserID,
	)
	var site entity.Site
	err = row.Scan(&site.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan site")
		return uuid.UUID{}, err
	}

	return site.ID, nil
}

func (r *siteRepo) GetOne(ctx context.Context, siteID uuid.UUID) (site *entity.Site, err error) {
	query := `SELECT 
			id, 
			name,
			description,
			user_id,
			created_at, updated_at, deleted_at 
		FROM sites
		WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, siteID)

	site = &entity.Site{}
	err = row.Scan(
		&site.ID,
		&site.Name,
		&site.Description,
		&site.UserID,
		&site.CreatedAt,
		&site.UpdatedAt,
		&site.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Msgf("site with ID %v not found", siteID)
			return nil, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to scan site")
		return
	}
	log.Debug().Msgf("retrieved site: %v", site)
	return
}

func (r *siteRepo) GetAll(ctx context.Context, userID uuid.UUID, limit, offset int) (sites []*entity.Site, err error) {
	query := `SELECT 
			id,
			name, 
			description, 
			user_id, 
			created_at, updated_at, deleted_at 
		FROM sites 
		WHERE user_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve sites")
		return
	}
	defer rows.Close()

	for rows.Next() {
		site := new(entity.Site)
		err = rows.Scan(
			&site.ID,
			&site.Name,
			&site.Description,
			&site.UserID,
			&site.CreatedAt, &site.UpdatedAt, &site.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan site")
			return
		}
		sites = append(sites, site)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(sites)
	if length == 0 {
		log.Info().Msgf("sites for user ID %v were not found", userID)
		return nil, ErrNotFound
	}

	log.Debug().Msgf("retrieved %d sites", length)
	return
}

func (r *siteRepo) Update(ctx context.Context, updateSiteDTO dto.PatchUpdateSiteDTO) (err error) {
	query := "UPDATE sites SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateSiteDTO.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, updateSiteDTO.Name)
		paramID++
	}
	if updateSiteDTO.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", paramID))
		params = append(params, updateSiteDTO.Description)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateSiteDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no site found with the UUID: %v", updateSiteDTO.ID)
		return ErrNotFound
	}

	return
}

// Checks if the site has been soft deleted
func (r *siteRepo) IsSiteSoftDeleted(ctx context.Context, siteID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM sites WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, siteID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no site found with UUID %v", siteID)
			return false, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve site")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is site deleted: %v", isDeleted)
	return
}

func (r *siteRepo) SoftDelete(ctx context.Context, siteID uuid.UUID) (err error) {
	query := `UPDATE sites SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, siteID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete site")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no site found with the UUID: %v", siteID)
		return ErrNotFound
	}

	log.Debug().Msg("site deleted_at timestamp updated successfully")
	return
}

func (r *siteRepo) Restore(ctx context.Context, siteID uuid.UUID) (err error) {
	query := `UPDATE sites SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, siteID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore site")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no site found with the UUID: %v", siteID)
		return ErrNotFound
	}

	log.Debug().Msg("site deleted_at timestamp set NULL successfully")
	return
}
