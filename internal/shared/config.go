package shared

type Config struct {
	Port                int    `env:"PORT" envDefault:"3000"`
	DatabaseURL         string `env:"DATABASE_URL"`
	CorsOrigin          string `env:"CORS_ORIGIN" envDefault:"*"`
	JWTRefreshSecret    string `env:"JWT_REFRESH_SECRET" envDefault:"secret"`
	JWTRefreshExpiredIn int    `env:"JWT_REFRESH_EXPIRED_IN" envDefault:"86400"`
	JWTAccessSecret     string `env:"JWT_ACCESS_SECRET" envDefault:"secret"`
	JWTAccessExpiredIn  int    `env:"JWT_ACCESS_EXPIRED_IN" envDefault:"600"`
	EmailUsername       string `env:"EMAIL_USERNAME"`
	EmailPassword       string `env:"EMAIL_PASSWORD"`
	EmailHost           string `env:"EMAIL_HOST"`
	EmailPort           string `env:"EMAIL_PORT"`
}
