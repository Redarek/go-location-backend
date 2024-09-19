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
	// userURL  = "/users/:user_id"
	userGroup   = "/user"
	registerURL = "/register"
	loginURL    = "/login"
)

type userHandler struct {
	usecase usecase.UserUsecase
}

// Регистрирует новый handler
func NewUserHandler(userUsecase usecase.UserUsecase) *userHandler {
	return &userHandler{usecase: userUsecase}
}

// Регистрирует маршруты для user
func (h *userHandler) Register(router *router.Router) {
	userGruop := router.V1.Group(userGroup)
	userGruop.Post(registerURL, h.RegisterUser)
	userGruop.Post(loginURL, h.Login)
}

// func (h *bookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	// books := h.bookService.GetAll(context.Background(), 0, 0)
// 	w.Write([]byte("books"))
// 	w.WriteHeader(http.StatusOK)
// }

// Регистрирует нового пользователя, если его не существует.
//
// Возвращаемые статусы:
//
//	201 Created – пользователь успешно создан
//	409 Conflict – пользователь уже существует
//	500 InternalServerError – ошибка сервера
func (h *userHandler) RegisterUser(ctx *fiber.Ctx) error {
	// DTO from client (HTTP/JSON)
	var dto dto.RegisterUserDTO
	err := ctx.BodyParser(&dto)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse user request body")
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
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
	userID, err := h.usecase.Register(dto)
	if err != nil {
		if errors.Is(err, usecase.ErrAlreadyRegistered) {
			return ctx.Status(fiber.StatusConflict).SendString("User is already registered")
		}

		log.Error().Err(err).Msg("Failed to register new user")
		// ? JSON RPC: TRANSPORT: 200, error: {msg, ..., dev_msg}
		// ? Возвращать JSON?
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create user")
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": userID})
}

// Авторизует пользователя.
//
// Возвращаемые статусы:
//
//	200 OK – пользователь авторизован
//	400 BadRequest – переданы некорректные данные
//	401 Unauthorized – неверные логин/пароль или пользователя не существует
//	500 InternalServerError – ошибка сервера
func (h *userHandler) Login(ctx *fiber.Ctx) error {
	var dto dto.LoginUserDTO
	err := ctx.BodyParser(&dto)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse user request body")
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// TODO validate

	// TODO already login err

	token, err := h.usecase.Login(dto)
	if err != nil {
		if errors.Is(err, usecase.ErrBadLogin) {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Wrong login or password")
		}

		log.Error().Err(err).Msg("failed to login")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to login")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
