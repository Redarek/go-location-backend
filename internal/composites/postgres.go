package composites

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"location-backend/internal/config"
	"location-backend/pkg/client/postgres"
)

type PostgresComposite struct {
	pool *pgxpool.Pool
}

func NewPostgresComposite() (сomposite *PostgresComposite, err error) {
	// Connect to PostgreSQL
	pool, err := postgres.ConnectPostgres(&config.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	// defer pool.Close()

	return &PostgresComposite{pool: pool}, err
}
