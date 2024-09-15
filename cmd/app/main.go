package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	// "location-backend/internal/adapters/db/repository"
	"location-backend/internal/config"
	// "location-backend/internal/controller/http/v1"
	// "location-backend/internal/domain/service"
	"location-backend/internal/server"
	"location-backend/pkg/client/postgres"
	"location-backend/pkg/logger"
)

func main() {
	logger.Setup()

	// Load the configuration
	config.LoadConfig()

	// Connect to PostgreSQL
	pool, err := postgres.ConnectPostgres(&config.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer pool.Close()

	// Sync tables
	err = postgres.SyncTables(pool)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to sync tables")
	}

	// log.Info().Msg("PostgreSQL connection and table sync completed successfully")

	// Initialize your services and handlers
	// userRepo := repository.NewUserRepo(pool)
	// userService := service.NewUserService(userRepo)
	// userHandler := v1.NewUserHandler(userUsecase)

	// Initialize and start the Fiber server
	router := server.New()
	go func() {
		if err := router.App.Listen(":" + config.App.Port); err != nil {
			log.Fatal().Err(err).Msg("Failed to start Fiber server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the application
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Info().Msg("Shutting down gracefully...")
	// TODO might need to add code here to shut down server properly
}
