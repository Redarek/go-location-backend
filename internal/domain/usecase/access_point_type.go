package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type IAccessPointTypeService interface {
	CreateAccessPointType(ctx context.Context, createDTO *dto.CreateAccessPointTypeDTO) (accessPointTypeID uuid.UUID, err error)
	GetAccessPointType(ctx context.Context, accessPointTypeID uuid.UUID) (accessPointType *entity.AccessPointType, err error)
	GetAccessPointTypeDetailed(ctx context.Context, getDTO dto.GetAccessPointTypeDetailedDTO) (accessPointTypeDetailed *entity.AccessPointTypeDetailed, err error)
	GetAccessPointTypes(ctx context.Context, getDTO dto.GetAccessPointTypesDTO) (accessPointTypes []*entity.AccessPointType, err error)

	UpdateAccessPointType(ctx context.Context, patchUpdateDTO *dto.PatchUpdateAccessPointTypeDTO) (err error)

	IsAccessPointTypeSoftDeleted(ctx context.Context, accessPointTypeID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteAccessPointType(ctx context.Context, accessPointTypeID uuid.UUID) (err error)
	RestoreAccessPointType(ctx context.Context, accessPointTypeID uuid.UUID) (err error)
}

type AccessPointTypeUsecase struct {
	accessPointTypeService IAccessPointTypeService
	// accessPointRadioTemplateService AccessPointRadioTemplateService
	siteService ISiteService
}

func NewAccessPointTypeUsecase(
	accessPointTypeService IAccessPointTypeService,
	// accessPointRadioTemplateService AccessPointRadioTemplateService,
	siteService ISiteService,
) *AccessPointTypeUsecase {
	return &AccessPointTypeUsecase{
		accessPointTypeService: accessPointTypeService,
		// accessPointRadioTemplateService: accessPointRadioTemplateService,
		siteService: siteService,
	}
}

func (u *AccessPointTypeUsecase) CreateAccessPointType(ctx context.Context, createDTO *dto.CreateAccessPointTypeDTO) (accessPointTypeID uuid.UUID, err error) {
	_, err = u.siteService.GetSite(ctx, createDTO.SiteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Err(err).Msg("failed to create accessPointType: the site with provided site ID does not exist")
			return accessPointTypeID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check site existing")
		return
	}

	accessPointTypeID, err = u.accessPointTypeService.CreateAccessPointType(ctx, createDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create accessPointType")
		return
	}

	log.Info().Msgf("accessPointType %v successfully created", accessPointTypeID)
	return
}

func (u *AccessPointTypeUsecase) GetAccessPointType(ctx context.Context, accessPointTypeID uuid.UUID) (accessPointType *entity.AccessPointType, err error) {
	accessPointType, err = u.accessPointTypeService.GetAccessPointType(ctx, accessPointTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get accessPointType")
			return
		}
	}

	return
}

func (u *AccessPointTypeUsecase) GetAccessPointTypeDetailed(ctx context.Context, getDTO dto.GetAccessPointTypeDetailedDTO) (accessPointTypeDetailed *entity.AccessPointTypeDetailed, err error) {
	accessPointTypeDetailed, err = u.accessPointTypeService.GetAccessPointTypeDetailed(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to get access point type detailed")
		return
	}

	return
}

func (u *AccessPointTypeUsecase) GetAccessPointTypes(ctx context.Context, getDTO dto.GetAccessPointTypesDTO) (accessPointTypes []*entity.AccessPointType, err error) {
	accessPointTypes, err = u.accessPointTypeService.GetAccessPointTypes(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to get access point types")
		return
	}

	return
}

func (u *AccessPointTypeUsecase) PatchUpdateAccessPointType(ctx context.Context, patchUpdateDTO *dto.PatchUpdateAccessPointTypeDTO) (err error) {
	_, err = u.accessPointTypeService.GetAccessPointType(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Error().Err(err).Msg("failed to check accessPointType existing")
			return ErrNotFound
		}
	}

	err = u.accessPointTypeService.UpdateAccessPointType(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Err(err).Msg("accessPointType was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update accessPointType")
		return
	}

	return
}

func (u *AccessPointTypeUsecase) SoftDeleteAccessPointType(ctx context.Context, accessPointTypeID uuid.UUID) (err error) {
	isDeleted, err := u.accessPointTypeService.IsAccessPointTypeSoftDeleted(ctx, accessPointTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if accessPointType is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.accessPointTypeService.SoftDeleteAccessPointType(ctx, accessPointTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete accessPointType")
		return
	}

	return
}

func (u *AccessPointTypeUsecase) RestoreAccessPointType(ctx context.Context, accessPointTypeID uuid.UUID) (err error) {
	isDeleted, err := u.accessPointTypeService.IsAccessPointTypeSoftDeleted(ctx, accessPointTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if accessPointType is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.accessPointTypeService.RestoreAccessPointType(ctx, accessPointTypeID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore accessPointType")
		return
	}

	return
}
