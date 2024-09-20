package router

import (
	"github.com/gofiber/fiber/v2"
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

	// TODO решить что с этим делать
	// v1.Use(jwtware.New(jwtware.Config{SigningKey: jwtware.SigningKey{Key: []byte(config.App.JWTSecret)}}))

	router := &Router{
		App: app,
		V1:  v1,
		// db:  db,
	}

	return router
}
