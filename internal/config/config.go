package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	RabbitMQ   RabbitMQConfig
	JWT        JWTConfig
	CORS       CORSConfig
	RateLimit  RateLimitConfig
	Log        LogConfig
	FileUpload FileUploadConfig
	Worker     WorkerConfig
}

type AppConfig struct {
	Env  string
	Port int
	Name string
}

type DatabaseConfig struct {
	Host               string
	Port               int
	User               string
	Password           string
	DBName             string
	SSLMode            string
	MaxConnections     int
	MaxIdleConnections int
	ConnMaxLifetime    time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

type RabbitMQConfig struct {
	URL         string
	Exchange    string
	QueuePrefix string
}

type JWTConfig struct {
	Secret              string
	ExpiryHours         int
	RefreshExpiryDays   int
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type RateLimitConfig struct {
	RequestsPerMinute int
	Burst             int
}

type LogConfig struct {
	Level  string
	Format string
}

type FileUploadConfig struct {
	MaxSize int64
}

type WorkerConfig struct {
	PoolSize  int
	QueueSize int
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	viper.AutomaticEnv()

	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("APP_NAME", "POS-System")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_MAX_CONNECTIONS", 100)
	viper.SetDefault("DB_MAX_IDLE_CONNECTIONS", 10)
	viper.SetDefault("DB_CONNECTION_MAX_LIFETIME", "1h")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_POOL_SIZE", 10)

	viper.SetDefault("RABBITMQ_EXCHANGE", "pos_events")
	viper.SetDefault("RABBITMQ_QUEUE_PREFIX", "pos_")

	viper.SetDefault("JWT_EXPIRY_HOURS", 24)
	viper.SetDefault("JWT_REFRESH_EXPIRY_DAYS", 7)

	viper.SetDefault("RATE_LIMIT_REQUESTS_PER_MINUTE", 60)
	viper.SetDefault("RATE_LIMIT_BURST", 10)

	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "json")

	viper.SetDefault("MAX_UPLOAD_SIZE", 10485760)

	viper.SetDefault("WORKER_POOL_SIZE", 10)
	viper.SetDefault("WORKER_QUEUE_SIZE", 100)

	connMaxLifetime, _ := time.ParseDuration(viper.GetString("DB_CONNECTION_MAX_LIFETIME"))

	config := &Config{
		App: AppConfig{
			Env:  viper.GetString("APP_ENV"),
			Port: viper.GetInt("APP_PORT"),
			Name: viper.GetString("APP_NAME"),
		},
		Database: DatabaseConfig{
			Host:               viper.GetString("DB_HOST"),
			Port:               viper.GetInt("DB_PORT"),
			User:               viper.GetString("DB_USER"),
			Password:           viper.GetString("DB_PASSWORD"),
			DBName:             viper.GetString("DB_NAME"),
			SSLMode:            viper.GetString("DB_SSLMODE"),
			MaxConnections:     viper.GetInt("DB_MAX_CONNECTIONS"),
			MaxIdleConnections: viper.GetInt("DB_MAX_IDLE_CONNECTIONS"),
			ConnMaxLifetime:    connMaxLifetime,
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetInt("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
			PoolSize: viper.GetInt("REDIS_POOL_SIZE"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:         viper.GetString("RABBITMQ_URL"),
			Exchange:    viper.GetString("RABBITMQ_EXCHANGE"),
			QueuePrefix: viper.GetString("RABBITMQ_QUEUE_PREFIX"),
		},
		JWT: JWTConfig{
			Secret:            viper.GetString("JWT_SECRET"),
			ExpiryHours:       viper.GetInt("JWT_EXPIRY_HOURS"),
			RefreshExpiryDays: viper.GetInt("JWT_REFRESH_EXPIRY_DAYS"),
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
			AllowedMethods: viper.GetStringSlice("CORS_ALLOWED_METHODS"),
			AllowedHeaders: viper.GetStringSlice("CORS_ALLOWED_HEADERS"),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: viper.GetInt("RATE_LIMIT_REQUESTS_PER_MINUTE"),
			Burst:             viper.GetInt("RATE_LIMIT_BURST"),
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
		FileUpload: FileUploadConfig{
			MaxSize: viper.GetInt64("MAX_UPLOAD_SIZE"),
		},
		Worker: WorkerConfig{
			PoolSize:  viper.GetInt("WORKER_POOL_SIZE"),
			QueueSize: viper.GetInt("WORKER_QUEUE_SIZE"),
		},
	}

	return config, nil
}