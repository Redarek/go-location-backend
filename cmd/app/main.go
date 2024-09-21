package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"location-backend/internal/composites"
	"location-backend/internal/config"
	"location-backend/internal/middleware"
	"location-backend/internal/router"
	"location-backend/pkg/logger"
)

func Start() {
	logger.Setup()

	// Load the configuration
	config.LoadConfig()

	// Connect to PostgreSQL
	postgresComposite, err := composites.NewPostgresComposite()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database composite")
	}

	// TODO в композит
	// Init middleware
	middleware.InitAuth()

	// Create router
	r := router.New()

	repositoryComposit := composites.NewRepositoryComposite(postgresComposite)
	serviceComposit := composites.NewServiceComposite(repositoryComposit)
	usecaseComposite := composites.NewUsecaseComposite(serviceComposit)
	handlerComposite := composites.NewHandlerComposite(usecaseComposite)

	// Register common routes
	router.RegisterRoutes(r, handlerComposite)

	// Register routes
	// TODO err
	// TODO вынести в отдельный файл

	// // Глобальные маршруты
	// api := r.App.Group("/api")
	// v1 := api.Group("/v1")

	// healthComposite := composites.NewHealthComposite(postgresComposite)
	// healthComposite.Handler.Register(&v1)

	// userComposite := composites.NewUserComposite(postgresComposite)
	// userComposite.Handler.Register(&v1)

	// Initialize and start the Fiber server
	go func() {
		if err := r.App.Listen(":" + config.App.Port); err != nil {
			log.Fatal().Err(err).Msg("failed to start Fiber server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the application
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Info().Msg("shutting down gracefully...")
	// TODO might need to add code here to shut down server properly
}
