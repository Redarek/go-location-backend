package user_usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"location-backend/internal/controller/http/dto"
	"location-backend/internal/domain/service"
)

//? Здесь был интерфейс сервиса (Перенесён в в сервисы)

type UserUsecase interface {
	CreateUser(ctx *fiber.Ctx, dto CreateUserDTO) (userID uuid.UUID, err error)
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

func (u userUsecase) CreateUser(ctx *fiber.Ctx, dto dto.CreateUserDTO) (userID uuid.UUID, err error) {
	return userID, nil
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
