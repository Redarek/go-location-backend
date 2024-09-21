package usecase

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	domain_dto "location-backend/internal/domain/dto"
	// "location-backend/internal/domain/entity"
	"location-backend/internal/domain/service"
)

type SiteUsecase interface {
	CreateSite(dto domain_dto.CreateSiteDTO) (siteID uuid.UUID, err error)
	// GetSite(dto domain_dto.GetSiteDTO) (site entity.Site, err error)
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
