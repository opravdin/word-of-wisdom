package env

import (
	"os"
	"strconv"
	"time"
)

// AppConfig holds all application configuration
type AppConfig struct {
	Server ServerConfig
	Redis  RedisConfig
	Pow    PowConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
}

// RedisConfig holds Redis-related configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// PowConfig holds all PoW-related configuration
type PowConfig struct {
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

	// Rate limiting
	MaxUnsolvedChallenges int
	BucketCapacity        int
}

// DefaultConfig returns the default application configuration
func DefaultConfig() *AppConfig {
	return &AppConfig{
		Server: ServerConfig{
			Port: "8080",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
		Pow: PowConfig{
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

			// Default rate limiting
			MaxUnsolvedChallenges: 10,
			BucketCapacity:        10,
		},
	}
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *AppConfig {
	config := DefaultConfig()

	// Load server configuration
	if port := getEnvWithDefault("PORT", config.Server.Port); port != "" {
		config.Server.Port = port
	}

	// Load Redis configuration
	if host := getEnvWithDefault("REDIS_HOST", config.Redis.Host); host != "" {
		config.Redis.Host = host
	}
	if port := getEnvWithDefault("REDIS_PORT", config.Redis.Port); port != "" {
		config.Redis.Port = port
	}
	if password := getEnvWithDefault("REDIS_PASSWORD", config.Redis.Password); password != "" {
		config.Redis.Password = password
	}
	if db := getEnvInt("REDIS_DB", config.Redis.DB); db >= 0 {
		config.Redis.DB = db
	}

	// Load PoW configuration
	if val := getEnvInt("POW_SCRYPT_N", config.Pow.ScryptN); val > 0 {
		config.Pow.ScryptN = val
	}
	if val := getEnvInt("POW_SCRYPT_R", config.Pow.ScryptR); val > 0 {
		config.Pow.ScryptR = val
	}
	if val := getEnvInt("POW_SCRYPT_P", config.Pow.ScryptP); val > 0 {
		config.Pow.ScryptP = val
	}
	if val := getEnvInt("POW_KEY_LEN", config.Pow.KeyLen); val > 0 {
		config.Pow.KeyLen = val
	}

	// Load challenge settings
	if val := getEnvDuration("POW_CHALLENGE_TTL", config.Pow.ChallengeTTL); val > 0 {
		config.Pow.ChallengeTTL = val
	}

	// Load difficulty settings
	if val := getEnvInt("POW_REQUESTS_PER_DIFFICULTY_INCREASE", config.Pow.RequestsPerDifficultyIncrease); val > 0 {
		config.Pow.RequestsPerDifficultyIncrease = val
	}
	if val := getEnvInt("POW_MAX_DIFFICULTY_LEVEL", config.Pow.MaxDifficultyLevel); val > 0 {
		config.Pow.MaxDifficultyLevel = val
	}

	// Load rate limiting settings
	if val := getEnvInt("POW_MAX_UNSOLVED_CHALLENGES", config.Pow.MaxUnsolvedChallenges); val > 0 {
		config.Pow.MaxUnsolvedChallenges = val
	}
	if val := getEnvInt("POW_BUCKET_CAPACITY", config.Pow.BucketCapacity); val > 0 {
		config.Pow.BucketCapacity = val
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
	if err != nil || val < 0 {
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
