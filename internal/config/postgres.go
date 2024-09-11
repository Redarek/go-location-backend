package config

type PostgresConfig struct {
	// URL      string `env:"DB_URL, required"` // docker run --name location-postgres -e DB_PASSWORD=postgres -e DB_USER=postgres -e DB_DATABASE=postgres -p 5432:5432 -d postgres
	Host     string `env:"DB_HOST, required"`
	Port     int    `env:"DB_PORT, required"`
	Database string `env:"DB_DATABASE, required"`
	Username string `env:"DB_USERNAME, required"`
	Password string `env:"DB_PASSWORD, required"`
}
