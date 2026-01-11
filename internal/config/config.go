package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	// Server settings
	ServerPort string
	ServerHost string

	// Default camera credentials (optional)
	DefaultUsername string
	DefaultPassword string

	// API settings
	APIPrefix string

	// Logging
	LogLevel string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "3000"),
		ServerHost:      getEnv("SERVER_HOST", "0.0.0.0"),
		DefaultUsername: getEnv("TAPO_DEFAULT_USERNAME", ""),
		DefaultPassword: getEnv("TAPO_DEFAULT_PASSWORD", ""),
		APIPrefix:       getEnv("API_PREFIX", "/api"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return c.ServerHost + ":" + c.ServerPort
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool gets a boolean environment variable with a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
