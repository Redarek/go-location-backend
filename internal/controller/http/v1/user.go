package v1

import (
	// "encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/usecase"
	"location-backend/internal/middleware"
	"location-backend/pkg/httperrors"
)

const (
	// TODO user -> users
	// TODO logout
	// userURL  = "/user/:user_id"
	getUserByNameURL = "/"
	registerURL      = "/register"
	loginURL         = "/login"
)

type userHandler struct {
	usecase usecase.UserUsecase
}

// Регистрирует новый handler
func NewUserHandler(usecase usecase.UserUsecase) *userHandler {
	return &userHandler{usecase: usecase}
}

// Регистрирует маршруты для user
func (h *userHandler) Register(r *fiber.Router) fiber.Router {
	router := *r

	router.Post(registerURL, h.RegisterUser)
	router.Post(loginURL, h.Login)
	router.Get(getRoleByNameURL, middleware.Auth, h.GetUserByName)

	return router
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
	var dto http_dto.RegisterUserDTO
	err := ctx.BodyParser(&dto)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse user request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse user request body",
			nil,
		))
	}

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.RegisterUserDTO{
		Username: dto.Username,
		Password: dto.Password,
	}

	// ? Нужно ли передавать ctx внутрь?
	// Call the use case to create the user
	userID, err := h.usecase.Register(domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrAlreadyExists) {
			return ctx.Status(fiber.StatusConflict).JSON(httperrors.NewErrorResponse(
				fiber.StatusConflict,
				"User is already registered",
				"",
				nil,
			))
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to register a new user")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to register a new user",
			"",
			nil,
		))
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
	var dto http_dto.LoginUserDTO
	err := ctx.BodyParser(&dto)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse user request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperrors.NewErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			"Failed to parse user request body",
			nil,
		))
	}

	// TODO validate

	// TODO already login err

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.LoginUserDTO{
		Username: dto.Username,
		Password: dto.Password,
	}

	token, err := h.usecase.Login(domain_dto.LoginUserDTO(domainDTO))
	if err != nil {
		if errors.Is(err, usecase.ErrBadLogin) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(httperrors.NewErrorResponse(
				fiber.StatusUnauthorized,
				"Wrong login or password",
				"",
				nil,
			))
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to log in")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to log in",
			"",
			nil,
		))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *userHandler) GetUserByName(ctx *fiber.Ctx) error {
	var dto http_dto.GetUserByNameDTO = http_dto.GetUserByNameDTO{
		Username: ctx.Query("username"),
	}
	// accessPointID, err := uuid.Parse(c.Query("id"))

	// err := ctx.BodyParser(&dto)
	// if err != nil {
	// 	log.Error().Err(err).Msg("failed to parse user request body")
	// 	return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	// }

	// TODO validate

	// Mapping http DTO -> domain DTO
	domainDTO := domain_dto.GetUserByNameDTO{
		Username: dto.Username,
	}

	user, err := h.usecase.GetUserByName(domainDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			ctx.Status(fiber.StatusNoContent)
			return nil
		}

		log.Error().Err(err).Msg("an unexpected error has occurred while trying to get user")
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperrors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"An unexpected error has occurred while trying to get user",
			"",
			nil,
		))
	}

	// Mapping domain DTO -> http DTO
	userDTO := http_dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": userDTO})
}
