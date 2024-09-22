package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/service"
)

type SiteUsecase interface {
	CreateSite(dto domain_dto.CreateSiteDTO) (siteID uuid.UUID, err error)
	GetSite(dto domain_dto.GetSiteDTO) (siteDTO domain_dto.SiteDTO, err error)
	GetSites(ctx context.Context, dto domain_dto.GetSitesDTO) (sitesDTO []domain_dto.SiteDTO, err error)
	// TODO GetSitesDetailed

	PatchUpdateSite(ctx context.Context, patchUpdateDTO domain_dto.PatchUpdateSiteDTO) (err error)
	SoftDeleteSite(ctx context.Context, siteID uuid.UUID) (err error)
	RestoreSite(ctx context.Context, siteID uuid.UUID) (err error)
}

type siteUsecase struct {
	siteService service.SiteService
}

func NewSiteUsecase(siteService service.SiteService) *siteUsecase {
	return &siteUsecase{siteService: siteService}
}

func (u *siteUsecase) CreateSite(dto domain_dto.CreateSiteDTO) (siteID uuid.UUID, err error) {
	var createSiteDTO domain_dto.CreateSiteDTO = domain_dto.CreateSiteDTO{
		Name:        dto.Name,
		Description: dto.Description,
		UserID:      dto.UserID,
	}

	siteID, err = u.siteService.CreateSite(createSiteDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create site")
		return
	}

	log.Info().Msgf("site %v successfully created", dto.Name)
	return
}

func (u *siteUsecase) GetSite(dto domain_dto.GetSiteDTO) (siteDTO domain_dto.SiteDTO, err error) {
	site, err := u.siteService.GetSite(dto.ID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return domain_dto.SiteDTO{}, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get site")
			return
		}
	}

	// Mapping domain entity -> domain DTO
	siteDTO = domain_dto.SiteDTO{
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

func (u *siteUsecase) GetSites(ctx context.Context, dto domain_dto.GetSitesDTO) (sitesDTO []domain_dto.SiteDTO, err error) {
	sites, err := u.siteService.GetSites(ctx, dto)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return []domain_dto.SiteDTO{}, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get sites")
			return
		}
	}

	for _, site := range sites {
		// Mapping domain entity -> domain DTO
		siteDTO := domain_dto.SiteDTO{
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

func (u *siteUsecase) PatchUpdateSite(ctx context.Context, patchUpdateDTO domain_dto.PatchUpdateSiteDTO) (err error) {
	_, err = u.siteService.GetSite(patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			log.Error().Err(err).Msg("failed to check site existing")
			return ErrNotFound
		}
	}

	err = u.siteService.UpdateSite(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, service.ErrNotUpdated) {
			log.Info().Err(err).Msg("site was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update site")
		return
	}

	return
}

func (u *siteUsecase) SoftDeleteSite(ctx context.Context, siteID uuid.UUID) (err error) {
	isDeleted, err := u.siteService.IsSiteSoftDeleted(ctx, siteID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
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

func (u *siteUsecase) RestoreSite(ctx context.Context, siteID uuid.UUID) (err error) {
	isDeleted, err := u.siteService.IsSiteSoftDeleted(ctx, siteID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
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
