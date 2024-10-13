package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
)

type AccessPointRadioRepo interface {
	Create(ctx context.Context, createAccessPointRadioDTO *dto.CreateAccessPointRadioDTO) (accessPointRadioID uuid.UUID, err error)
	GetOne(ctx context.Context, accessPointRadioID uuid.UUID) (accessPointRadio *entity.AccessPointRadio, err error)
	GetAll(ctx context.Context, accessPointID uuid.UUID, limit, offset int) (accessPointRadios []*entity.AccessPointRadio, err error)

	Update(ctx context.Context, updateAccessPointRadioDTO *dto.PatchUpdateAccessPointRadioDTO) (err error)

	IsAccessPointRadioSoftDeleted(ctx context.Context, accessPointRadioID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, accessPointRadioID uuid.UUID) (err error)
	Restore(ctx context.Context, accessPointRadioID uuid.UUID) (err error)
}

type accessPointRadioService struct {
	repository AccessPointRadioRepo
}

func NewAccessPointRadioService(repository AccessPointRadioRepo) *accessPointRadioService {
	return &accessPointRadioService{repository: repository}
}

func (s *accessPointRadioService) CreateAccessPointRadio(ctx context.Context, createAccessPointRadioDTO *dto.CreateAccessPointRadioDTO) (accessPointRadioID uuid.UUID, err error) {
	accessPointRadioID, err = s.repository.Create(ctx, createAccessPointRadioDTO)
	return
}

func (s *accessPointRadioService) GetAccessPointRadio(ctx context.Context, accessPointRadioID uuid.UUID) (accessPointRadio *entity.AccessPointRadio, err error) {
	accessPointRadio, err = s.repository.GetOne(ctx, accessPointRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return accessPointRadio, usecase.ErrNotFound
		}
	}

	return
}

func (s *accessPointRadioService) GetAccessPointRadios(ctx context.Context, dto dto.GetAccessPointRadiosDTO) (accessPointRadios []*entity.AccessPointRadio, err error) {
	accessPointRadios, err = s.repository.GetAll(ctx, dto.AccessPointID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return accessPointRadios, usecase.ErrNotFound
		}
	}

	return
}

// TODO PUT update
func (s *accessPointRadioService) UpdateAccessPointRadio(ctx context.Context, updateAccessPointRadioDTO *dto.PatchUpdateAccessPointRadioDTO) (err error) {
	err = s.repository.Update(ctx, updateAccessPointRadioDTO)
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

func (s *accessPointRadioService) IsAccessPointRadioSoftDeleted(ctx context.Context, accessPointRadioID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.repository.IsAccessPointRadioSoftDeleted(ctx, accessPointRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, usecase.ErrNotFound
		}
	}

	return
}

func (s *accessPointRadioService) SoftDeleteAccessPointRadio(ctx context.Context, accessPointRadioID uuid.UUID) (err error) {
	err = s.repository.SoftDelete(ctx, accessPointRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
	}

	return
}

func (s *accessPointRadioService) RestoreAccessPointRadio(ctx context.Context, accessPointRadioID uuid.UUID) (err error) {
	err = s.repository.Restore(ctx, accessPointRadioID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}
	}

	return
}
