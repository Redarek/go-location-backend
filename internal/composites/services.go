package composites

import (
	"location-backend/internal/domain/service"
)

type ServiceComposite struct {
	healthService service.HealthService
	userService   service.UserService
}

func NewServiceComposite(composite *RepositoryComposite) (serviceComposite *ServiceComposite) {
	return &ServiceComposite{
		healthService: service.NewHealthService(composite.healthRepo),
		userService:   service.NewUserService(composite.userRepo),
	}
}
