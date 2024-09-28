package usecase

import (
	"context"
)

type HealthService interface {
	Health(ctx context.Context) (err error)
}

type HealthUsecase struct {
	healthService HealthService
}

func NewHealthUsecase(healthService HealthService) *HealthUsecase {
	return &HealthUsecase{healthService: healthService}
}

// Health pings database
func (u *HealthUsecase) Health(ctx context.Context) (err error) {
	err = u.healthService.Health(ctx)
	return
}
