package server

import (
	"github.com/gofiber/fiber/v2"
	"location-backend/internal/db"
)

type Fiber struct {
	App *fiber.App
	db  db.Service
}

func New(db db.Service) *Fiber {
	server := &Fiber{
		App: fiber.New(),
		db:  db,
	}

	return server
}
