package v1

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
	"location-backend/pkg/httperrors"
)

const (
	getRoleByNameURL = "/"
	createRoleURL    = "/"
)

type roleHandler struct {
	usecase *usecase.RoleUsecase
}

// Регистрирует новый handler
func NewRoleHandler(usecase *usecase.RoleUsecase) *roleHandler {
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
	var dtoObj dto.CreateRoleDTO
	err := ctx.BodyParser(&dtoObj)
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

	roleID, err := h.usecase.CreateRole(context.Background(), &dtoObj)
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
	var role *entity.Role

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

		// TODO validate

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
		var dtoObj dto.GetRoleByNameDTO = dto.GetRoleByNameDTO{
			Name: ctx.Query("name"),
		}

		// TODO validate

		var err error
		role, err = h.usecase.GetRoleByName(context.Background(), dtoObj.Name)
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

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": role})
}
