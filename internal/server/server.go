package server

import (
	// "location-backend/internal/db"

	"github.com/gofiber/fiber/v2"
	// "location-backend/internal/handlers" // Import your handlers package
	// "location-backend/internal/adapters/db/repository"
	// "location-backend/internal/controller/http/v1"
)

type Router struct {
	App *fiber.App
	// db  db.Service
}

// type Handlers struct {
//     UserHandler *v1.UserHandler
//     // Add other handlers here
// }

func New() *Router {
	router := &Router{
		App: fiber.New(),
		// db:  db,
	}

	// userRepo := repository.NewUserRepo()
	// server.registerRoutes()

	return router
}
