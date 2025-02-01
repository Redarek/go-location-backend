package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	// "github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
	"location-backend/pkg/utils"
)

type IAccessPointRadioTemplateRepo interface {
	Create(ctx context.Context, createAccessPointRadioTemplateDTO *dto.CreateAccessPointRadioTemplateDTO) (accessPointRadioTemplateID uuid.UUID, err error)
	GetOne(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (accessPointRadioTemplate *entity.AccessPointRadioTemplate, err error)
	GetAll(ctx context.Context, accessPointTypeID uuid.UUID, limit, offset int) (accessPointRadioTemplates []*entity.AccessPointRadioTemplate, err error)

	Update(ctx context.Context, updateAccessPointRadioTemplateDTO *dto.PatchUpdateAccessPointRadioTemplateDTO) (err error)

	IsAccessPointRadioTemplateSoftDeleted(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (err error)
	Restore(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (err error)
}

type accessPointRadioTemplateService struct {
	repository IAccessPointRadioTemplateRepo
}

func NewAccessPointRadioTemplateService(repository IAccessPointRadioTemplateRepo) *accessPointRadioTemplateService {
	return &accessPointRadioTemplateService{repository: repository}
}

func (s *accessPointRadioTemplateService) CreateAccessPointRadioTemplate(ctx context.Context, createAccessPointRadioTemplateDTO *dto.CreateAccessPointRadioTemplateDTO) (accessPointRadioTemplateID uuid.UUID, err error) {
	accessPointRadioTemplateID, err = s.repository.Create(ctx, createAccessPointRadioTemplateDTO)
	return
}

func (s *accessPointRadioTemplateService) GetAccessPointRadioTemplate(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (accessPointRadioTemplate *entity.AccessPointRadioTemplate, err error) {
	accessPointRadioTemplate, err = s.repository.GetOne(ctx, accessPointRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return accessPointRadioTemplate, usecase.ErrNotFound
		}
	}

	return
}

func (s *accessPointRadioTemplateService) GetAccessPointRadioTemplates(ctx context.Context, dto *dto.GetAccessPointRadioTemplatesDTO) (accessPointRadioTemplates []*entity.AccessPointRadioTemplate, err error) {
	accessPointRadioTemplates, err = s.repository.GetAll(ctx, dto.AccessPointTypeID, dto.Size, utils.GetOffset(dto.Page, dto.Size))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return accessPointRadioTemplates, usecase.ErrNotFound
		}
	}

	return
}

// TODO PUT update
func (s *accessPointRadioTemplateService) UpdateAccessPointRadioTemplate(ctx context.Context, updateAccessPointRadioTemplateDTO *dto.PatchUpdateAccessPointRadioTemplateDTO) (err error) {
	err = s.repository.Update(ctx, updateAccessPointRadioTemplateDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
		if errors.Is(err, ErrNotUpdated) {
			return usecase.ErrNotUpdated
		}
	}

	return
}

func (s *accessPointRadioTemplateService) IsAccessPointRadioTemplateSoftDeleted(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.repository.IsAccessPointRadioTemplateSoftDeleted(ctx, accessPointRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, usecase.ErrNotFound
		}
	}

	return
}

func (s *accessPointRadioTemplateService) SoftDeleteAccessPointRadioTemplate(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (err error) {
	err = s.repository.SoftDelete(ctx, accessPointRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
	}

	return
}

func (s *accessPointRadioTemplateService) RestoreAccessPointRadioTemplate(ctx context.Context, accessPointRadioTemplateID uuid.UUID) (err error) {
	err = s.repository.Restore(ctx, accessPointRadioTemplateID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
	}

	return
}
