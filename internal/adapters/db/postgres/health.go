package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type HealthRepo interface {
	Health() (err error)
}

type healthRepo struct {
	pool *pgxpool.Pool
}

// TODO сделать единую точку создания репозиториев
func NewHealthRepo(pool *pgxpool.Pool) *healthRepo {
	return &healthRepo{pool: pool}
}

// Health pings database
func (r *healthRepo) Health() (err error) {
	// Creating a context with a timeout ensures that the health check does not hang indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the database to check connectivity.
	err = r.pool.Ping(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot ping postgres database")
	}

	log.Info().Msg("success ping postgres database")

	return
}
