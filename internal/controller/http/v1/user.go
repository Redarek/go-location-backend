package v1

import (
	// "encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"location-backend/internal/controller/http/dto"
	"location-backend/internal/domain/usecase"
	"location-backend/internal/router"
)

const (
	// bookURL  = "/users/:user_id"
	// booksURL = "/users"
	registerURL = "/api/v1/user/register"
	loginURL    = "/api/v1/user/login"
)

//? Здесь был интерфейс UserUsercase из бизнес логики

// ? Раньше не было
// type UserHandler interface {
// 	CreateUser(ctx *fiber.Ctx) error
// }

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *userHandler {
	return &userHandler{userUsecase: userUsecase}
}

func (h *userHandler) Register(router *router.Router) {
	router.App.Post(registerURL, h.RegisterUser)
}

// func (h *bookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	// books := h.bookService.GetAll(context.Background(), 0, 0)
// 	w.Write([]byte("books"))
// 	w.WriteHeader(http.StatusOK)
// }

// RegisterUser регистрирует нового пользователя, если его не существует.
//
// Возвращаемые статусы:
//
//	201 Created – пользователь успешно создан
//	409 Conflict – пользователь уже существует
//	500 InternalServerError – ошибка сервера
func (h *userHandler) RegisterUser(ctx *fiber.Ctx) error {
	// DTO from client (HTTP/JSON)
	var dto dto.CreateUserDTO
	err := ctx.BodyParser(&dto)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse user request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// TODO validate

	// ? use only if more than 1 controller
	// // Mapping dto.CreateUserDTO --> user_usecase.CreateUserDTO
	// usecaseDTO := user_usecase.CreateUserDTO{
	// 	Username: d.Username,
	// 	Password: d.Password,
	// }

	// ? Нужно ли передавать ctx внутрь?
	// Call the use case to create the user
	userID, err := h.userUsecase.Register(dto)
	if err != nil {
		if errors.Is(err, usecase.ErrAlreadyRegistered) {
			return ctx.Status(fiber.StatusConflict).SendString("User is already registered")
		}

		log.Error().Err(err).Msg("Failed to create user (usecase)")
		// ? JSON RPC: TRANSPORT: 200, error: {msg, ..., dev_msg}
		// ? Возвращать JSON?
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create user")
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": userID})
}
