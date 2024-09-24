package v1

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/usecase"
	"location-backend/pkg/httperrors"
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
	router.Get(getRoleByNameURL, h.GetRoleByIDorName)
	router.Post(createRoleURL, h.CreateRole)

	return router
}

func (h *roleHandler) CreateRole(ctx *fiber.Ctx) error {
	var dto http_dto.CreateRoleDTO
	err := ctx.BodyParser(&dto)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse user request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse role request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := &domain_dto.CreateRoleDTO{
		Name: dto.Name,
	}

	roleID, err := h.usecase.CreateRole(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrAlreadyExists) {
			return ctx.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"Role is already exists",
				"",
				nil,
			))
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the role")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the role",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": roleID})
}

func (h *roleHandler) GetRoleByIDorName(ctx *fiber.Ctx) error {
	var role *domain_dto.RoleDTO

	if ctx.Query("id") != "" {
		roleID, err := uuid.Parse(ctx.Query("id"))
		if err != nil {
			log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
			return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"Invalid ID",
				"Failed to parse 'id' as UUID",
				nil,
			))
		}

		// var dto http_dto.GetRoleDTO = http_dto.GetRoleDTO{
		// 	ID: uuid,
		// }

		// TODO validate

		// Mapping http DTO -> domain DTO
		// domainDTO := &domain_dto.GetRoleDTO{
		// 	ID: dto.ID,
		// }

		role, err = h.usecase.GetRole(context.Background(), roleID)
		if err != nil {
			if errors.Is(err, usecase.ErrNotFound) {
				ctx.Status(fiber.StatusNoContent)
				return nil
			}

			log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the role by ID")
			return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
				fiber.StatusInternalServerError,
				"An unexpected error has occurred while trying to retrieve the role by ID",
				"",
				nil,
			))
		}
	} else if ctx.Query("name") != "" {
		var dto http_dto.GetRoleByNameDTO = http_dto.GetRoleByNameDTO{
			Name: ctx.Query("name"),
		}

		// TODO validate

		// Mapping http DTO -> domain DTO
		// domainDTO := domain_dto.GetRoleByNameDTO{
		// 	Name: dto.Name,
		// }

		var err error
		role, err = h.usecase.GetRoleByName(context.Background(), dto.Name)
		if err != nil {
			if errors.Is(err, usecase.ErrNotFound) {
				ctx.Status(fiber.StatusNoContent)
				return nil
			}

			log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the role by name")
			return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
				fiber.StatusInternalServerError,
				"An unexpected error has occurred while trying to retrieve the role by name",
				"",
				nil,
			))
		}
	} else {
		log.Warn().Msg("invalid request body: either the 'id' or 'name' parameter was not provided")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Either the 'id' or 'name' parameter was not provided",
			nil,
		))
	}

	roleDTO := http_dto.RoleDTO{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: role.DeletedAt,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": roleDTO})
}
