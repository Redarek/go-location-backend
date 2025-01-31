package v1

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/usecase"
	"location-backend/pkg/httperrors"
)

const (
	createSiteURL       = "/"
	getSiteURL          = "/"
	getSitesURL         = "/all"
	getSitesDetailedURL = "/all/detailed"

	patchUpdateSiteURL = "/"

	softDeleteSiteURL = "/sd"
	restoreSiteURL    = "/restore"
)

type siteHandler struct {
	usecase *usecase.SiteUsecase
}

// Регистрирует новый handler
func NewSiteHandler(usecase *usecase.SiteUsecase) *siteHandler {
	return &siteHandler{
		usecase: usecase,
	}
}

// Регистрирует маршруты для user
func (h *siteHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Post(createSiteURL, h.CreateSite)
	router.Get(getSiteURL, h.GetSite)
	router.Get(getSitesURL, h.GetSites)
	router.Get(getSitesDetailedURL, h.GetSitesDetailed)

	router.Patch(patchUpdateSiteURL, h.PatchUpdateSite)

	router.Patch(softDeleteSiteURL, h.SoftDeleteSite)
	router.Patch(restoreSiteURL, h.RestoreSite)

	// TODO Get list detailed
	return router
}

func (h *siteHandler) CreateSite(ctx *fiber.Ctx) error {
	// UserID не передаётся!
	var dtoObj dto.CreateSiteDTO
	err := ctx.BodyParser(&dtoObj)
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

	dtoObj.UserID = userID

	siteID, err := h.usecase.CreateSite(context.Background(), &dtoObj)
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

	var dtoObj dto.GetSiteDTO = dto.GetSiteDTO{
		ID: siteID,
	}

	// TODO validate

	site, err := h.usecase.GetSite(context.Background(), dtoObj.ID)
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

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": site})
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
	var dtoObj dto.GetSitesDTO = dto.GetSitesDTO{
		UserID: userID,
		Page:   1,
		Size:   100,
	}

	// TODO validate

	sites, err := h.usecase.GetSites(context.Background(), &dtoObj)
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": sites})
}

func (h *siteHandler) GetSitesDetailed(c *fiber.Ctx) error {
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
	var dtoObj dto.GetSitesDetailedDTO = dto.GetSitesDetailedDTO{
		UserID: userID,
		Page:   1,
		Size:   100,
	}

	// TODO validate

	sitesDetailed, err := h.usecase.GetSitesDetailed(context.Background(), &dtoObj)
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": sitesDetailed})
}

func (h *siteHandler) PatchUpdateSite(c *fiber.Ctx) error {
	var dtoObj dto.PatchUpdateSiteDTO
	err := c.BodyParser(&dtoObj)
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

	err = h.usecase.PatchUpdateSite(context.Background(), &dtoObj)
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
