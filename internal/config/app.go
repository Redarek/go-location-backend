package config

type AppConfig struct {
	Port         string `env:"PORT, required"`
	JWTSecret    string `env:"JWT_SECRET, required"`
	IsProduction bool   `env:"PRODUCTION, required"`
	ClientURL    string `env:"CLIENT_URL, required"`
}
