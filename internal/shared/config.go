package shared

type Config struct {
	Port        int    `env:"PORT" envDefault:"3000"`
	DatabaseURL string `env:"DATABASE_URL"`
	CorsOrigin  string `env:"CORS_ORIGIN" envDefault:"*"`
}
