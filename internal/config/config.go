package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	GRPC     GRPCConfig
}

type AppConfig struct {
	Name    string
	Env     string
	Port    string
	URL     string
	Version string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Name     string
	UseSSL   bool
	Timezone string
}

type GRPCConfig struct {
	Port               string
	ExternalServiceURL string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "go-gin-boilerplate"),
			Env:     getEnv("APP_ENV", "development"),
			Port:    getEnv("APP_PORT", "8000"),
			URL:     getEnv("APP_URL", "http://localhost:8000"),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Pass:     getEnv("DB_PASS", "postgres"),
			Name:     getEnv("DB_NAME", "go_gin_boilerplate"),
			UseSSL:   getEnvBool("DB_USE_SSL", false),
			Timezone: getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		},
		GRPC: GRPCConfig{
			Port:               getEnv("GRPC_PORT", "50051"),
			ExternalServiceURL: getEnv("GRPC_EXTERNAL_SERVICE_URL", "localhost:50052"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		parsed, err := strconv.ParseBool(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
