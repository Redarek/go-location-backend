package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"location-backend/internal/db"
)

type Fiber struct {
	App *fiber.App
	db  db.Service
}

func New(db db.Service) *Fiber {
	app := fiber.New()
	app.Use(cors.New())
	server := &Fiber{
		App: app,
		db:  db,
	}

	return server
}
