package pow

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/opravdin/word-of-wisdom/internal/configuration/env"
	"github.com/opravdin/word-of-wisdom/internal/logger"
	"golang.org/x/crypto/scrypt"
)

const (
	// randomSeedLength is the length of the random seed in bytes
	randomSeedLength = 16
)

// PoWUtils provides utility functions for Proof of Work operations
type PoWUtils struct {
	config *env.PowConfig
	logger logger.Logger
}

// NewPoWUtils creates a new PoWUtils instance with the provided configuration
func NewPoWUtils(config *env.PowConfig, log logger.Logger) *PoWUtils {
	return &PoWUtils{
		config: config,
		logger: log,
	}
}

// CalculateDifficultyLevel calculates a difficulty level based on request count
// The difficulty increases based on the configured requests per difficulty increase
func (u *PoWUtils) CalculateDifficultyLevel(requestCount int64) int {
	// Calculate the difficulty level
	level := int(requestCount / int64(u.config.RequestsPerDifficultyIncrease))

	// Ensure we don't exceed the maximum difficulty level
	if level > u.config.MaxDifficultyLevel {
		level = u.config.MaxDifficultyLevel
	}

	u.logger.Debug("Calculated difficulty level", "requestCount", requestCount, "level", level)
	return level
}

// GenerateRandomSeed generates a random seed for the PoW challenge
func (u *PoWUtils) GenerateRandomSeed() (string, error) {
	randomBytes := make([]byte, randomSeedLength)
	if _, err := rand.Read(randomBytes); err != nil {
		u.logger.Error("Failed to generate random seed", "error", err)
		return "", err
	}
	seed := hex.EncodeToString(randomBytes)
	u.logger.Debug("Generated random seed", "seed", seed)
	return seed, nil
}

// VerifySolution verifies that a solution meets the required difficulty level
// Returns true if the solution is valid, false otherwise
func (u *PoWUtils) VerifySolution(challengeID, seed, nonce string, difficultyLevel int) bool {
	u.logger.Debug("Verifying solution", "challengeID", challengeID, "nonce", nonce, "difficultyLevel", difficultyLevel)

	// Combine the challenge ID, seed, and nonce
	input := []byte(challengeID + seed + nonce)

	// Use configured scrypt parameters for verification
	hash, err := scrypt.Key(
		input,
		[]byte(challengeID),
		u.config.ScryptN,
		u.config.ScryptR,
		u.config.ScryptP,
		u.config.KeyLen,
	)
	if err != nil {
		u.logger.Error("Error computing scrypt hash", "error", err)
		return false
	}

	// Convert hash to hex string
	hashHex := hex.EncodeToString(hash)
	// Check if the hash has the required number of leading zeros
	prefix := strings.Repeat("0", difficultyLevel)
	isValid := strings.HasPrefix(hashHex, prefix)

	u.logger.Debug("Solution verification result",
		"challengeID", challengeID,
		"valid", isValid,
		"requiredPrefix", prefix,
		"hashPrefix", hashHex[:min(len(prefix)+2, len(hashHex))])

	return isValid
}

// min returns the smaller of x or y
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
