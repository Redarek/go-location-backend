package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/service"
)

type RoleUsecase interface {
	CreateRole(ctx context.Context, dto domain_dto.CreateRoleDTO) (roleID uuid.UUID, err error)
	GetRole(ctx context.Context, dto domain_dto.GetRoleDTO) (roleDTO *domain_dto.RoleDTO, err error)
	GetRoleByName(ctx context.Context, dto domain_dto.GetRoleByNameDTO) (roleDTO *domain_dto.RoleDTO, err error)
}

type roleUsecase struct {
	roleService service.RoleService
}

func NewRoleUsecase(roleService service.RoleService) *roleUsecase {
	return &roleUsecase{roleService: roleService}
}

func (u *roleUsecase) CreateRole(ctx context.Context, dto domain_dto.CreateRoleDTO) (roleID uuid.UUID, err error) {
	_, err = u.roleService.GetRoleByName(ctx, dto.Name)
	if err != nil {
		// If error except ErrNotFound
		if !errors.Is(err, service.ErrNotFound) {
			log.Error().Err(err).Msg("failed to check user existing")
			return
		}
	} else { // If already exists
		return roleID, ErrAlreadyExists
	}

	var createRoleDTO domain_dto.CreateRoleDTO = domain_dto.CreateRoleDTO{
		Name: dto.Name,
	}

	roleID, err = u.roleService.CreateRole(ctx, createRoleDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create role")
		return
	}

	log.Info().Msgf("role %v successfully created", dto.Name)
	return
}

func (u *roleUsecase) GetRole(ctx context.Context, dto domain_dto.GetRoleDTO) (roleDTO *domain_dto.RoleDTO, err error) {
	role, err := u.roleService.GetRole(ctx, dto.ID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get role")
			return
		}
	}

	// Mapping domain entity -> domain DTO
	roleDTO = &domain_dto.RoleDTO{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: role.DeletedAt,
	}

	return
}

func (u *roleUsecase) GetRoleByName(ctx context.Context, dto domain_dto.GetRoleByNameDTO) (roleDTO *domain_dto.RoleDTO, err error) {
	role, err := u.roleService.GetRoleByName(ctx, dto.Name)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get role")
			return
		}
	}

	// Mapping domain entity -> domain DTO
	roleDTO = &domain_dto.RoleDTO{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: role.DeletedAt,
	}

	return
}
