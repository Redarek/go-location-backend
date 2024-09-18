package composites

import (
	repository "location-backend/internal/adapters/db/postgres"
	handler "location-backend/internal/controller/http/v1"
	"location-backend/internal/domain/service"
	"location-backend/internal/domain/usecase"
)

type HealthComposite struct {
	Repository repository.HealthRepo
	Service    service.HealthService
	Usecase    usecase.HealthUsecase
	Handler    handler.Handler
}

func NewHealthComposite(composite *PostgresComposite) (userComposite *HealthComposite) {
	healthRepo := repository.NewHealthRepo(composite.pool)
	healthService := service.NewHealthService(healthRepo)
	healthUsecase := usecase.NewHealthUsecase(healthService)
	healthHandler := handler.NewHealthHandler(healthUsecase)

	return &HealthComposite{
		Repository: healthRepo,
		Service:    healthService,
		Usecase:    healthUsecase,
		Handler:    healthHandler,
	}
}
