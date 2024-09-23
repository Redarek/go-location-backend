package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	repository "location-backend/internal/adapters/db/postgres"
	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

//? Здесь был интерфейсы репозитория UserRepo (перенесён в репозиторий)

type UserService interface {
	// GetAllForList(ctx context.Context) []entity.BookView
	// GetByID(ctx context.Context, id uuid.UUID) entity.User
	GetUserByName(ctx context.Context, username string) (user *entity.User, err error)
	CreateUser(ctx context.Context, userCreate dto.CreateUserDTO) (userID uuid.UUID, err error)

	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type userService struct {
	repository repository.UserRepo
}

func NewUserService(repository repository.UserRepo) *userService {
	return &userService{repository: repository}
}

func (s userService) CreateUser(ctx context.Context, userCreate dto.CreateUserDTO) (userID uuid.UUID, err error) {
	userID, err = s.repository.Create(ctx, userCreate)
	if err != nil {
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to create user")
		return
	}

	return userID, nil
}

// ? Нужен ли ctx *fiber.Ctx здесь?
func (s userService) GetUserByName(ctx context.Context, username string) (user *entity.User, err error) {
	user, err = s.repository.GetOneByName(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return user, ErrNotFound
		}
		// TODO улучшить лог
		log.Error().Err(err).Msg("failed to retrieve user")
		return
	}

	return
}

// Хэширует пароль
func (s userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	// log.Debug().Msgf("Password: %v", password)
	// log.Debug().Msgf("HashPassword: %v", bytes)
	return string(bytes), err
}

// Сравнивает пароль и его хэш. Если верно – true, иначе – false.
func (s userService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Debug().Msgf("failed to compare hash and password (password: '%v' \t hash: '%v')", password, hash)
	}

	return err == nil
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
