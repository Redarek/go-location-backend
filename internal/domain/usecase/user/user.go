package user_usecase

import (
	"context"

	"github.com/google/uuid"

	"location-backend/internal/controller/http/dto"
	"location-backend/internal/domain/entity"
)

type Service interface {
	// GetAllForList(ctx context.Context) []entity.BookView
	GetByID(ctx context.Context, id uuid.UUID) entity.User
	GetByUsername(ctx context.Context, username string) entity.User
}

// type AuthorService interface {
// 	GetByID(ctx context.Context, id string) entity.User
// }

// type GenreService interface {
// 	GetByID(ctx context.Context, id string) entity.User
// }

type userUsecase struct {
	userService Service
	// authorService UserService
	// genreService  GenreService
}

func (u userUsecase) CreateBook(ctx context.Context, dto dto.CreateUserDTO) (string, error) {
	return "", nil
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
