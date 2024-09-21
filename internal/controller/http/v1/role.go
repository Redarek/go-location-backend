package v1

// import (
// 	// "encoding/json"
// 	"errors"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/rs/zerolog/log"

// 	"location-backend/internal/controller/http/dto"
// 	"location-backend/internal/domain/usecase"
// 	"location-backend/internal/router"
// )

// const (
// 	roleGroup   = "/role"
// 	createURL = "/"
// )

// type roleHandler struct {
// 	usecase usecase.RoleUsecase
// }

// // Регистрирует новый handler
// func NewRoleHandler(usecase usecase.RoleUsecase) *roleHandler {
// 	return &roleHandler{usecase: usecase}
// }

// // Регистрирует маршруты для role
// func (h *roleHandler) Register(router *router.Router) {
// 	role := router.V1.Group(roleGroup)
// 	role.Post(createURL, h.CreateRole)
// 	// role.Post(loginURL, h.Login)
// }

// func (h *roleHandler) CreateRole(ctx *fiber.Ctx) error {
// 	var dto dto.CreateRoleDTO
// 	err := ctx.BodyParser(&dto)
// 	if err != nil {
// 		log.Error().Err(err).Msg("failed to parse user request body")
// 		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
// 	}

// 	// TODO validate

// 	// TODO already login err

// 	token, err := h.usecase.Login(dto)
// 	if err != nil {
// 		if errors.Is(err, usecase.ErrBadLogin) {
// 			return ctx.Status(fiber.StatusUnauthorized).SendString("Wrong login or password")
// 		}

// 		log.Error().Err(err).Msg("failed to login")
// 		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to login")
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
// }
