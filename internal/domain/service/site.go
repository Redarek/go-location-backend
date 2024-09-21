package service

import (
	// "context"

	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	repository "location-backend/internal/adapters/db/postgres"
	// "location-backend/internal/controller/http/dto"
	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type SiteService interface {
	CreateSite(userCreate dto.CreateSiteDTO) (siteID uuid.UUID, err error)
	GetSite(siteID uuid.UUID) (site entity.Site, err error)
	// GetSiteByName(name string) (user entity.Site, err error)
}

type siteService struct {
	repository repository.SiteRepo
}

func NewSiteService(repository repository.SiteRepo) *siteService {
	return &siteService{repository: repository}
}

func (s *siteService) CreateSite(createSiteDTO dto.CreateSiteDTO) (siteID uuid.UUID, err error) {
	siteID, err = s.repository.Create(createSiteDTO)
	if err != nil {
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to create site")
		return
	}

	return siteID, nil
}

func (s *siteService) GetSite(siteID uuid.UUID) (site entity.Site, err error) {
	site, err = s.repository.GetOne(siteID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return site, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve site")
		return
	}

	return
}
