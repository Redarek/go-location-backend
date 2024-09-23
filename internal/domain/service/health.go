package service

import (
	"context"

	repository "location-backend/internal/adapters/db/postgres"
)

type HealthService interface {
	Health(ctx context.Context) (err error)
}

type healthService struct {
	repository repository.HealthRepo
}

func NewHealthService(repository repository.HealthRepo) *healthService {
	return &healthService{repository: repository}
}

func (s *healthService) Health(ctx context.Context) (err error) {
	err = s.repository.Health(ctx)
	return
}
