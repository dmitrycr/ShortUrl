package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config содержит конфигурацию приложения
type Config struct {
	// HTTP Server
	ServerPort string
	BaseURL    string

	// Database
	DatabaseURL string

	// URL Shortener
	CodeLength int

	// Environment
	Environment string // dev, staging, production
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		CodeLength:  getEnvAsInt("CODE_LENGTH", 6),
		Environment: getEnv("ENVIRONMENT", "dev"),
	}

	// Валидация обязательных параметров
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает переменную окружения как int
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// IsDevelopment проверяет, запущено ли приложение в dev режиме
func (c *Config) IsDevelopment() bool {
	return c.Environment == "dev"
}

// IsProduction проверяет, запущено ли приложение в production режиме
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
