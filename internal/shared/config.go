package shared

type Config struct {
	Port                int    `env:"PORT" envDefault:"3000"`
	DatabaseURL         string `env:"DATABASE_URL"`
	AmqpURL             string `env:"AMQP_URL"`
	TaskQueue           string `env:"TASK_QUEUE" envDefault:"task"`
	CorsOrigin          string `env:"CORS_ORIGIN" envDefault:"*"`
	JWTRefreshSecret    string `env:"JWT_REFRESH_SECRET" envDefault:"secret"`
	JWTRefreshExpiredIn int    `env:"JWT_REFRESH_EXPIRED_IN" envDefault:"86400"`
	JWTAccessSecret     string `env:"JWT_ACCESS_SECRET" envDefault:"secret"`
	JWTAccessExpiredIn  int    `env:"JWT_ACCESS_EXPIRED_IN" envDefault:"600"`
	EmailUsername       string `env:"EMAIL_USERNAME"`
	EmailPassword       string `env:"EMAIL_PASSWORD"`
	EmailHost           string `env:"EMAIL_HOST"`
	EmailPort           string `env:"EMAIL_PORT"`
	AwsAccessKeyID      string `env:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey  string `env:"AWS_SECRET_ACCESS_KEY"`
	S3Bucket            string `env:"S3_BUCKET"`
	AwsRegion           string `env:"AWS_REGION"`
	AwsEndpointURL      string `env:"AWS_ENDPOINT_URL"`
	JudgeURL            string `env:"JUDGE_URL"`
}
