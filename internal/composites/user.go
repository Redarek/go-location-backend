package composites

import (
	"location-backend/internal/adapters/db/repository"
	handler "location-backend/internal/controller/http/v1"
	"location-backend/internal/domain/service"
	"location-backend/internal/domain/usecase"
)

type UserComposite struct {
	Repository repository.UserRepo
	Service    service.UserService
	Usecase    usecase.UserUsecase
	Handler    handler.Handler
}

func NewUserComposite(composite *PostgresComposite) (userComposite *UserComposite, err error) {
	userRepo := repository.NewUserRepo(composite.pool)
	userService := service.NewUserService(userRepo)
	userUsecase := usecase.NewUserUsecase(userService)
	userHandler := handler.NewUserHandler(userUsecase)

	return &UserComposite{
		Repository: userRepo,
		Service:    userService,
		Usecase:    userUsecase,
		Handler:    userHandler,
	}, nil
}
