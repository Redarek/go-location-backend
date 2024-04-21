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

type PostgresConfig struct {
	URL string `env:"DB_URL,required"` // docker run --name location-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=postgres -p 5432:5432 -d postgres
}

type AppConfig struct {
	Port         string `env:"PORT,required"`
	JWTSecret    string `env:"JWT_SECRET,required"`
	IsProduction bool   `env:"PRODUCTION,required"`
	ClientURL    string `env:"CLIENT_URL,required"`
}

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
