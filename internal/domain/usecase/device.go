package usecase

import (
	"context"

	"github.com/google/uuid"

	"location-backend/internal/domain/entity"
)

type IDeviceService interface {
	GetDevices(ctx context.Context, mac string, floorID uuid.UUID, limit, offset int) (devices []*entity.Device, err error)
}

type DeviceUsecase struct {
	deviceService IDeviceService
}

func NewDeviceUsecase(deviceService IDeviceService) *DeviceUsecase {
	return &DeviceUsecase{
		deviceService: deviceService,
	}
}
