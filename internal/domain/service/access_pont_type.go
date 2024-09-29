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

type AccessPointTypeRepo interface {
	Create(ctx context.Context, createAccessPointTypeDTO *dto.CreateAccessPointTypeDTO) (accessPointTypeID uuid.UUID, err error)
	GetOne(ctx context.Context, accessPointTypeID uuid.UUID) (accessPointType *entity.AccessPointType, err error)
	// GetOneDetailed(ctx context.Context, accessPointTypeID uuid.UUID) (accessPointType *entity.AccessPointTypeDetailed, err error) // TODO

	GetAll(ctx context.Context, siteID uuid.UUID, limit, offset int) (accessPointTypes []*entity.AccessPointType, err error)

	Update(ctx context.Context, updateAccessPointTypeDTO *dto.PatchUpdateAccessPointTypeDTO) (err error)

	IsAccessPointTypeSoftDeleted(ctx context.Context, accessPointTypeID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, accessPointTypeID uuid.UUID) (err error)
	Restore(ctx context.Context, accessPointTypeID uuid.UUID) (err error)
}

type accessPointTypeService struct {
	repository AccessPointTypeRepo
}

func NewAccessPointTypeService(repository AccessPointTypeRepo) *accessPointTypeService {
	return &accessPointTypeService{repository: repository}
}

func (s *accessPointTypeService) CreateAccessPointType(ctx context.Context, createAccessPointTypeDTO *dto.CreateAccessPointTypeDTO) (accessPointTypeID uuid.UUID, err error) {
	accessPointTypeID, err = s.repository.Create(ctx, createAccessPointTypeDTO)
	if err != nil {
		// TODO улучшить лог
		log.Error().Msg("failed to create accessPointType")
		return
	}

	return accessPointTypeID, nil
}

func (s *accessPointTypeService) GetAccessPointType(ctx context.Context, accessPointTypeID uuid.UUID) (accessPointType *entity.AccessPointType, err error) {
	accessPointType, err = s.repository.GetOne(ctx, accessPointTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return accessPointType, usecase.ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Msg("failed to retrieve accessPointType")
		return
	}

	return
}

func (s *accessPointTypeService) GetAccessPointTypes(ctx context.Context, dto dto.GetAccessPointTypesDTO) (accessPointTypes []*entity.AccessPointType, err error) {
	accessPointTypes, err = s.repository.GetAll(ctx, dto.SiteID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return accessPointTypes, usecase.ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Msg("failed to retrieve accessPointType")
		return
	}

	return
}

// TODO PUT update
func (s *accessPointTypeService) UpdateAccessPointType(ctx context.Context, updateAccessPointTypeDTO *dto.PatchUpdateAccessPointTypeDTO) (err error) {
	err = s.repository.Update(ctx, updateAccessPointTypeDTO)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
		if errors.Is(err, ErrNotUpdated) {
			return usecase.ErrNotUpdated
		}
		// TODO улучшить лог
		log.Error().Msg("failed to update accessPointType")
		return
	}

	return
}

func (s *accessPointTypeService) IsAccessPointTypeSoftDeleted(ctx context.Context, accessPointTypeID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.repository.IsAccessPointTypeSoftDeleted(ctx, accessPointTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, usecase.ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Msg("failed to retrieve accessPointType")
		return
	}

	return
}

func (s *accessPointTypeService) SoftDeleteAccessPointType(ctx context.Context, accessPointTypeID uuid.UUID) (err error) {
	err = s.repository.SoftDelete(ctx, accessPointTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Msg("failed to soft delete accessPointType")
		return
	}

	return
}

func (s *accessPointTypeService) RestoreAccessPointType(ctx context.Context, accessPointTypeID uuid.UUID) (err error) {
	err = s.repository.Restore(ctx, accessPointTypeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Msg("failed to restore accessPointType")
		return
	}

	return
}
