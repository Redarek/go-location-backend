package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"location-backend/internal/config"
	"time"
)

type postgres struct {
	*pgxpool.Pool
}

// New initializes a new postgres connection.
func New() Service {

	pool, err := pgxpool.New(context.Background(), config.Postgres.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to postgres")
	}
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Unable to take conn from pool")
	}
	defer conn.Release()

	// It's not recommended to defer conn.Close() here because this will close the connection immediately after New() finishes
	// Instead, ensure that the connection is closed outside of this function when it's no longer needed

	query := `
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    UNIQUE (user_id, role_id)
);


CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token VARCHAR(1500) NOT NULL,
    expiry TIMESTAMPTZ NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL
);


CREATE TABLE IF NOT EXISTS sites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS buildings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    country VARCHAR NOT NULL,
    city VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    site_id UUID NOT NULL REFERENCES sites(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS floors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    number INTEGER NOT NULL,
    image VARCHAR,
    heatmap VARCHAR,
    width_in_pixels INTEGER NOT NULL DEFAULT 0,
    height_in_pixels INTEGER NOT NULL DEFAULT 0,
    scale FLOAT NOT NULL CHECK (scale > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    building_id UUID NOT NULL REFERENCES buildings(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS access_point_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    color VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    site_id UUID NOT NULL REFERENCES sites(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS radio_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number INTEGER NOT NULL CHECK (number > 0),
    channel INTEGER NOT NULL CHECK (channel > 0),
    wifi VARCHAR NOT NULL,
    power INTEGER NOT NULL,
    bandwidth VARCHAR NOT NULL,
    guard_interval INTEGER NOT NULL CHECK (guard_interval > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    access_point_type_id UUID NOT NULL REFERENCES access_point_types(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS access_points (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    x INTEGER NOT NULL CHECK (x > 0),
    y INTEGER NOT NULL CHECK (y > 0),
    z FLOAT NOT NULL CHECK (z > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    floor_id UUID NOT NULL REFERENCES floors(id) ON DELETE SET NULL,
    access_point_type_id UUID NOT NULL REFERENCES access_point_types(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS radios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number INTEGER NOT NULL CHECK (number > 0),
    channel INTEGER NOT NULL CHECK (channel > 0),
    wifi VARCHAR NOT NULL,
    power INTEGER NOT NULL,
    bandwidth VARCHAR NOT NULL,
    guard_interval INTEGER NOT NULL CHECK (guard_interval > 0),
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    access_point_id UUID NOT NULL REFERENCES access_points(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS sensor_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    color VARCHAR NOT NULL,
    alias VARCHAR NOT NULL,
    interface_0 VARCHAR NOT NULL,
    interface_1 VARCHAR NOT NULL,
    interface_2 VARCHAR NOT NULL,
    rx_ant_gain FLOAT NOT NULL DEFAULT 0, -- TODO: add check
    hor_rotation_offset INTEGER NOT NULL DEFAULT 0, -- TODO: add check
    vert_rotation_offset INTEGER NOT NULL DEFAULT 0, -- TODO: add check
    correction_factor_24 FLOAT NOT NULL DEFAULT 0, -- TODO: add check
    correction_factor_5 FLOAT NOT NULL DEFAULT 0, -- TODO: add check
    correction_factor_6 FLOAT NOT NULL DEFAULT 0, -- TODO: add check
    diagram JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    site_id UUID NOT NULL REFERENCES sites(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS sensors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    x INTEGER NOT NULL CHECK (x > 0),
    y INTEGER NOT NULL CHECK (y > 0),
    z FLOAT NOT NULL CHECK (z > 0),
    mac VARCHAR UNIQUE NOT NULL,
    ip VARCHAR NOT NULL,
    alias VARCHAR NOT NULL,
    interface_0 VARCHAR NOT NULL,
    interface_1 VARCHAR NOT NULL,
    interface_2 VARCHAR NOT NULL,
    rx_ant_gain FLOAT NOT NULL DEFAULT 0, -- TODO: add check
    hor_rotation_offset INTEGER NOT NULL DEFAULT 0, -- TODO: add check
    vert_rotation_offset INTEGER NOT NULL DEFAULT 0, -- TODO: add check
    correction_factor_24 FLOAT NOT NULL DEFAULT 0, -- TODO: add check
    correction_factor_5 FLOAT NOT NULL DEFAULT 0, -- TODO: add check
    correction_factor_6 FLOAT NOT NULL DEFAULT 0, -- TODO: add check
    diagram JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    floor_id UUID NOT NULL REFERENCES floors(id) ON DELETE SET NULL,
	sensor_type_id UUID NOT NULL REFERENCES sensor_types(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS wall_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    color VARCHAR NOT NULL,
    attenuation_24 FLOAT NOT NULL CHECK (attenuation_24 > 0),
    attenuation_5 FLOAT NOT NULL CHECK (attenuation_5 > 0),
    attenuation_6 FLOAT NOT NULL CHECK (attenuation_6 > 0),
    thickness FLOAT NOT NULL CHECK (thickness > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    site_id UUID NOT NULL REFERENCES sites(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS walls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    x1 INTEGER NOT NULL CHECK (x1 > 0),
    y1 INTEGER NOT NULL CHECK (y1 > 0),
    x2 INTEGER NOT NULL CHECK (x2 > 0),
    y2 INTEGER NOT NULL CHECK (y2 > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    floor_id UUID NOT NULL REFERENCES floors(id) ON DELETE SET NULL,
    wall_type_id UUID NOT NULL REFERENCES wall_types(id) ON DELETE SET NULL
);

-- Relation between walls and wall types
--ALTER TABLE walls ADD COLUMN wall_type_id UUID REFERENCES wall_types(id) ON DELETE SET NULL;

-- Активация расширения для генерации UUID
CREATE EXTENSION IF NOT EXISTS pgcrypto;`

	_, err = conn.Exec(context.Background(), query)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create tables")
		return nil
	}
	db := &postgres{pool}
	return db
}

// Health pings database
func (p *postgres) Health() map[string]string {
	// Creating a context with a timeout ensures that the health check does not hang indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the database to check connectivity.
	err := p.Ping(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("db down!")
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
