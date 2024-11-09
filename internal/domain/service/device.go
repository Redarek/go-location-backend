package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
)

type IDeviceRepo interface {
	GetAll(ctx context.Context, mac string, floorID uuid.UUID, limit, offset int) (devices []*entity.Device, err error)
}

type deviceService struct {
	repository IDeviceRepo
}

func NewDeviceService(repository IDeviceRepo) *deviceService {
	return &deviceService{repository: repository}
}

func (s *deviceService) GetDevices(ctx context.Context, mac string, floorID uuid.UUID, limit, offset int) (devices []*entity.Device, err error) {
	devices, err = s.repository.GetAll(ctx, mac, floorID, limit, offset)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return devices, usecase.ErrNotFound
		}

		return
	}

	return
}
