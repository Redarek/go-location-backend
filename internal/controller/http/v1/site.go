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
	createSiteURL = "/"
	getSiteURL    = "/"
	getSitesURL   = "/all"

	patchUpdateSiteURL = "/"

	softDeleteSiteURL = "/sd"
	restoreSiteURL    = "/restore"
)

type siteHandler struct {
	usecase *usecase.SiteUsecase
}

// Регистрирует новый handler
func NewSiteHandler(usecase *usecase.SiteUsecase) *siteHandler {
	return &siteHandler{usecase: usecase}
}

// Регистрирует маршруты для user
func (h *siteHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createSiteURL, h.CreateSite)
	router.Get(getSiteURL, h.GetSite)
	router.Get(getSitesURL, h.GetSites)

	router.Patch(patchUpdateSiteURL, h.PatchUpdateSite)

	router.Patch(softDeleteSiteURL, h.SoftDeleteSite)
	router.Patch(restoreSiteURL, h.RestoreSite)

	// TODO Get list detailed
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
	userID, err := GetUserIDFromJWT(ctx)
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
	domainDTO := &domain_dto.CreateSiteDTO{
		Name:        dto.Name,
		Description: dto.Description,
		UserID:      userID,
	}

	siteID, err := h.usecase.CreateSite(context.Background(), domainDTO)
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
	// domainDTO := domain_dto.GetSiteDTO{
	// 	ID: dto.ID,
	// }

	site, err := h.usecase.GetSite(context.Background(), dto.ID)
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

func (h *siteHandler) GetSites(c *fiber.Ctx) error {
	// Получение ID пользователя из JWT
	userID, err := GetUserIDFromJWT(c)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse user ID from JWT")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Unable to retrieve user information due to an issue with your authentication token",
			"",
			nil,
		))
	}

	// TODO реализовать передачу page и size
	var dto http_dto.GetSitesDTO = http_dto.GetSitesDTO{
		UserID: userID,
		Page:   1,
		Size:   100,
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.GetSitesDTO{
		UserID: dto.UserID,
		Limit:  dto.Size,
		Offset: (dto.Page - 1) * dto.Size,
	}

	sites, err := h.usecase.GetSites(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the site")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the site",
			"",
			nil,
		))
	}

	var sitesDTO []http_dto.SiteDTO
	for _, site := range sites {
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

		sitesDTO = append(sitesDTO, siteDTO)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": sitesDTO})
}

func (h *siteHandler) PatchUpdateSite(c *fiber.Ctx) error {
	var dto http_dto.PatchUpdateSiteDTO
	err := c.BodyParser(&dto)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse site request body")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse site request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := &domain_dto.PatchUpdateSiteDTO{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
	}

	err = h.usecase.PatchUpdateSite(context.Background(), domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The site was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrNotUpdated) {
			c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"The site was not updated",
				"It usually occurs when there are no editable fields provided",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to patch update the site")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to patch update the site",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *siteHandler) SoftDeleteSite(c *fiber.Ctx) error {
	siteID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.SoftDeleteSiteDTO{
	// 	ID: siteID,
	// }

	err = h.usecase.SoftDeleteSite(context.Background(), siteID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The site was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadySoftDeleted) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The site has already soft deleted",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to soft delete the site")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to soft delete the site",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *siteHandler) RestoreSite(c *fiber.Ctx) error {
	siteID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return c.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	// domainDTO := domain_dto.SoftDeleteSiteDTO{
	// 	ID: siteID,
	// }

	err = h.usecase.RestoreSite(context.Background(), siteID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.Status(fiber.StatusNotFound).JSON(httperrors.NewErrorResponse(
				fiber.StatusNotFound,
				"The site was not found",
				"",
				nil,
			))
			return nil
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			c.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"The site is already restored",
				"",
				nil,
			))
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to restore the site")
		return c.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to restore the site",
			"",
			nil,
		))
	}

	c.Status(fiber.StatusOK)
	return nil
}
