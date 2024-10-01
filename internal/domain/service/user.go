package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/usecase"
)

type UserRepo interface {
	Create(ctx context.Context, dto *dto.CreateUserDTO) (userID uuid.UUID, err error)
	GetOneByName(ctx context.Context, username string) (user *entity.User, err error)
}

type userService struct {
	repository UserRepo
}

func NewUserService(repository UserRepo) *userService {
	return &userService{repository: repository}
}

func (s userService) CreateUser(ctx context.Context, createUserDTO *dto.CreateUserDTO) (userID uuid.UUID, err error) {
	userID, err = s.repository.Create(ctx, createUserDTO)
	return
}

// ? Нужен ли ctx *fiber.Ctx здесь?
func (s userService) GetUserByName(ctx context.Context, username string) (user *entity.User, err error) {
	user, err = s.repository.GetOneByName(ctx, username)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return user, usecase.ErrNotFound
		}

		return
	}

	return
}
