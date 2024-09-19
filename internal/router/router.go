package router

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"

	"location-backend/internal/config"
)

type Router struct {
	App *fiber.App
	V1  fiber.Router
	// db  db.Service
}

func New() *Router {
	app := fiber.New()

	RegisterRoutes(app)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Use(jwtware.New(jwtware.Config{SigningKey: jwtware.SigningKey{Key: []byte(config.App.JWTSecret)}}))

	router := &Router{
		App: app,
		V1:  v1,
		// db:  db,
	}

	return router
}
