package service

import (
	"context"
)

type HealthRepo interface {
	Health(ctx context.Context) (err error)
}

type healthService struct {
	repository HealthRepo
}

func NewHealthService(repository HealthRepo) *healthService {
	return &healthService{repository: repository}
}

func (s *healthService) Health(ctx context.Context) (err error) {
	err = s.repository.Health(ctx)
	return
}
