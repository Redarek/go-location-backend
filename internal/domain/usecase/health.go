package usecase

import (
	"context"
)

type IHealthService interface {
	Health(ctx context.Context) (err error)
}

type HealthUsecase struct {
	healthService IHealthService
}

func NewHealthUsecase(healthService IHealthService) *HealthUsecase {
	return &HealthUsecase{healthService: healthService}
}

// Health pings database
func (u *HealthUsecase) Health(ctx context.Context) (err error) {
	err = u.healthService.Health(ctx)
	return
}
