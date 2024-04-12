package config

import (
	"github.com/caarlos0/env/v10"
	"log"
)

var (
	Postgres PostgresConfig
	App      AppConfig
)

type PostgresConfig struct {
	URL string `env:"DB_URL" envDefault:"postgresql://postgres:postgres@localhost:5432/postgres"` // docker run --name location-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=postgres -p 5432:5432 -d postgres
}

type AppConfig struct {
	Port         string `env:"PORT" envDefault:"3000"`
	JWTSecret    string `env:"JWT_SECRET""`
	IsProduction bool   `env:"PRODUCTION"`
}

func Init() {
	if err := env.Parse(&App); err != nil {
		log.Printf("%+v\n", err)
	}
	log.Printf("%+v\n", App)

	if err := env.Parse(&Postgres); err != nil {
		log.Printf("%+v\n", err)
	}

}
