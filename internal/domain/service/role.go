package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	repository "location-backend/internal/adapters/db/postgres"
	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type RoleService interface {
	CreateRole(ctx context.Context, userCreate dto.CreateRoleDTO) (roleID uuid.UUID, err error)
	GetRole(ctx context.Context, roleID uuid.UUID) (role *entity.Role, err error)
	GetRoleByName(ctx context.Context, name string) (role *entity.Role, err error)
}

type roleService struct {
	repository repository.RoleRepo
}

func NewRoleService(repository repository.RoleRepo) *roleService {
	return &roleService{repository: repository}
}

func (s *roleService) CreateRole(ctx context.Context, createRoleDTO dto.CreateRoleDTO) (roleID uuid.UUID, err error) {
	roleID, err = s.repository.Create(ctx, createRoleDTO)
	if err != nil {
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to create role")
		return
	}

	return roleID, nil
}

func (s *roleService) GetRole(ctx context.Context, roleID uuid.UUID) (role *entity.Role, err error) {
	role, err = s.repository.GetOne(ctx, roleID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve role")
		return
	}

	return
}

func (s *roleService) GetRoleByName(ctx context.Context, name string) (role *entity.Role, err error) {
	role, err = s.repository.GetOneByName(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve role")
		return
	}

	return
}
