package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"location-backend/pkg/client/postgres"
)

var (
	Postgres postgres.PostgresConfig
	App      AppConfig
)

// TODO
type Config struct {
	IsDebug bool `env:"IS_DEBUG" env-default:"false"`
}

func LoadConfig() {
	// Load from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("error loading .env file")
	}

	// Parse environment variables
	if err := env.Parse(&Postgres); err != nil {
		log.Fatal().Err(err).Msg("failed to parse PostgreSQL config")
	}

	if err := env.Parse(&App); err != nil {
		log.Fatal().Err(err).Msg("failed to parse App config")
	}

	log.Info().Msg("configuration loaded successfully")
}
