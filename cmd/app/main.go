package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"location-backend/internal/config"
	"location-backend/pkg/client/postgres"
	"location-backend/pkg/logger"
)

func main() {
	logger.Setup()

	// Load the configuration
	config.LoadConfig()

	// Connect to PostgreSQL
	conn, err := postgres.ConnectPostgres(&config.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer conn.Close(context.Background())

	// Sync tables
	err = postgres.SyncTables(conn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to sync tables")
	}

	log.Info().Msg("PostgreSQL connection and table sync completed successfully")
}
