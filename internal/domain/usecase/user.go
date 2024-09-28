package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"location-backend/internal/config"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type UserService interface {
	// GetAllForList(ctx context.Context) []entity.BookView
	// GetByID(ctx context.Context, id uuid.UUID) entity.User
	GetUserByName(ctx context.Context, username string) (user *entity.User, err error)
	CreateUser(ctx context.Context, createUserDTO *domain_dto.CreateUserDTO) (userID uuid.UUID, err error)

	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type UserUsecase struct {
	userService UserService
	// authorService UserService
	// genreService  GenreService
}

// ? TEST. Изначально этого не было
func NewUserUsecase(userService UserService) *UserUsecase {
	return &UserUsecase{userService: userService}
}

// ? Нужен ли ctx *fiber.Ctx
func (u *UserUsecase) Register(ctx context.Context, dto *domain_dto.RegisterUserDTO) (userID uuid.UUID, err error) {
	_, err = u.userService.GetUserByName(ctx, dto.Username)
	if err != nil {
		// If error except ErrNotFound
		if !errors.Is(err, ErrNotFound) {
			log.Error().Err(err).Msg("failed to check user existing")
			return
		}
	} else { // If user already exists
		return userID, ErrAlreadyExists
	}

	hash, err := u.userService.HashPassword(dto.Password)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash user password")
		return
	}

	var createUserDTO domain_dto.CreateUserDTO = domain_dto.CreateUserDTO{
		Username:     dto.Username,
		PasswordHash: hash,
	}

	userID, err = u.userService.CreateUser(ctx, &createUserDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return
	}

	log.Info().Msgf("user %v successfully registered", dto.Username)

	return
}

func (u *UserUsecase) Login(ctx context.Context, dto *domain_dto.LoginUserDTO) (signedString string, err error) {
	user, err := u.userService.GetUserByName(ctx, dto.Username)
	if err != nil {
		// Return ErrBadLogin if user not found
		if errors.Is(err, ErrNotFound) {
			return "", ErrBadLogin
		} else {
			log.Error().Err(err).Msg("failed to check user existing")
			return
		}
	}

	if !u.userService.CheckPasswordHash(dto.Password, user.PasswordHash) {
		log.Info().Msg("wrong password")
		return "", ErrBadLogin
	}

	// TODO already login err

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // TODO вынести в конфиг
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	signedString, err = token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		log.Error().Err(err).Msg("failed to sign token")
		return
		// return c.SendStatus(fiber.StatusInternalServerError)
	}

	return
	// return c.JSON(fiber.Map{"token": signedString})
}

func (u *UserUsecase) GetUserByName(ctx context.Context, username string) (user *entity.User, err error) {
	user, err = u.userService.GetUserByName(ctx, username)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get user")
			return
		}
	}

	return
}

// func (u userUsecase) ListAllBooks(ctx context.Context) []entity.BookView {
// 	// отобразить список книг с именем Жанра и именем Автора
// 	return u.userService.GetAllForList(ctx)
// }

// func (u bookUsecase) GetFullBook(ctx context.Context, id string) entity.FullBook {
// 	book := u.bookService.GetByID(ctx, id)
// 	author := u.authorService.GetByID(ctx, book.AuthorID)
// 	genre := u.genreService.GetByID(ctx, book.GenreID)

// 	return entity.FullBook{
// 		Book:   book,
// 		Author: author,
// 		Genre:  genre,
// 	}
// }

// // pagination
// func (u bookUsecase) GetBooksWithAllAuthors(ctx context.Context, id string) []entity.BookView {
// 	// Book{Authors: [all authors]}
// 	// book, author(book_id) -=-
// 	return nil
// }
