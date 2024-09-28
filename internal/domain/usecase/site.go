package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type SiteService interface {
	CreateSite(ctx context.Context, createSiteDTO *domain_dto.CreateSiteDTO) (siteID uuid.UUID, err error)
	GetSite(ctx context.Context, siteID uuid.UUID) (site *entity.Site, err error)
	GetSites(ctx context.Context, getSiteDTO domain_dto.GetSitesDTO) (sites []*entity.Site, err error)
	// TODO get site list detailed

	UpdateSite(ctx context.Context, patchUpdateSiteDTO *domain_dto.PatchUpdateSiteDTO) (err error)

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

func (u *SiteUsecase) CreateSite(ctx context.Context, dto *domain_dto.CreateSiteDTO) (siteID uuid.UUID, err error) {
	siteID, err = u.siteService.CreateSite(ctx, dto)
	if err != nil {
		log.Error().Err(err).Msg("failed to create site")
		return
	}

	log.Info().Msgf("site %v successfully created", dto.Name)
	return
}

func (u *SiteUsecase) GetSite(ctx context.Context, siteID uuid.UUID) (siteDTO *domain_dto.SiteDTO, err error) {
	site, err := u.siteService.GetSite(ctx, siteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get site")
			return
		}
	}

	// Mapping domain entity -> domain DTO
	siteDTO = &domain_dto.SiteDTO{
		ID:          site.ID,
		Name:        site.Name,
		Description: site.Description,
		UserID:      site.UserID,
		CreatedAt:   site.CreatedAt,
		UpdatedAt:   site.UpdatedAt,
		DeletedAt:   site.DeletedAt,
	}

	return
}

func (u *SiteUsecase) GetSites(ctx context.Context, dto domain_dto.GetSitesDTO) (sitesDTO []*domain_dto.SiteDTO, err error) {
	sites, err := u.siteService.GetSites(ctx, dto)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get sites")
			return
		}
	}

	for _, site := range sites {
		// Mapping domain entity -> domain DTO
		siteDTO := &domain_dto.SiteDTO{
			ID:          site.ID,
			Name:        site.Name,
			Description: site.Description,
			UserID:      site.UserID,
			CreatedAt:   site.CreatedAt,
			UpdatedAt:   site.UpdatedAt,
			DeletedAt:   site.DeletedAt,
		}

		sitesDTO = append(sitesDTO, siteDTO)
	}

	return
}

func (u *SiteUsecase) PatchUpdateSite(ctx context.Context, dto *domain_dto.PatchUpdateSiteDTO) (err error) {
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
