package usecase

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"location-backend/internal/config"

	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/service"
)

var (

	// Occurs when login with wrong login or password, or if user does not exist
	ErrBadLogin = errors.New("incorrect login or password, or no such user")
)

//? Здесь был интерфейс сервиса (Перенесён в в сервисы)

type UserUsecase interface {
	Register(dto http_dto.RegisterUserDTO) (userID uuid.UUID, err error)
	Login(dto http_dto.LoginUserDTO) (signedString string, err error)
	GetUserByName(dto http_dto.GetUserByNameDTO) (user entity.User, err error)
	// ListAllBooks(ctx context.Context) []entity.BookView
	// GetFullBook(ctx context.Context, id string) entity.FullBook
}

// type AuthorService interface {
// 	GetByID(ctx context.Context, id string) entity.User
// }

// type GenreService interface {
// 	GetByID(ctx context.Context, id string) entity.User
// }

type userUsecase struct {
	userService service.UserService
	// authorService UserService
	// genreService  GenreService
}

// ? TEST. Изначально этого не было
func NewUserUsecase(userService service.UserService) *userUsecase {
	return &userUsecase{userService: userService}
}

// ? Нужен ли ctx *fiber.Ctx
func (u userUsecase) Register(dto http_dto.RegisterUserDTO) (userID uuid.UUID, err error) {
	_, err = u.userService.GetUserByName(dto.Username)
	if err != nil {
		// If error except ErrNotFound
		if !errors.Is(err, service.ErrNotFound) {
			log.Error().Err(err).Msg("failed to check user existing")
			return
		}
	} else { // If user already exists
		return userID, ErrAlreadyExists
	}

	hash, err := hashPassword(dto.Password)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash user password")
		return
	}

	var createUserDTO domain_dto.CreateUserDTO = domain_dto.CreateUserDTO{
		Username:     dto.Username,
		PasswordHash: hash,
	}

	userID, err = u.userService.CreateUser(createUserDTO)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return
	}

	log.Info().Msgf("user %v successfully registered", dto.Username)

	return
}

func (u userUsecase) Login(dto http_dto.LoginUserDTO) (signedString string, err error) {
	user, err := u.userService.GetUserByName(dto.Username)
	if err != nil {
		// Return ErrBadLogin if user not found
		if errors.Is(err, service.ErrNotFound) {
			return "", ErrBadLogin
		} else {
			log.Error().Err(err).Msg("failed to check user existing")
			return
		}
	}

	if !checkPasswordHash(dto.Password, user.PasswordHash) {
		log.Info().Msg("wrong password")
		return "", ErrBadLogin
	}

	// TODO already login err

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
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

func (u userUsecase) GetUserByName(dto http_dto.GetUserByNameDTO) (user entity.User, err error) {
	user, err = u.userService.GetUserByName(dto.Username)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return entity.User{}, ErrNotFound
		} else {
			log.Error().Err(err).Msg("failed to get user")
			return
		}
	}

	return
}

// Хэширует пароль
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	// log.Debug().Msgf("Password: %v", password)
	// log.Debug().Msgf("HashPassword: %v", bytes)
	return string(bytes), err
}

// Сравнивает пароль и его хэш. Если верно – true, иначе – false.
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Debug().Msgf("failed to compare hash and password (password: '%v' \t hash: '%v')", password, hash)
	}

	return err == nil
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
