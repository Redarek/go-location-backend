package composites

import (
	"location-backend/internal/domain/usecase"
)

type UsecaseComposite struct {
	healthUsecase usecase.HealthUsecase
	userUsecase   usecase.UserUsecase
	roleUsecase   usecase.RoleUsecase
}

func NewUsecaseComposite(composite *ServiceComposite) (serviceComposite *UsecaseComposite) {
	return &UsecaseComposite{
		healthUsecase: usecase.NewHealthUsecase(composite.healthService),
		userUsecase:   usecase.NewUserUsecase(composite.userService),
		roleUsecase:   usecase.NewRoleUsecase(composite.roleService),
	}
}
