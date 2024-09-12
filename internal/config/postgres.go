package config

type PostgresConfig struct {
	// URL      string `env:"DB_URL, required"` // docker run --name location-postgres -e DB_PASSWORD=postgres -e DB_USER=postgres -e DB_DATABASE=postgres -p 5432:5432 -d postgres
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	Database string `env:"DB_DATABASE"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
}
