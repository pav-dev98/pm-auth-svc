package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	GRPCPort  string
	DSN       string
	JWTSecret string
	JWTExpiration string
}

func Load() *Config {
	godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "pm_auth"),
	)

	return &Config{
		GRPCPort:      getEnv("GRPC_PORT", "50051"),
		DSN:           dsn,
		JWTSecret:     getEnv("JWT_SECRET", "secret"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}