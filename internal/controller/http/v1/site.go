package v1

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/usecase"
	"location-backend/pkg/httperrors"
)

const (
	getSiteURL    = "/"
	createSiteURL = "/"
)

type siteHandler struct {
	usecase usecase.SiteUsecase
}

// Регистрирует новый handler
func NewSiteHandler(usecase usecase.SiteUsecase) *siteHandler {
	return &siteHandler{usecase: usecase}
}

// Регистрирует маршруты для user
func (h *siteHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createSiteURL, h.CreateSite)
	router.Get(getSiteURL, h.GetSite)

	return router
}

func (h *siteHandler) CreateSite(ctx *fiber.Ctx) error {
	var dto http_dto.CreateSiteDTO
	err := ctx.BodyParser(&dto)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse site request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse site request body",
			nil,
		))
	}

	// TODO validate

	// Получение ID пользователя из JWT
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse user ID from JWT")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Unable to retrieve user information due to an issue with your authentication token",
			"",
			nil,
		))
	}

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.CreateSiteDTO{
		Name:        dto.Name,
		Description: dto.Description,
		UserID:      userID,
	}

	siteID, err := h.usecase.CreateSite(domainDTO)
	if err != nil {
		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create the site")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create the site",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": siteID})
}

func (h *siteHandler) GetSite(ctx *fiber.Ctx) error {
	siteID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	var dto http_dto.GetSiteDTO = http_dto.GetSiteDTO{
		ID: siteID,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.GetSiteDTO{
		ID: dto.ID,
	}

	site, err := h.usecase.GetSite(domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the site")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the site",
			"",
			nil,
		))
	}

	// Mapping domain DTO -> http DTO
	siteDTO := http_dto.SiteDTO{
		ID:          site.ID,
		Name:        site.Name,
		Description: site.Description,
		UserID:      site.UserID,
		CreatedAt:   site.CreatedAt,
		UpdatedAt:   site.UpdatedAt,
		DeletedAt:   site.DeletedAt,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": siteDTO})
}
