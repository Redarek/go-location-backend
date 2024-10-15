package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
)

type AccessPointRepo interface {
	Create(ctx context.Context, createAccessPointDTO *dto.CreateAccessPointDTO) (accessPointID uuid.UUID, err error)
	GetOne(ctx context.Context, accessPointID uuid.UUID) (accessPoint *entity.AccessPoint, err error)
	GetOneDetailed(ctx context.Context, accessPointID uuid.UUID) (apDetailed *entity.AccessPointDetailed, err error)
	GetAll(ctx context.Context, floorID uuid.UUID, limit, offset int) (accessPoints []*entity.AccessPoint, err error)
	GetAllDetailed(ctx context.Context, floorID uuid.UUID, limit, offset int) (accessPointsDetailed []*entity.AccessPointDetailed, err error)

	Update(ctx context.Context, updateAccessPointDTO *dto.PatchUpdateAccessPointDTO) (err error)

	IsAccessPointSoftDeleted(ctx context.Context, accessPointID uuid.UUID) (isDeleted bool, err error)
	SoftDelete(ctx context.Context, accessPointID uuid.UUID) (err error)
	Restore(ctx context.Context, accessPointID uuid.UUID) (err error)
}

type accessPointService struct {
	accessPointRepo      AccessPointRepo
	accessPointTypeRepo  AccessPointTypeRepo
	accessPointRadioRepo AccessPointRadioRepo
}

func NewAccessPointService(
	accessPointRepo AccessPointRepo,
	accessPointTypeRepo AccessPointTypeRepo,
	accessPointRadioRepo AccessPointRadioRepo,
) *accessPointService {
	return &accessPointService{
		accessPointRepo:      accessPointRepo,
		accessPointTypeRepo:  accessPointTypeRepo,
		accessPointRadioRepo: accessPointRadioRepo,
	}
}

func (s *accessPointService) CreateAccessPoint(ctx context.Context, createAccessPointDTO *dto.CreateAccessPointDTO) (accessPointID uuid.UUID, err error) {
	accessPointID, err = s.accessPointRepo.Create(ctx, createAccessPointDTO)
	return
}

func (s *accessPointService) GetAccessPoint(ctx context.Context, accessPointID uuid.UUID) (accessPoint *entity.AccessPoint, err error) {
	accessPoint, err = s.accessPointRepo.GetOne(ctx, accessPointID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return accessPoint, usecase.ErrNotFound
		}

		return
	}

	return
}

// func (s *accessPointService) GetAccessPointDetailed(ctx context.Context, dto dto.GetAccessPointDetailedDTO) (accessPointDetailed *entity.AccessPointDetailed, err error) {
// 	accessPoint, err := s.accessPointRepo.GetOne(ctx, dto.ID)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			log.Debug().Msg("access point was not found")
// 			return accessPointDetailed, usecase.ErrNotFound
// 		}
// 		// TODO улучшить лог
// 		log.Error().Msg("failed to retrieve access point")
// 		return
// 	}

// 	accessPointType, err := s.accessPointTypeRepo.GetOne(ctx, accessPoint.AccessPointTypeID)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			log.Error().Msg("access point type was not found")
// 			return
// 		}

// 		// TODO улучшить лог
// 		log.Error().Msg("failed to retrieve access point type")
// 		return
// 	}

// 	accessPointRadios, err := s.accessPointRadioRepo.GetAll(ctx, accessPoint.ID, dto.Limit, dto.Offset)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			log.Debug().Msg("access point radios were not found")
// 			err = nil
// 		} else {
// 			// TODO улучшить лог
// 			log.Error().Msg("failed to retrieve access point radios")
// 			return
// 		}
// 	}

// 	accessPointDetailed = &entity.AccessPointDetailed{
// 		AccessPoint:     *accessPoint,
// 		AccessPointType: *accessPointType,
// 		Radios:          accessPointRadios,
// 	}

// 	return
// }

func (s *accessPointService) GetAccessPointDetailed(ctx context.Context, dto dto.GetAccessPointDetailedDTO) (apDetailed *entity.AccessPointDetailed, err error) {
	apDetailed, err = s.accessPointRepo.GetOneDetailed(ctx, dto.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *accessPointService) GetAccessPoints(ctx context.Context, dto dto.GetAccessPointsDTO) (accessPoints []*entity.AccessPoint, err error) {
	accessPoints, err = s.accessPointRepo.GetAll(ctx, dto.FloorID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return accessPoints, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *accessPointService) GetAccessPointsDetailed(ctx context.Context, dto dto.GetAccessPointsDetailedDTO) (accessPoints []*entity.AccessPointDetailed, err error) {
	accessPoints, err = s.accessPointRepo.GetAllDetailed(ctx, dto.FloorID, dto.Limit, dto.Offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return accessPoints, usecase.ErrNotFound
		}

		return
	}

	return
}

// TODO PUT update
func (s *accessPointService) UpdateAccessPoint(ctx context.Context, updateAccessPointDTO *dto.PatchUpdateAccessPointDTO) (err error) {
	err = s.accessPointRepo.Update(ctx, updateAccessPointDTO)
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

func (s *accessPointService) IsAccessPointSoftDeleted(ctx context.Context, accessPointID uuid.UUID) (isDeleted bool, err error) {
	isDeleted, err = s.accessPointRepo.IsAccessPointSoftDeleted(ctx, accessPointID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *accessPointService) SoftDeleteAccessPoint(ctx context.Context, accessPointID uuid.UUID) (err error) {
	err = s.accessPointRepo.SoftDelete(ctx, accessPointID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}

		return
	}

	return
}

func (s *accessPointService) RestoreAccessPoint(ctx context.Context, accessPointID uuid.UUID) (err error) {
	err = s.accessPointRepo.Restore(ctx, accessPointID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return usecase.ErrNotFound
		}

		return
	}

	return
}
