package auth

type Config struct {
	Secret    string `env:"JWT_SECRET" envRequired:"true"`
	ExpiredIn int    `env:"JWT_EXPIRED_IN" envDefault:"5"` // hours
}
