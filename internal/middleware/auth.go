package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"location-backend/internal/config"
	"location-backend/pkg/httperrors"
)

// Auth returns the pre-initialized JWT middleware
var Auth fiber.Handler

func InitAuth() {
	log.Info().Msg("initializing JWT middleware")
	// Initialize the JWT middleware once
	Auth = jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.App.JWTSecret)},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			log.Info().Err(err).Msg("token validation failed")
			return ctx.Status(fiber.StatusUnauthorized).JSON(httperrors.NewErrorResponse(
				fiber.StatusUnauthorized,
				"Invalid or expired JWT",
				"",
				nil,
			))
		},
	})
}

// // Auth returns the pre-initialized JWT middleware
// func Auth() fiber.Handler {
// 	return jwtMiddleware
// }
