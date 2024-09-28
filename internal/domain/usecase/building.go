package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type BuildingService interface {
	CreateBuilding(ctx context.Context, createBuildingDTO *domain_dto.CreateBuildingDTO) (buildingID uuid.UUID, err error)
	GetBuilding(ctx context.Context, buildingID uuid.UUID) (building *entity.Building, err error)
	GetBuildings(ctx context.Context, dto domain_dto.GetBuildingsDTO) (buildings []*entity.Building, err error)
	// TODO get building list detailed

	UpdateBuilding(ctx context.Context, updateBuildingDTO *domain_dto.PatchUpdateBuildingDTO) (err error)

	IsBuildingSoftDeleted(ctx context.Context, buildingID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteBuilding(ctx context.Context, buildingID uuid.UUID) (err error)
	RestoreBuilding(ctx context.Context, buildingID uuid.UUID) (err error)
}

type BuildingUsecase interface {
	CreateBuilding(ctx context.Context, dto *domain_dto.CreateBuildingDTO) (buildingID uuid.UUID, err error)
	GetBuilding(ctx context.Context, buildingID uuid.UUID) (buildingDTO *domain_dto.BuildingDTO, err error)
	GetBuildings(ctx context.Context, dto domain_dto.GetBuildingsDTO) (buildingsDTO []*domain_dto.BuildingDTO, err error)
	// TODO GetBuildingsDetailed

	PatchUpdateBuilding(ctx context.Context, patchUpdateDTO *domain_dto.PatchUpdateBuildingDTO) (err error)

	SoftDeleteBuilding(ctx context.Context, buildingID uuid.UUID) (err error)
	RestoreBuilding(ctx context.Context, buildingID uuid.UUID) (err error)
}

type buildingUsecase struct {
	buildingService BuildingService
	siteService     SiteService
}

func NewBuildingUsecase(buildingService BuildingService, siteService SiteService) *buildingUsecase {
	return &buildingUsecase{
		buildingService: buildingService,
		siteService:     siteService,
	}
}

func (u *buildingUsecase) CreateBuilding(ctx context.Context, dto *domain_dto.CreateBuildingDTO) (buildingID uuid.UUID, err error) {
	_, err = u.siteService.GetSite(ctx, dto.SiteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Err(err).Msg("failed to create building: the site with provided site ID does not exist")
			return buildingID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check site existing")
		return
	}

	buildingID, err = u.buildingService.CreateBuilding(ctx, dto)
	if err != nil {
		log.Error().Err(err).Msg("failed to create building")
		return
	}

	log.Info().Msgf("building %v successfully created", buildingID)
	return
}

func (u *buildingUsecase) GetBuilding(ctx context.Context, buildingID uuid.UUID) (buildingDTO *domain_dto.BuildingDTO, err error) {
	building, err := u.buildingService.GetBuilding(ctx, buildingID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get building")
			return
		}
	}

	// Mapping domain entity -> domain DTO
	buildingDTO = &domain_dto.BuildingDTO{
		ID:          building.ID,
		Name:        building.Name,
		Description: building.Description,
		Country:     building.Country,
		City:        building.City,
		Address:     building.Address,
		SiteID:      building.SiteID,
		CreatedAt:   building.CreatedAt,
		UpdatedAt:   building.UpdatedAt,
		DeletedAt:   building.DeletedAt,
	}

	return
}

func (u *buildingUsecase) GetBuildings(ctx context.Context, dto domain_dto.GetBuildingsDTO) (buildingsDTO []*domain_dto.BuildingDTO, err error) {
	buildings, err := u.buildingService.GetBuildings(ctx, dto)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get buildings")
			return
		}
	}

	for _, building := range buildings {
		// Mapping domain entity -> domain DTO
		buildingDTO := &domain_dto.BuildingDTO{
			ID:          building.ID,
			Name:        building.Name,
			Description: building.Description,
			Country:     building.Country,
			City:        building.City,
			Address:     building.Address,
			SiteID:      building.SiteID,
			CreatedAt:   building.CreatedAt,
			UpdatedAt:   building.UpdatedAt,
			DeletedAt:   building.DeletedAt,
		}

		buildingsDTO = append(buildingsDTO, buildingDTO)
	}

	return
}

func (u *buildingUsecase) PatchUpdateBuilding(ctx context.Context, patchUpdateDTO *domain_dto.PatchUpdateBuildingDTO) (err error) {
	_, err = u.buildingService.GetBuilding(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Error().Err(err).Msg("failed to check building existing")
			return ErrNotFound
		}
	}

	err = u.buildingService.UpdateBuilding(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Err(err).Msg("building was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update building")
		return
	}

	return
}

func (u *buildingUsecase) SoftDeleteBuilding(ctx context.Context, buildingID uuid.UUID) (err error) {
	isDeleted, err := u.buildingService.IsBuildingSoftDeleted(ctx, buildingID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if building is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.buildingService.SoftDeleteBuilding(ctx, buildingID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete building")
		return
	}

	return
}

func (u *buildingUsecase) RestoreBuilding(ctx context.Context, buildingID uuid.UUID) (err error) {
	isDeleted, err := u.buildingService.IsBuildingSoftDeleted(ctx, buildingID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if building is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.buildingService.RestoreBuilding(ctx, buildingID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore building")
		return
	}

	return
}
