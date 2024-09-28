package usecase

import (
	"context"
)

type HealthService interface {
	Health(ctx context.Context) (err error)
}

type HealthUsecase interface {
	Health(ctx context.Context) (err error)
}

type healthUsecase struct {
	healthService HealthService
}

func NewHealthUsecase(healthService HealthService) *healthUsecase {
	return &healthUsecase{healthService: healthService}
}

// Health pings database
func (u *healthUsecase) Health(ctx context.Context) (err error) {
	err = u.healthService.Health(ctx)
	return
}
