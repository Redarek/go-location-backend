package usecase

import (
	"location-backend/internal/domain/service"
)

type HealthUsecase interface {
	Health() (err error)
}

type healthUsecase struct {
	healthService service.HealthService
}

func NewHealthUsecase(healthService service.HealthService) *healthUsecase {
	return &healthUsecase{healthService: healthService}
}

// Health pings database
func (u *healthUsecase) Health() (err error) {
	err = u.healthService.Health()
	return
}
