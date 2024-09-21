package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/service"
)

type RoleUsecase interface {
	CreateRole(dto domain_dto.CreateRoleDTO) (roleID uuid.UUID, err error)
	GetRole(dto domain_dto.GetRoleDTO) (role entity.Role, err error)
	GetRoleByName(dto domain_dto.GetRoleByNameDTO) (role entity.Role, err error)
}

type roleUsecase struct {
	roleService service.RoleService
}

func NewRoleUsecase(roleService service.RoleService) *roleUsecase {
	return &roleUsecase{roleService: roleService}
}

func (u *roleUsecase) CreateRole(dto domain_dto.CreateRoleDTO) (roleID uuid.UUID, err error) {
	_, err = u.roleService.GetRoleByName(dto.Name)
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

	roleID, err = u.roleService.CreateRole(createRoleDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create role")
		return
	}

	log.Info().Msgf("role %v successfully created", dto.Name)
	return
}

func (u *roleUsecase) GetRole(dto domain_dto.GetRoleDTO) (role entity.Role, err error) {
	role, err = u.roleService.GetRole(dto.ID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return entity.Role{}, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get role")
			return
		}
	}

	return
}

func (u *roleUsecase) GetRoleByName(dto domain_dto.GetRoleByNameDTO) (role entity.Role, err error) {
	role, err = u.roleService.GetRoleByName(dto.Name)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return entity.Role{}, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get role")
			return
		}
	}

	return
}
