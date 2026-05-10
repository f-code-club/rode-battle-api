package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Port        int    `env:"PORT" envDefault:"3000"`
	CorsOrigin  string `env:"CORS_ORIGIN" envDefault:"*"`
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://user:password@localhost/db"`
}

func Load() (Config, error) {
	return env.ParseAs[Config]()
}
