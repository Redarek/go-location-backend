package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
)

type RoleRepo interface {
	Create(ctx context.Context, createRoleDTO *dto.CreateRoleDTO) (roleID uuid.UUID, err error)
	GetOne(ctx context.Context, roleID uuid.UUID) (role *entity.Role, err error)
	GetOneByName(ctx context.Context, name string) (role *entity.Role, err error)
}

type roleService struct {
	repository RoleRepo
}

func NewRoleService(repository RoleRepo) *roleService {
	return &roleService{repository: repository}
}

func (s *roleService) CreateRole(ctx context.Context, createRoleDTO *dto.CreateRoleDTO) (roleID uuid.UUID, err error) {
	roleID, err = s.repository.Create(ctx, createRoleDTO)
	return
}

func (s *roleService) GetRole(ctx context.Context, roleID uuid.UUID) (role *entity.Role, err error) {
	role, err = s.repository.GetOne(ctx, roleID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, usecase.ErrNotFound
		}
		return
	}

	return
}

func (s *roleService) GetRoleByName(ctx context.Context, name string) (role *entity.Role, err error) {
	role, err = s.repository.GetOneByName(ctx, name)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, usecase.ErrNotFound
		}
		return
	}

	return
}
