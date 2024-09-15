package repository

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"location-backend/internal/domain/entity"
)

type UserRepo interface {
	Create(ctx *fiber.Ctx) entity.User
	GetOne(id uuid.UUID) entity.User
	// GetOneByName(username string) entity.User
	// GetAll(limit, offset int) []entity.User
	// Delete(book entity.Book) error
}

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *userRepo {
	return &userRepo{pool: pool}
}

func (bs *userRepo) Create(ctx *fiber.Ctx) entity.User {
	return entity.User{}
}

func (bs *userRepo) GetOne(id uuid.UUID) entity.User {
	// var entity.User user
	return entity.User{}
}

//	func (bs *userRepo) GetAll(limit, offset int) []*entity.User {
//		return nil
//	}

// func (bs *userRepo) Delete(user *entity.User) error {
// 	return nil
// }
