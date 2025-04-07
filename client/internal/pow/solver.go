package pow

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/opravdin/word-of-wisdom/client/internal/config"
	"github.com/opravdin/word-of-wisdom/client/internal/logger"
	"golang.org/x/crypto/scrypt"
)

// Default PoW parameters (used as fallback)
const (
	DefaultScryptN      = 16384
	DefaultScryptR      = 8
	DefaultScryptP      = 1
	DefaultScryptKeyLen = 32

	// Maximum number of attempts to find a solution
	MaxAttempts = 1000000
)

// solverImpl implements the Solver interface
type solverImpl struct {
	rand   *rand.Rand
	logger logger.Logger
	config *config.ClientConfig
}

// solverFactory implements the SolverFactory interface
type solverFactory struct {
	logger logger.Logger
	config *config.ClientConfig
}

// NewSolverFactory creates a new solver factory
func NewSolverFactory(logger logger.Logger, config *config.ClientConfig) SolverFactory {
	return &solverFactory{
		logger: logger,
		config: config,
	}
}

// NewSolver creates a new PoW solver
func (f *solverFactory) NewSolver() Solver {
	return &solverImpl{
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
		logger: f.logger.With("component", "pow_solver"),
		config: f.config,
	}
}

// Solve solves a PoW challenge by finding a nonce that produces a hash with the required number of leading zeros
func (s *solverImpl) Solve(challengeID, seed string, difficultyLevel, n, r, p, keyLen int) (string, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.config.SolveTimeout)
	defer cancel()

	return s.SolveWithContext(ctx, challengeID, seed, difficultyLevel, n, r, p, keyLen)
}

// SolveWithContext solves a PoW challenge with a context for cancellation/timeout
func (s *solverImpl) SolveWithContext(ctx context.Context, challengeID, seed string, difficultyLevel, n, r, p, keyLen int) (string, error) {
	s.logger.Debug("Solving PoW challenge",
		"challengeID", challengeID,
		"difficultyLevel", difficultyLevel,
		"n", n,
		"r", r,
		"p", p,
		"keyLen", keyLen)

	// Use default values if parameters are invalid
	if n <= 0 {
		n = DefaultScryptN
	}
	if r <= 0 {
		r = DefaultScryptR
	}
	if p <= 0 {
		p = DefaultScryptP
	}
	if keyLen <= 0 {
		keyLen = DefaultScryptKeyLen
	}

	// Generate random nonces and check if they produce a valid hash
	for attempt := 0; attempt < MaxAttempts; attempt++ {
		// Check if context is done (timeout or cancellation)
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("solving timed out after %v", s.config.SolveTimeout)
		default:
			// Continue with the next attempt
		}

		// Generate a random nonce
		nonce := generateNonce(s.rand)

		// Calculate hash with the nonce
		input := []byte(challengeID + seed + nonce)
		hash, err := scrypt.Key(input, []byte(challengeID), n, r, p, keyLen)
		if err != nil {
			return "", fmt.Errorf("error calculating scrypt: %w", err)
		}

		// Check if the hash meets the difficulty requirement
		hashHex := hex.EncodeToString(hash)
		if strings.HasPrefix(hashHex, strings.Repeat("0", difficultyLevel)) {
			s.logger.Debug("Found solution",
				"attempt", attempt,
				"nonce", nonce,
				"hash", hashHex[:16]+"...")
			return nonce, nil
		}

		// Log progress periodically
		if attempt > 0 && attempt%10000 == 0 {
			s.logger.Debug("Still searching for solution", "attempts", attempt)
		}
	}

	return "", fmt.Errorf("failed to find solution after %d attempts", MaxAttempts)
}

// SolveWithDefaults solves a PoW challenge using default parameters
// This is for backward compatibility with older servers
func (s *solverImpl) SolveWithDefaults(challengeID, seed string) (string, error) {
	s.logger.Debug("Solving PoW challenge with default parameters",
		"challengeID", challengeID)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.config.SolveTimeout)
	defer cancel()

	// For backward compatibility, assume difficulty level 1
	return s.SolveWithContext(ctx, challengeID, seed, 1, DefaultScryptN, DefaultScryptR, DefaultScryptP, DefaultScryptKeyLen)
}

// generateNonce generates a random nonce for PoW
func generateNonce(r *rand.Rand) string {
	// Define nonce length as a constant to avoid magic number
	const nonceLength = 8

	nonceBytes := make([]byte, nonceLength)
	for i := range nonceBytes {
		nonceBytes[i] = byte(r.Intn(256))
	}
	return hex.EncodeToString(nonceBytes)
}
