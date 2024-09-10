package service

import (
	"context"

	"github.com/google/uuid"

	"location-backend/internal/domain/entity"
)

type UserRepo interface {
	GetOne(id uuid.UUID) entity.User
	GetOneByName(username string) entity.User
	GetAll(limit, offset int) []entity.User
	Create(book entity.User) entity.User
	// Delete(book entity.Book) error
}

type userService struct {
	repository UserRepo
}

func NewUserService(repository UserRepo) *userService {
	return &userService{repository: repository}
}

func (s userService) Create(ctx context.Context) entity.User {
	return entity.User{}
}

func (s userService) GetByID(ctx context.Context, id uuid.UUID) entity.User {
	return s.repository.GetOne(id)
}

func (s userService) GetByUsername(ctx context.Context, username string) entity.User {
	return s.repository.GetOneByName(username)
}

// func (s userService) GetAll(ctx context.Context, limit, offset int) []entity.Book {
// 	return s.repository.GetAll(limit, offset)
// }

// func (s userService) GetAllForList(ctx context.Context) []entity.BookView {
// 	// TODO implement
// 	return nil
// }
