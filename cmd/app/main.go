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
	log.Info().Msg("connecting to PostgreSQL...")
	postgresComposite, err := composites.NewPostgresComposite()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database composite")
	}

	// TODO в композит
	// Init middleware
	log.Info().Msg("initializing middleware...")
	middleware.InitAuth()

	// Create router
	log.Info().Msg("initializing router...")
	r := router.New()

	// Create composites
	log.Info().Msg("initializing composites...")
	repositoryComposit := composites.NewRepositoryComposite(postgresComposite)
	serviceComposit := composites.NewServiceComposite(repositoryComposit)
	usecaseComposite := composites.NewUsecaseComposite(serviceComposit)
	handlerComposite := composites.NewHandlerComposite(usecaseComposite)

	// Register routes
	log.Info().Msg("registering routes...")
	router.RegisterRoutes(r, handlerComposite)

	// Initialize and start the Fiber server
	log.Info().Msg("initializing and starting the Fiber server...")
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
