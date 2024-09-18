package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"location-backend/internal/config"
)

type Router struct {
	App *fiber.App
	V1  fiber.Router
	// db  db.Service
}

func New() *Router {
	app := fiber.New()

	// Cors
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowOrigins:     config.App.ClientURL,
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Static("/public", "../../public")

	api := app.Group("/api")
	v1 := api.Group("/v1")

	router := &Router{
		App: app,
		V1:  v1,
		// db:  db,
	}

	return router
}
