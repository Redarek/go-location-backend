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

	// Глобальные маршруты
	api := app.Group("/api")
	v1 := api.Group("/v1")

	router := &Router{
		App: app,
		V1:  v1,
		// db:  db,
	}

	return router
}
