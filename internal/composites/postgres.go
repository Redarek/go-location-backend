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

	return &PostgresComposite{pool: pool}, err
}
