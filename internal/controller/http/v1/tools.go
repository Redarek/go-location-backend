package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetUserIDFromJWT(ctx *fiber.Ctx) (userID uuid.UUID, err error) {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, err = uuid.Parse(claims["id"].(string))

	return
}
