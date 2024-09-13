package v1

import (
	// "encoding/json"

	user_usecase "location-backend/internal/domain/usecase/user"
	// "location-backend/internal/controller/http/dto"
	"location-backend/internal/server"

	"github.com/gofiber/fiber/v2"
	// "github.com/julienschmidt/httprouter"
)

const (
	// bookURL  = "/users/:user_id"
	// booksURL = "/users"
	registerURL = "/register"
	loginURL    = "/login"
)

type UserUsecase interface {
	CreateBook(ctx *fiber.Ctx, dto user_usecase.CreateUserDTO) (string, error)
	// ListAllBooks(ctx context.Context) []entity.BookView
	// GetFullBook(ctx context.Context, id string) entity.FullBook
}

type userHandler struct {
	userUsecase UserUsecase
}

func NewUserHandler(userUsecase UserUsecase) *userHandler {
	return &userHandler{userUsecase: userUsecase}
}

func (h *userHandler) Register(router *server.Fiber) {
	// router.GET(booksURL, h.GetAllBooks)
}

// func (h *bookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	// books := h.bookService.GetAll(context.Background(), 0, 0)
// 	w.Write([]byte("books"))
// 	w.WriteHeader(http.StatusOK)
// }

// func (h *userHandler) CreateBook(ctx *fiber.Ctx, dto user_usecase.CreateUserDTO) {

// 	var d dto.CreateUserDTO
// 	defer r.Body.Close()
// 	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
// 		return // error
// 	}

// 	// validate

// 	// MAPPING dto.CreateBookDTO --> book_usecase.CreateBookDTO
// 	usecaseDTO := book_usecase.CreateBookDTO{
// 		Name:       "",
// 		Year:       0,
// 		AuthorUUID: "",
// 		GenreUUID:  "",
// 	}
// 	book, err := h.userUsecase.CreateBook(r.Context(), usecaseDTO)
// 	if err != nil {
// 		// JSON RPC: TRANSPORT: 200, error: {msg, ..., dev_msg}
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(book))
// }
