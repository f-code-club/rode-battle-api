package shared

type Config struct {
	Port                int    `env:"PORT" envDefault:"3000"`
	DatabaseURL         string `env:"DATABASE_URL"`
	CorsOrigin          string `env:"CORS_ORIGIN" envDefault:"*"`
	JWTRefreshSecret    string `env:"JWT_REFRESH_SECRET" envDefault:"secret"`
	JWTRefreshExpiredIn int    `env:"JWT_REFRESH_EXPIRED_IN" envDefault:"86400"`
	JWTSecret           string `env:"JWT_SECRET" envDefault:"secret"`
	JWTExpiredIn        int    `env:"JWT_EXPIRED_IN" envDefault:"600"`
}
