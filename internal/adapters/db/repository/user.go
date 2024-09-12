package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"location-backend/internal/domain/entity"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{db: db}
}

func (bs *userRepo) GetOne(id string) *entity.User {
	return nil
}
func (bs *userRepo) GetAll(limit, offset int) []*entity.User {
	return nil
}
func (bs *userRepo) Create(user *entity.User) *entity.User {
	return nil
}
func (bs *userRepo) Delete(user *entity.User) error {
	return nil
}
