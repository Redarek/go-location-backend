package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type SiteService interface {
	CreateSite(ctx context.Context, createSiteDTO *dto.CreateSiteDTO) (siteID uuid.UUID, err error)
	GetSite(ctx context.Context, siteID uuid.UUID) (site *entity.Site, err error)
	GetSiteDetailed(ctx context.Context, getDTO dto.GetSiteDetailedDTO) (siteDetailed *entity.SiteDetailed, err error)
	GetSites(ctx context.Context, getSiteDTO dto.GetSitesDTO) (sites []*entity.Site, err error)
	GetSitesDetailed(ctx context.Context, getDTO dto.GetSitesDTO) (sitesDetailed []*entity.SiteDetailed, err error)

	UpdateSite(ctx context.Context, patchUpdateSiteDTO *dto.PatchUpdateSiteDTO) (err error)

	IsSiteSoftDeleted(ctx context.Context, siteID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteSite(ctx context.Context, siteID uuid.UUID) (err error)
	RestoreSite(ctx context.Context, siteID uuid.UUID) (err error)
}

type SiteUsecase struct {
	siteService SiteService
}

func NewSiteUsecase(siteService SiteService) *SiteUsecase {
	return &SiteUsecase{siteService: siteService}
}

func (u *SiteUsecase) CreateSite(ctx context.Context, dto *dto.CreateSiteDTO) (siteID uuid.UUID, err error) {
	siteID, err = u.siteService.CreateSite(ctx, dto)
	if err != nil {
		log.Error().Err(err).Msg("failed to create site")
		return
	}

	log.Info().Msgf("site %v successfully created", dto.Name)
	return
}

func (u *SiteUsecase) GetSite(ctx context.Context, siteID uuid.UUID) (site *entity.Site, err error) {
	site, err = u.siteService.GetSite(ctx, siteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get site")
			return
		}
	}

	return
}

func (u *SiteUsecase) GetSiteDetailed(ctx context.Context, getDTO dto.GetSiteDetailedDTO) (siteDetailed *entity.SiteDetailed, err error) {
	siteDetailed, err = u.siteService.GetSiteDetailed(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get site detailed")
			return
		}
	}

	return
}

func (u *SiteUsecase) GetSites(ctx context.Context, dto dto.GetSitesDTO) (sites []*entity.Site, err error) {
	sites, err = u.siteService.GetSites(ctx, dto)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get sites")
			return
		}
	}

	return
}

func (u *SiteUsecase) GetSitesDetailed(ctx context.Context, getDTO dto.GetSitesDTO) (sitesDetailed []*entity.SiteDetailed, err error) {
	sitesDetailed, err = u.siteService.GetSitesDetailed(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get sites detailed")
			return
		}
	}

	return
}

func (u *SiteUsecase) PatchUpdateSite(ctx context.Context, dto *dto.PatchUpdateSiteDTO) (err error) {
	_, err = u.siteService.GetSite(ctx, dto.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Error().Err(err).Msg("failed to check site existing")
			return ErrNotFound
		}
	}

	err = u.siteService.UpdateSite(ctx, dto)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Err(err).Msg("site was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update site")
		return
	}

	return
}

func (u *SiteUsecase) SoftDeleteSite(ctx context.Context, siteID uuid.UUID) (err error) {
	isDeleted, err := u.siteService.IsSiteSoftDeleted(ctx, siteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if site is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.siteService.SoftDeleteSite(ctx, siteID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete site")
		return
	}

	return
}

func (u *SiteUsecase) RestoreSite(ctx context.Context, siteID uuid.UUID) (err error) {
	isDeleted, err := u.siteService.IsSiteSoftDeleted(ctx, siteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if site is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.siteService.RestoreSite(ctx, siteID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore site")
		return
	}

	return
}
