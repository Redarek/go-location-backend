package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"location-backend/internal/adapters/db/repository"
	"location-backend/internal/controller/http/dto"
	"location-backend/internal/domain/entity"
)

//? Здесь был интерфейсы репозитория UserRepo (перенесён в репозиторий)

type UserService interface {
	// GetAllForList(ctx context.Context) []entity.BookView
	// GetByID(ctx context.Context, id uuid.UUID) entity.User
	// GetByUsername(ctx context.Context, username string) entity.User
	CreateUser(ctx *fiber.Ctx, dto dto.CreateUserDTO) (userID uuid.UUID, err error)
}

type userService struct {
	repository repository.UserRepo
}

func NewUserService(repository repository.UserRepo) *userService {
	return &userService{repository: repository}
}

func (s userService) Create(ctx context.Context) entity.User {
	return entity.User{}
}

func (s userService) GetByID(ctx context.Context, id uuid.UUID) entity.User {
	return s.repository.GetOne(id)
}

// func (s userService) GetByUsername(ctx context.Context, username string) entity.User {
// 	return s.repository.GetOneByName(username)
// }

// func (s userService) GetAll(ctx context.Context, limit, offset int) []entity.Book {
// 	return s.repository.GetAll(limit, offset)
// }

// func (s userService) GetAllForList(ctx context.Context) []entity.BookView {
// 	// TODO implement
// 	return nil
// }
