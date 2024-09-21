package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/usecase"
)

const (
	healthURL = "/health"
)

// type HealthHandler interface {
// 	Health(ctx *fiber.Ctx) error
// }

type healthHandler struct {
	healthUsecase usecase.HealthUsecase
}

func NewHealthHandler(healthUsecase usecase.HealthUsecase) *healthHandler {
	return &healthHandler{healthUsecase: healthUsecase}
}

func (h *healthHandler) Register(r *fiber.Router) fiber.Router {
	router := *r
	router.Get(healthURL, h.Health)
	return router
}

func (h healthHandler) Health(ctx *fiber.Ctx) error {
	err := h.healthUsecase.Health()
	if err != nil {
		log.Error().Err(err).Msg("failed to check database health")
		// ? Точно ли статус 500?
		return ctx.Status(fiber.StatusInternalServerError).SendString("It's not healthy")
	}

	return ctx.Status(fiber.StatusOK).SendString("It's healthy")
}
