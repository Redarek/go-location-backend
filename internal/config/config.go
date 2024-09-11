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

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}

	if err := env.Parse(&App); err != nil {
		log.Error().Err(err)
	}
	log.Debug().Msgf("%+v\n", App)

	if err := env.Parse(&Postgres); err != nil {
		log.Error().Err(err)
	}
	log.Debug().Msgf("%+v\n", Postgres)
}
