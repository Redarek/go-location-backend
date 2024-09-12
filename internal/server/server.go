package server

import (
	// "location-backend/internal/db"

	"github.com/gofiber/fiber/v2"
	// "location-backend/internal/handlers" // Import your handlers package
)

type Fiber struct {
	App *fiber.App
	// db  db.Service
}

func New() *Fiber {
	server := &Fiber{
		App: fiber.New(),
		// db:  db,
	}

	// Initialize routes
	// server.registerRoutes()

	return server
}

// func (f *Fiber) registerRoutes() {
//     // Example of setting up routes
//     f.App.Get("/users", handlers.GetUsers(f.db))
//     f.App.Post("/users", handlers.CreateUser(f.db))
//     // Add more routes as needed
// }
