package pow

import (
	"os"
	"strconv"
	"time"
)

// Config holds all PoW-related configuration
type Config struct {
	// Scrypt parameters
	ScryptN int
	ScryptR int
	ScryptP int
	KeyLen  int

	// Challenge settings
	ChallengeTTL time.Duration

	// Difficulty settings
	RequestsPerDifficultyIncrease int
	MaxDifficultyLevel            int
}

// DefaultConfig returns the default PoW configuration
func DefaultConfig() *Config {
	return &Config{
		// Default Scrypt parameters
		ScryptN: 16384,
		ScryptR: 8,
		ScryptP: 1,
		KeyLen:  32,

		// Default challenge settings
		ChallengeTTL: 5 * time.Minute,

		// Default difficulty settings
		RequestsPerDifficultyIncrease: 10,
		MaxDifficultyLevel:            8,
	}
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	config := DefaultConfig()

	// Load Scrypt parameters
	if val := getEnvInt("POW_SCRYPT_N", config.ScryptN); val > 0 {
		config.ScryptN = val
	}
	if val := getEnvInt("POW_SCRYPT_R", config.ScryptR); val > 0 {
		config.ScryptR = val
	}
	if val := getEnvInt("POW_SCRYPT_P", config.ScryptP); val > 0 {
		config.ScryptP = val
	}
	if val := getEnvInt("POW_KEY_LEN", config.KeyLen); val > 0 {
		config.KeyLen = val
	}

	// Load challenge settings
	if val := getEnvDuration("POW_CHALLENGE_TTL", config.ChallengeTTL); val > 0 {
		config.ChallengeTTL = val
	}

	// Load difficulty settings
	if val := getEnvInt("POW_REQUESTS_PER_DIFFICULTY_INCREASE", config.RequestsPerDifficultyIncrease); val > 0 {
		config.RequestsPerDifficultyIncrease = val
	}
	if val := getEnvInt("POW_MAX_DIFFICULTY_LEVEL", config.MaxDifficultyLevel); val > 0 {
		config.MaxDifficultyLevel = val
	}

	return config
}

// Helper functions for environment variable parsing

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	valStr := getEnvWithDefault(key, strconv.Itoa(defaultValue))
	val, err := strconv.Atoi(valStr)
	if err != nil || val <= 0 {
		return defaultValue
	}
	return val
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	valStr := getEnvWithDefault(key, "")
	if valStr == "" {
		return defaultValue
	}

	// Try to parse as seconds first
	seconds, err := strconv.Atoi(valStr)
	if err == nil {
		return time.Duration(seconds) * time.Second
	}

	// Try to parse as duration string
	duration, err := time.ParseDuration(valStr)
	if err != nil {
		return defaultValue
	}
	return duration
}
