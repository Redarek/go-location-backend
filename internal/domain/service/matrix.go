package service

import (
	"context"

	"github.com/google/uuid"

	"location-backend/internal/domain/entity"
)

type IMatrixRepo interface {
	Create(ctx context.Context, points []*entity.Point, matrixPoints []*entity.MatrixPoint) (err error)

	Delete(ctx context.Context, floorID uuid.UUID) (deletedCount int64, err error)
}

type matrixService struct {
	repository IMatrixRepo
}

func NewMatrixService(repository IMatrixRepo) *matrixService {
	return &matrixService{repository: repository}
}

func (s *matrixService) CreateMatrix(ctx context.Context, points []*entity.Point, matrixPoints []*entity.MatrixPoint) (err error) {
	err = s.repository.Create(ctx, points, matrixPoints)
	return
}

func (s *matrixService) DeletePoints(ctx context.Context, floorID uuid.UUID) (deletedCount int64, err error) {
	deletedCount, err = s.repository.Delete(ctx, floorID)
	return
}
