package config

type AppConfig struct {
	Port         string `env:"PORT"`
	JWTSecret    string `env:"JWT_SECRET"`
	IsProduction bool   `env:"PRODUCTION"`
	ClientURL    string `env:"CLIENT_URL"`

	ManualMapIdRequired bool `env:"MANUAL_MAP_ID_REQUIRED"`
}
