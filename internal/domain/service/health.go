package service

import (
	repository "location-backend/internal/adapters/db/postgres"
)

type HealthService interface {
	Health() (err error)
}

type healthService struct {
	repository repository.HealthRepo
}

func NewHealthService(repository repository.HealthRepo) *healthService {
	return &healthService{repository: repository}
}

func (s *healthService) Health() (err error) {
	err = s.repository.Health()
	return
}
