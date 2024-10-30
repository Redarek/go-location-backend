package service

import (
	"context"

	"location-backend/internal/domain/dto"
)

type IMatrixRepo interface {
	Create(ctx context.Context, createMatrixDTOs []*dto.CreateMatrixDTO) (err error)
}

type matrixService struct {
	repository IMatrixRepo
}

func NewMatrixService(repository IMatrixRepo) *matrixService {
	return &matrixService{repository: repository}
}

func (s *matrixService) CreateMatrix(ctx context.Context, createMatrixDTO []*dto.CreateMatrixDTO) (err error) {
	err = s.repository.Create(ctx, createMatrixDTO)
	return
}
