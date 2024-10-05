package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type RoleService interface {
	CreateRole(ctx context.Context, createRoleDTO *dto.CreateRoleDTO) (roleID uuid.UUID, err error)
	GetRole(ctx context.Context, roleID uuid.UUID) (role *entity.Role, err error)
	GetRoleByName(ctx context.Context, name string) (role *entity.Role, err error)
}

type RoleUsecase struct {
	roleService RoleService
}

func NewRoleUsecase(roleService RoleService) *RoleUsecase {
	return &RoleUsecase{roleService: roleService}
}

func (u *RoleUsecase) CreateRole(ctx context.Context, dto *dto.CreateRoleDTO) (roleID uuid.UUID, err error) {
	_, err = u.roleService.GetRoleByName(ctx, dto.Name)
	if err != nil {
		// If error except ErrNotFound
		if !errors.Is(err, ErrNotFound) {
			log.Error().Err(err).Msg("failed to check user existing")
			return
		}
	} else { // If already exists
		return roleID, ErrAlreadyExists
	}

	roleID, err = u.roleService.CreateRole(ctx, dto)
	if err != nil {
		log.Error().Err(err).Msg("failed to create role")
		return
	}

	log.Info().Msgf("role %v successfully created", dto.Name)
	return
}

func (u *RoleUsecase) GetRole(ctx context.Context, roleID uuid.UUID) (role *entity.Role, err error) {
	role, err = u.roleService.GetRole(ctx, roleID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get role")
			return
		}
	}

	return
}

func (u *RoleUsecase) GetRoleByName(ctx context.Context, name string) (role *entity.Role, err error) {
	role, err = u.roleService.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get role")
			return
		}
	}

	return
}
