package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type SiteRepo interface {
	Create(ctx context.Context, createSiteDTO *dto.CreateSiteDTO) (siteID uuid.UUID, err error)
	GetOne(ctx context.Context, siteID uuid.UUID) (site *entity.Site, err error)
	GetAll(ctx context.Context, userID uuid.UUID, limit, offset int) (sites []*entity.Site, err error)

	Update(ctx context.Context, patchUpdateSiteDTO *dto.PatchUpdateSiteDTO) (err error)

	IsSiteSoftDeleted(ctx context.Context, siteID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, siteID uuid.UUID) (err error)
	Restore(ctx context.Context, siteID uuid.UUID) (err error)
}

type SiteService interface {
	CreateSite(ctx context.Context, createSiteDTO *dto.CreateSiteDTO) (siteID uuid.UUID, err error)
	GetSite(ctx context.Context, siteID uuid.UUID) (site *entity.Site, err error)
	GetSites(ctx context.Context, getSiteDTO dto.GetSitesDTO) (sites []*entity.Site, err error)
	// TODO get site list detailed

	UpdateSite(ctx context.Context, patchUpdateSiteDTO *dto.PatchUpdateSiteDTO) (err error)

	IsSiteSoftDeleted(ctx context.Context, siteID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteSite(ctx context.Context, siteID uuid.UUID) (err error)
	RestoreSite(ctx context.Context, siteID uuid.UUID) (err error)
}

type siteService struct {
	repository SiteRepo
}

func NewSiteService(repository SiteRepo) *siteService {
	return &siteService{repository: repository}
}

func (s *siteService) CreateSite(ctx context.Context, createSiteDTO *dto.CreateSiteDTO) (siteID uuid.UUID, err error) {
	siteID, err = s.repository.Create(ctx, createSiteDTO)
	if err != nil {
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to create site")
		return
	}

	return siteID, nil
}

func (s *siteService) GetSite(ctx context.Context, siteID uuid.UUID) (site *entity.Site, err error) {
	site, err = s.repository.GetOne(ctx, siteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return site, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve site")
		return
	}

	return
}

func (s *siteService) GetSites(ctx context.Context, dto dto.GetSitesDTO) (sites []*entity.Site, err error) {
	sites, err = s.repository.GetAll(ctx, dto.UserID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return sites, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve site")
		return
	}

	return
}

// TODO PUT update
func (s *siteService) UpdateSite(ctx context.Context, patchUpdateSiteDTO *dto.PatchUpdateSiteDTO) (err error) {
	err = s.repository.Update(ctx, patchUpdateSiteDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		if errors.Is(err, ErrNotUpdated) {
			return ErrNotUpdated
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to update site")
		return
	}

	return
}

func (s *siteService) IsSiteSoftDeleted(ctx context.Context, siteID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.repository.IsSiteSoftDeleted(ctx, siteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve site")
		return
	}

	return
}

func (s *siteService) SoftDeleteSite(ctx context.Context, siteID uuid.UUID) (err error) {
	err = s.repository.SoftDelete(ctx, siteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to soft delete site")
		return
	}

	return
}

func (s *siteService) RestoreSite(ctx context.Context, siteID uuid.UUID) (err error) {
	err = s.repository.Restore(ctx, siteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to restore site")
		return
	}

	return
}
