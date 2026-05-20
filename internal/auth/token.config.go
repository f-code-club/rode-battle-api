package auth

import "time"

type Config struct {
	Secret    string        `env:"JWT_SECRET" envRequired:"true"`
	ExpiredIn time.Duration `env:"JWT_EXPIRED_IN" envDefault:"5h"`
}
