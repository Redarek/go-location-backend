package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"location-backend/internal/controller/http/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/service"
)

// ErrAlreadyRegistered occurs when register the already registered user
var ErrAlreadyRegistered = errors.New("user is already registered")

//? Здесь был интерфейс сервиса (Перенесён в в сервисы)

type UserUsecase interface {
	Register(dto dto.CreateUserDTO) (userID uuid.UUID, err error)
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

type createUserDTO struct {
	Username     string
	PasswordHash string
}

// ? TEST. Изначально этого не было
func NewUserUsecase(userService service.UserService) *userUsecase {
	return &userUsecase{userService: userService}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	// log.Debug().Msgf("Password: %v", password)
	// log.Debug().Msgf("HashPassword: %v", bytes)
	return string(bytes), err
}

// ? Нужен ли ctx *fiber.Ctx
func (u userUsecase) Register(dto dto.CreateUserDTO) (userID uuid.UUID, err error) {
	_, err = u.userService.GetUserByName(dto.Username)
	if err != nil {
		// If error except ErrNotFound
		if !errors.Is(err, service.ErrNotFound) {
			log.Error().Err(err).Msg("Failed to check user existing")
			return
		}
	} else { // If user already exists
		return userID, ErrAlreadyRegistered
	}

	hash, err := hashPassword(dto.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash user password")
		return
	}

	var userCreate entity.UserCreate = entity.UserCreate{
		Username:     dto.Username,
		PasswordHash: hash,
	}

	userID, err = u.userService.CreateUser(userCreate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return
	}

	log.Info().Msgf("User %v successfully registered", dto.Username)

	return
	// return ctx.JSON(fiber.Map{"id": userID})
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
