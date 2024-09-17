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
	// router := &Router{
	// 	App: fiber.New(),
	// 	// db:  db,
	// }
	app := fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")

	router := &Router{
		App: app,
		V1:  v1,
		// db:  db,
	}

	return router
}

func (f *Router) registerRoutes() {
	f.App.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowOrigins:     config.App.ClientURL,
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Example of setting up routes
	// f.App.Get("/users", handlers.GetUsers(f.db))
	// f.App.Post("/users", handlers.CreateUser(f.db))
}
