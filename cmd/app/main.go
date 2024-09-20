package main

import (
	"os"
	"os/signal"
	"syscall"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/rs/zerolog/log"

	"location-backend/internal/composites"
	"location-backend/internal/config"
	"location-backend/internal/router"
	"location-backend/pkg/logger"
)

func main() {
	logger.Setup()

	// Load the configuration
	config.LoadConfig()

	// Connect to PostgreSQL
	postgresComposite, err := composites.NewPostgresComposite()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database composite")
	}

	// Initialize and start the Fiber server
	router := router.New()
	go func() {
		if err := router.App.Listen(":" + config.App.Port); err != nil {
			log.Fatal().Err(err).Msg("failed to start Fiber server")
		}
	}()

	// TODO err
	// TODO вынести в отдельный файл
	healthComposite := composites.NewHealthComposite(postgresComposite)
	healthComposite.Handler.Register(router)

	userComposite := composites.NewUserComposite(postgresComposite)
	userComposite.Handler.Register(router)

	// TODO структурировать!
	router.V1.Use(jwtware.New(jwtware.Config{SigningKey: jwtware.SigningKey{Key: []byte(config.App.JWTSecret)}}))

	// Wait for interrupt signal to gracefully shutdown the application
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Info().Msg("Shutting down gracefully...")
	// TODO might need to add code here to shut down server properly
}
