package config

import (
	"os"
	"strconv"
	"time"
)

// ClientConfig holds the configuration for the client
type ClientConfig struct {
	// Server connection settings
	ServerAddress string
	HTTPAddress   string

	// Timeouts
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	// PoW settings
	DefaultScryptN      int
	DefaultScryptR      int
	DefaultScryptP      int
	DefaultScryptKeyLen int
	SolveTimeout        time.Duration
}

// LoadConfig loads the configuration from environment variables with defaults
func LoadConfig() *ClientConfig {
	return &ClientConfig{
		// Server connection settings
		ServerAddress: getEnv("SERVER_ADDRESS", "localhost:8080"),
		HTTPAddress:   getEnv("HTTP_ADDRESS", "localhost:3000"),

		// Timeouts
		ConnectTimeout: getDurationEnv("CONNECT_TIMEOUT", 10*time.Second),
		ReadTimeout:    getDurationEnv("READ_TIMEOUT", 30*time.Second),
		WriteTimeout:   getDurationEnv("WRITE_TIMEOUT", 30*time.Second),

		// PoW settings
		DefaultScryptN:      getIntEnv("DEFAULT_SCRYPT_N", 16384),
		DefaultScryptR:      getIntEnv("DEFAULT_SCRYPT_R", 8),
		DefaultScryptP:      getIntEnv("DEFAULT_SCRYPT_P", 1),
		DefaultScryptKeyLen: getIntEnv("DEFAULT_SCRYPT_KEY_LEN", 32),
		SolveTimeout:        getDurationEnv("SOLVE_TIMEOUT", 30*time.Second),
	}
}

// Validate validates the configuration
func (c *ClientConfig) Validate() error {
	// In a real implementation, we would validate the configuration here
	return nil
}

// Helper functions for environment variables

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getIntEnv gets an integer environment variable or returns a default value
func getIntEnv(key string, defaultValue int) int {
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

// getDurationEnv gets a duration environment variable or returns a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
