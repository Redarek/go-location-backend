package v1

import (
	// "encoding/json"

	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
)

const (
	getRoleByNameURL = "/"
	createRoleURL    = "/"
)

type roleHandler struct {
	usecase usecase.RoleUsecase
}

// Регистрирует новый handler
func NewRoleHandler(usecase usecase.RoleUsecase) *roleHandler {
	return &roleHandler{usecase: usecase}
}

// Регистрирует маршруты для user
func (h *roleHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Get(getRoleByNameURL, h.GetRoleByName)
	router.Post(createRoleURL, h.CreateRole)

	return router
}

func (h *roleHandler) CreateRole(ctx *fiber.Ctx) error {
	var dto http_dto.CreateRoleDTO
	err := ctx.BodyParser(&dto)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse user request body")
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.CreateRoleDTO{
		Name: dto.Name,
	}

	roleID, err := h.usecase.CreateRole(domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrAlreadyExists) {
			return ctx.Status(fiber.StatusConflict).SendString("Role is already exists")
		}

		log.Error().Err(err).Msg("failed to create tole")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create role")
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": roleID})
}

func (h *roleHandler) GetRoleByName(ctx *fiber.Ctx) error {
	var role entity.Role

	if ctx.Query("id") != "" {
		uuid, err := uuid.Parse(ctx.Query("id"))
		if err != nil {
			log.Error().Err(err).Msg("failed to parse id as UUID")
			return ctx.Status(fiber.StatusBadRequest).SendString("Wrong ID: this is not the UUID")
		}

		var dto http_dto.GetRoleDTO = http_dto.GetRoleDTO{
			ID: uuid,
		}

		// TODO validate

		// Mapping http DTO -> domain DTO
		domainDTO := domain_dto.GetRoleDTO{
			ID: dto.ID,
		}

		role, err = h.usecase.GetRole(domainDTO)
		if err != nil {
			if errors.Is(err, usecase.ErrNotFound) {
				return ctx.Status(fiber.StatusNoContent).SendString("Role not found")
			}

			log.Error().Err(err).Msg("failed to get role by name")
			return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to get role by name")
		}
	} else if ctx.Query("name") != "" {
		var dto http_dto.GetRoleByNameDTO = http_dto.GetRoleByNameDTO{
			Name: ctx.Query("name"),
		}

		// TODO validate

		// Mapping http DTO -> domain DTO
		domainDTO := domain_dto.GetRoleByNameDTO{
			Name: dto.Name,
		}

		var err error
		role, err = h.usecase.GetRoleByName(domainDTO)
		if err != nil {
			if errors.Is(err, usecase.ErrNotFound) {
				return ctx.Status(fiber.StatusNoContent).SendString("Role not found")
			}

			log.Error().Err(err).Msg("failed to get role by name")
			return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to get role by name")
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body: not provided id or name parameter")
	}

	roleDTO := http_dto.RoleDTO{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: role.DeletedAt,
	}

	return ctx.Status(fiber.StatusOK).JSON(roleDTO)
}
