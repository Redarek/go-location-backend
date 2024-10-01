package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type AccessPointRadioTemplateService interface {
	CreateAccessPointRadioTemplate(ctx context.Context, createDTO *dto.CreateAccessPointRadioTemplateDTO) (accessPointRadioTemplateID uuid.UUID, err error)
	GetAccessPointRadioTemplate(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (accessPointRadioTemplate *entity.AccessPointRadioTemplate, err error)
	GetAccessPointRadioTemplates(ctx context.Context, getDTO dto.GetAccessPointRadioTemplatesDTO) (accessPointRadioTemplates []*entity.AccessPointRadioTemplate, err error)

	UpdateAccessPointRadioTemplate(ctx context.Context, patchUpdateDTO *dto.PatchUpdateAccessPointRadioTemplateDTO) (err error)

	IsAccessPointRadioTemplateSoftDeleted(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (isDeleted bool, err error)
	SoftDeleteAccessPointRadioTemplate(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (err error)
	RestoreAccessPointRadioTemplate(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (err error)
}

type AccessPointRadioTemplateUsecase struct {
	accessPointRadioTemplateService AccessPointRadioTemplateService
	accessPointTypeService          AccessPointTypeService
}

func NewAccessPointRadioTemplateUsecase(accessPointRadioTemplateService AccessPointRadioTemplateService, accessPointTypeService AccessPointTypeService) *AccessPointRadioTemplateUsecase {
	return &AccessPointRadioTemplateUsecase{
		accessPointRadioTemplateService: accessPointRadioTemplateService,
		accessPointTypeService:          accessPointTypeService,
	}
}

func (u *AccessPointRadioTemplateUsecase) CreateAccessPointRadioTemplate(ctx context.Context, createDTO *dto.CreateAccessPointRadioTemplateDTO) (accessPointRadioTemplateID uuid.UUID, err error) {
	_, err = u.accessPointTypeService.GetAccessPointType(ctx, createDTO.AccessPointTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Msg("failed to create access point radio template: the access point type with provided ID does not exist")
			return accessPointRadioTemplateID, ErrNotFound
		}

		log.Error().Err(err).Msg("failed to check access point type existing")
		return
	}

	accessPointRadioTemplateID, err = u.accessPointRadioTemplateService.CreateAccessPointRadioTemplate(ctx, createDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create access point radio template")
		return
	}

	log.Debug().Msgf("access point radio template %v successfully created", accessPointRadioTemplateID)
	return
}

func (u *AccessPointRadioTemplateUsecase) GetAccessPointRadioTemplate(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (accessPointRadioTemplate *entity.AccessPointRadioTemplate, err error) {
	accessPointRadioTemplate, err = u.accessPointRadioTemplateService.GetAccessPointRadioTemplate(ctx, accessPointRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get access point radio template")
			return
		}
	}

	return
}

func (u *AccessPointRadioTemplateUsecase) GetAccessPointRadioTemplates(ctx context.Context, getDTO dto.GetAccessPointRadioTemplatesDTO) (accessPointRadioTemplates []*entity.AccessPointRadioTemplate, err error) {
	accessPointRadioTemplates, err = u.accessPointRadioTemplateService.GetAccessPointRadioTemplates(ctx, getDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get access point radio templates")
			return
		}
	}

	return
}

func (u *AccessPointRadioTemplateUsecase) PatchUpdateAccessPointRadioTemplate(ctx context.Context, patchUpdateDTO *dto.PatchUpdateAccessPointRadioTemplateDTO) (err error) {
	_, err = u.accessPointRadioTemplateService.GetAccessPointRadioTemplate(ctx, patchUpdateDTO.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Msg("access point radio template was not found")
			return ErrNotFound
		}
	}

	err = u.accessPointRadioTemplateService.UpdateAccessPointRadioTemplate(ctx, patchUpdateDTO)
	if err != nil {
		if errors.Is(err, ErrNotUpdated) {
			log.Info().Msg("accessPointRadioTemplate was not updated")
			return ErrNotUpdated
		}
		log.Error().Err(err).Msg("failed to patch update access point radio template")
		return
	}

	return
}

func (u *AccessPointRadioTemplateUsecase) SoftDeleteAccessPointRadioTemplate(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (err error) {
	isDeleted, err := u.accessPointRadioTemplateService.IsAccessPointRadioTemplateSoftDeleted(ctx, accessPointRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if access point radio template is soft deleted")
			return
		}
	}

	if isDeleted {
		return ErrAlreadySoftDeleted
	}

	err = u.accessPointRadioTemplateService.SoftDeleteAccessPointRadioTemplate(ctx, accessPointRadioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete access point radio template")
		return
	}

	return
}

func (u *AccessPointRadioTemplateUsecase) RestoreAccessPointRadioTemplate(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (err error) {
	isDeleted, err := u.accessPointRadioTemplateService.IsAccessPointRadioTemplateSoftDeleted(ctx, accessPointRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to check if access point radio template is soft deleted")
			return
		}
	}

	if !isDeleted {
		return ErrAlreadyExists
	}

	err = u.accessPointRadioTemplateService.RestoreAccessPointRadioTemplate(ctx, accessPointRadioTemplateID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore access point radio template")
		return
	}

	return
}
