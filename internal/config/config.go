package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var (
	Postgres PostgresConfig
	App      AppConfig
)

// TODO
type Config struct {
	IsDebug bool `env:"IS_DEBUG" env-default:"false"`
}

func LoadConfig() {
	// Load from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	// Parse environment variables
	if err := env.Parse(&Postgres); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse PostgreSQL config")
	}

	if err := env.Parse(&App); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse App config")
	}

	log.Info().Msg("Configuration loaded successfully")
}
