package v1

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	// 	http_dto "location-backend/internal/controller/http/dto"
	// 	"location-backend/internal/controller/http/mapper"
	// 	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/usecase"
	"location-backend/pkg/httperrors"
)

const (
	createMatrixURL = "/"
	findPointsURL   = "/find"
)

type matrixHandler struct {
	floorUsecase  *usecase.FloorUsecase
	matrixUsecase *usecase.MatrixUsecase
	// sensorMapper *mapper.SensorMapper
	// aprtMapper *mapper.SensorRadioTemplateMapper
}

// Регистрирует новый handler
func NewMatrixHandler(floorUsecase *usecase.FloorUsecase, matrixUsecase *usecase.MatrixUsecase) *matrixHandler {
	return &matrixHandler{
		floorUsecase:  floorUsecase,
		matrixUsecase: matrixUsecase,
		// sensorMapper: &mapper.SensorMapper{},
		// aprtMapper: &mapper.SensorRadioTemplateMapper{},
	}
}

// Регистрирует маршруты для user
func (h *matrixHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	// Create
	router.Post(createSensorURL, h.CreateMatrix)

	return router
}

func (h *matrixHandler) CreateMatrix(ctx *fiber.Ctx) (err error) { // TODO перенести часть логики в usecase
	floorID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse 'id' as UUID")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid ID",
			"Failed to parse 'id' as UUID",
			nil,
		))
	}

	// TODO вернуть floor
	_, err = h.floorUsecase.GetFloor(context.Background(), floorID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			log.Info().Err(err).Msg("the floor with provided 'id' does not exist")
			return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
				fiber.StatusBadRequest,
				"Invalid request body",
				"The floor with provided 'id' does not exist",
				nil,
			))
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to retrieve the floor")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to retrieve the floor",
			"",
			nil,
		))
	}

	err = h.matrixUsecase.CreateMatrix(context.Background(), floorID)
	if err != nil {
		log.Error().Err(err).Msg("an unexpected error has occurred while trying to create matrix")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to create matrix",
			"",
			nil,
		))
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *matrixHandler) FindPoints(ctx *fiber.Ctx) (err error) {
	var domainDTO dto.FindPointsDTO
	err = ctx.BodyParser(&domainDTO)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse find points")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse find points",
			nil,
		))
	}

	// TODO validate

	points, err := h.matrixUsecase.FindPoints(context.Background(), &domainDTO)
	if err != nil {
		// if errors.Is(err, usecase.ErrNotFound) {
		// 	log.Info().Err(err).Msg("the floor with provided 'floor_id' does not exist")
		// 	return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
		// 		fiber.StatusBadRequest,
		// 		"Invalid request body",
		// 		"The floor with provided 'floor_id' does not exist",
		// 		nil,
		// 	))
		// }

		// if errors.Is(err, usecase.ErrAlreadyExists) {
		// 	log.Info().Err(err).Msg("the sensor with provided 'mac' already exists")
		// 	return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
		// 		fiber.StatusBadRequest,
		// 		"Invalid request body",
		// 		"The sensor with provided 'mac' already exists",
		// 		nil,
		// 	))
		// }

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to find points")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to find points",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"points": points})
}
