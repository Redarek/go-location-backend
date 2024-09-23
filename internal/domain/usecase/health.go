package usecase

import (
	"context"

	"location-backend/internal/domain/service"
)

type HealthUsecase interface {
	Health(ctx context.Context) (err error)
}

type healthUsecase struct {
	healthService service.HealthService
}

func NewHealthUsecase(healthService service.HealthService) *healthUsecase {
	return &healthUsecase{healthService: healthService}
}

// Health pings database
func (u *healthUsecase) Health(ctx context.Context) (err error) {
	err = u.healthService.Health(ctx)
	return
}
