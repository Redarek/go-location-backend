package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
)

type ISiteRepo interface {
	Create(ctx context.Context, createSiteDTO *dto.CreateSiteDTO) (siteID uuid.UUID, err error)
	GetOne(ctx context.Context, siteID uuid.UUID) (site *entity.Site, err error)
	GetAll(ctx context.Context, userID uuid.UUID, limit, offset int) (sites []*entity.Site, err error)

	Update(ctx context.Context, patchUpdateSiteDTO *dto.PatchUpdateSiteDTO) (err error)

	IsSiteSoftDeleted(ctx context.Context, siteID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, siteID uuid.UUID) (err error)
	Restore(ctx context.Context, siteID uuid.UUID) (err error)
}

type siteService struct {
	siteRepo            ISiteRepo
	buildingRepo        IBuildingRepo
	wallTypeRepo        IWallTypeRepo
	accessPointTypeRepo IAccessPointTypeRepo
	sensorTypeRepo      ISensorTypeRepo
}

func NewSiteService(
	siteRepo ISiteRepo,
	buildingRepo IBuildingRepo,
	wallTypeRepo IWallTypeRepo,
	accessPointTypeRepo IAccessPointTypeRepo,
	sensorTypeRepo ISensorTypeRepo,
) *siteService {
	return &siteService{
		siteRepo:            siteRepo,
		buildingRepo:        buildingRepo,
		wallTypeRepo:        wallTypeRepo,
		accessPointTypeRepo: accessPointTypeRepo,
		sensorTypeRepo:      sensorTypeRepo,
	}
}

func (s *siteService) CreateSite(ctx context.Context, createSiteDTO *dto.CreateSiteDTO) (siteID uuid.UUID, err error) {
	siteID, err = s.siteRepo.Create(ctx, createSiteDTO)
	if err != nil {
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to create site")
		return
	}

	return siteID, nil
}

func (s *siteService) GetSite(ctx context.Context, siteID uuid.UUID) (site *entity.Site, err error) {
	site, err = s.siteRepo.GetOne(ctx, siteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return site, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *siteService) GetSiteDetailed(ctx context.Context, getDTO dto.GetSiteDetailedDTO) (siteDetailed *entity.SiteDetailed, err error) {
	site, err := s.GetSite(ctx, getDTO.ID)
	if err != nil {
		return
	}

	buildings, err := s.buildingRepo.GetAll(ctx, getDTO.ID, getDTO.Limit, getDTO.Offset)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return
		}
	}

	wallTypes, err := s.wallTypeRepo.GetAll(ctx, getDTO.ID, getDTO.Limit, getDTO.Offset)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return
		}
	}

	accessPointTypes, err := s.accessPointTypeRepo.GetAll(ctx, getDTO.ID, getDTO.Limit, getDTO.Offset)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return
		}
	}

	sensorTypes, err := s.sensorTypeRepo.GetAll(ctx, getDTO.ID, getDTO.Limit, getDTO.Offset)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return
		}
	}

	siteDetailed = &entity.SiteDetailed{
		Site:             *site,
		Buildings:        buildings,
		WallTypes:        wallTypes,
		AccessPointTypes: accessPointTypes,
		SensorTypes:      sensorTypes,
	}

	return
}

func (s *siteService) GetSites(ctx context.Context, dto dto.GetSitesDTO) (sites []*entity.Site, err error) {
	sites, err = s.siteRepo.GetAll(ctx, dto.UserID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return sites, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *siteService) GetSitesDetailed(ctx context.Context, getDTO dto.GetSitesDTO) (sitesDetailed []*entity.SiteDetailed, err error) {
	sites, err := s.GetSites(ctx, getDTO)
	if err != nil {
		return
	}

	for _, site := range sites {
		buildings, err := s.buildingRepo.GetAll(ctx, site.ID, getDTO.Limit, getDTO.Offset)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				return nil, err
			}
		}

		wallTypes, err := s.wallTypeRepo.GetAll(ctx, site.ID, getDTO.Limit, getDTO.Offset)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				return nil, err
			}
		}

		accessPointTypes, err := s.accessPointTypeRepo.GetAll(ctx, site.ID, getDTO.Limit, getDTO.Offset)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				return nil, err
			}
		}

		sensorTypes, err := s.sensorTypeRepo.GetAll(ctx, site.ID, getDTO.Limit, getDTO.Offset)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				return nil, err
			}
		}

		siteDetailed := &entity.SiteDetailed{
			Site:             *site,
			Buildings:        buildings,
			WallTypes:        wallTypes,
			AccessPointTypes: accessPointTypes,
			SensorTypes:      sensorTypes,
		}

		sitesDetailed = append(sitesDetailed, siteDetailed)
	}

	return
}

// TODO PUT update
func (s *siteService) UpdateSite(ctx context.Context, patchUpdateSiteDTO *dto.PatchUpdateSiteDTO) (err error) {
	err = s.siteRepo.Update(ctx, patchUpdateSiteDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
		if errors.Is(err, ErrNotUpdated) {
			return usecase.ErrNotUpdated
		}

		return
	}

	return
}

func (s *siteService) IsSiteSoftDeleted(ctx context.Context, siteID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.siteRepo.IsSiteSoftDeleted(ctx, siteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, usecase.ErrNotFound
		}
		return
	}

	return
}

func (s *siteService) SoftDeleteSite(ctx context.Context, siteID uuid.UUID) (err error) {
	err = s.siteRepo.SoftDelete(ctx, siteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
		return
	}

	return
}

func (s *siteService) RestoreSite(ctx context.Context, siteID uuid.UUID) (err error) {
	err = s.siteRepo.Restore(ctx, siteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
		return
	}

	return
}
