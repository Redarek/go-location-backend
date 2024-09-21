package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type SiteRepo interface {
	Create(createSiteDTO dto.CreateSiteDTO) (siteID uuid.UUID, err error)
	GetOne(siteID uuid.UUID) (site entity.Site, err error)
}

type siteRepo struct {
	pool *pgxpool.Pool
}

func NewSiteRepo(pool *pgxpool.Pool) *siteRepo {
	return &siteRepo{pool: pool}
}

func (r *siteRepo) Create(createSiteDTO dto.CreateSiteDTO) (siteID uuid.UUID, err error) {
	query := `INSERT INTO sites (
			name, 
			description, 
			user_id
		)
		VALUES ($1, $2, $3)
		RETURNING id`
	row := r.pool.QueryRow(context.Background(), query,
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

func (r *siteRepo) GetOne(siteID uuid.UUID) (site entity.Site, err error) {
	query := `SELECT 
			id, 
			name,
			description,
			user_id,
			created_at, 
			updated_at, 
			deleted_at 
		FROM sites
		WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(context.Background(), query, siteID)
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
			return site, ErrNotFound
		}
		log.Error().Err(err).Msg("failed to scan site")
		return
	}
	log.Debug().Msgf("retrieved site: %v", site)
	return
}
