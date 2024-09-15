package v1

import (
	// "encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"location-backend/internal/controller/http/dto"
	user_usecase "location-backend/internal/domain/usecase/user"
	"location-backend/internal/server"
)

const (
	// bookURL  = "/users/:user_id"
	// booksURL = "/users"
	registerURL = "/register"
	loginURL    = "/login"
)

//? Здесь был интерфейс UserUsercase из бизнес логики

type userHandler struct {
	userUsecase user_usecase.UserUsecase
}

func NewUserHandler(userUsecase user_usecase.UserUsecase) *userHandler {
	return &userHandler{userUsecase: userUsecase}
}

func (h *userHandler) Register(router *server.Router) {
	router.App.Post(registerURL, h.CreateUser)
}

// func (h *bookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	// books := h.bookService.GetAll(context.Background(), 0, 0)
// 	w.Write([]byte("books"))
// 	w.WriteHeader(http.StatusOK)
// }

func (h *userHandler) CreateUser(ctx *fiber.Ctx) error {
	// DTO from client (HTTP/JSON)
	var d dto.CreateUserDTO
	err := ctx.BodyParser(&d)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse user request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// TODO validate

	// Mapping dto.CreateUserDTO --> user_usecase.CreateUserDTO
	usecaseDTO := user_usecase.CreateUserDTO{
		Username: d.Username,
		Password: d.Password,
	}

	// Call the use case to create the user
	userID, err := h.userUsecase.CreateUser(ctx, usecaseDTO)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		//? JSON RPC: TRANSPORT: 200, error: {msg, ..., dev_msg}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(user))

	return ctx.Status(fiber.StatusOK).JSON(userID)
	// return nil
}
