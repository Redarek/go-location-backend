package service

import (
	// "context"

	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	// "github.com/google/uuid"

	repository "location-backend/internal/adapters/db/postgres"
	// "location-backend/internal/controller/http/dto"
	"location-backend/internal/domain/entity"
)

//? Здесь был интерфейсы репозитория UserRepo (перенесён в репозиторий)

type UserService interface {
	// GetAllForList(ctx context.Context) []entity.BookView
	// GetByID(ctx context.Context, id uuid.UUID) entity.User
	GetUserByName(username string) (user entity.User, err error)
	CreateUser(userCreate entity.UserCreate) (userID uuid.UUID, err error)
}

type userService struct {
	repository repository.UserRepo
}

func NewUserService(repository repository.UserRepo) *userService {
	return &userService{repository: repository}
}

func (s userService) CreateUser(userCreate entity.UserCreate) (userID uuid.UUID, err error) {
	userID, err = s.repository.Create(userCreate)
	if err != nil {
		// TODO улучшить лог
		log.Error().Err(err).Msg("Failed to create user")
		return
	}

	return userID, nil
}

// ? Нужен ли ctx *fiber.Ctx здесь?
func (s userService) GetUserByName(username string) (user entity.User, err error) {
	user, err = s.repository.GetOneByName(username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return user, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("Failed to retrieve user")
		return
	}

	return
}

// func (s userService) GetByID(ctx context.Context, id uuid.UUID) entity.User {
// 	return s.repository.GetOne(id)
// }

// func (s userService) GetAll(ctx context.Context, limit, offset int) []entity.Book {
// 	return s.repository.GetAll(limit, offset)
// }

// func (s userService) GetAllForList(ctx context.Context) []entity.BookView {
// 	// TODO implement
// 	return nil
// }
