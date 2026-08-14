package shared

type Config struct {
	Port                int    `env:"PORT" envDefault:"3000"`
	DatabaseURL         string `env:"DATABASE_URL"`
	CorsOrigin          string `env:"CORS_ORIGIN" envDefault:"*"`
	JWTRefreshSecret    string `env:"JWT_REFRESH_SECRET" envDefault:"secret"`
	JWTRefreshExpiredIn int    `env:"JWT_REFRESH_EXPIRED_IN" envDefault:"86400"`
	JWTAccessSecret     string `env:"JWT_ACCESS_SECRET" envDefault:"secret"`
	JWTAccessExpiredIn  int    `env:"JWT_ACCESS_EXPIRED_IN" envDefault:"600"`
	MailerUsername      string `env:"MAIL_USERNAME" envDefault:"admin@example.com"`
	MailerPassword      string `env:"MAIL_PASSWORD" envDefault:"password"`
	MailerHostName      string `env:"MAIL_HOST_NAME" envDefault:"smtp.gmail.com"`
	MailerPortNumber    string `env:"MAIL_PORT_NUMBER" envDefault:"587"`
}
