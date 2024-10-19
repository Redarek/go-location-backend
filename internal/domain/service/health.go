package service

import (
	"context"
)

type IHealthRepo interface {
	Health(ctx context.Context) (err error)
}

type healthService struct {
	repository IHealthRepo
}

func NewHealthService(repository IHealthRepo) *healthService {
	return &healthService{repository: repository}
}

func (s *healthService) Health(ctx context.Context) (err error) {
	err = s.repository.Health(ctx)
	return
}
