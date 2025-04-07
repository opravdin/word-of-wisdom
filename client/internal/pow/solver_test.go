package pow

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/opravdin/word-of-wisdom/client/internal/config"
	"github.com/opravdin/word-of-wisdom/client/internal/logger"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/scrypt"
)

// Test constants
const (
	testChallengeID    = "test-challenge-id"
	testSeed           = "test-seed"
	testDifficultyLow  = 1
	testDifficultyHigh = 10
	testTimeout        = 100 * time.Millisecond
)

func TestSolver_Solve(t *testing.T) {
	tests := []struct {
		name          string
		challengeID   string
		seed          string
		difficulty    int
		n             int
		r             int
		p             int
		keyLen        int
		shouldFail    bool
		errorContains string
	}{
		{
			name:        "should_solve_challenge_with_low_difficulty",
			challengeID: testChallengeID,
			seed:        testSeed,
			difficulty:  testDifficultyLow,
			n:           DefaultScryptN,
			r:           DefaultScryptR,
			p:           DefaultScryptP,
			keyLen:      DefaultScryptKeyLen,
			shouldFail:  false,
		},
		{
			name:        "should_use_default_parameters_when_zero",
			challengeID: testChallengeID,
			seed:        testSeed,
			difficulty:  testDifficultyLow,
			n:           0,
			r:           0,
			p:           0,
			keyLen:      0,
			shouldFail:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			log := logger.NewStdLogger()
			cfg := &config.ClientConfig{
				SolveTimeout: 5 * time.Second, // Longer timeout for tests
			}
			factory := NewSolverFactory(log, cfg)
			solver := factory.NewSolver()

			// Execute
			nonce, err := solver.Solve(tt.challengeID, tt.seed, tt.difficulty, tt.n, tt.r, tt.p, tt.keyLen)

			// Assert
			if tt.shouldFail {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, nonce)

				// Verify the solution meets the difficulty requirement
				n := DefaultScryptN
				r := DefaultScryptR
				p := DefaultScryptP
				keyLen := DefaultScryptKeyLen
				if tt.n > 0 {
					n = tt.n
				}
				if tt.r > 0 {
					r = tt.r
				}
				if tt.p > 0 {
					p = tt.p
				}
				if tt.keyLen > 0 {
					keyLen = tt.keyLen
				}

				hash, err := scryptKey([]byte(tt.challengeID+tt.seed+nonce), []byte(tt.challengeID), n, r, p, keyLen)
				assert.NoError(t, err)
				assert.True(t, strings.HasPrefix(hash, strings.Repeat("0", tt.difficulty)))
			}
		})
	}
}

func TestSolver_SolveWithTimeout(t *testing.T) {
	tests := []struct {
		name          string
		challengeID   string
		seed          string
		difficulty    int
		timeout       time.Duration
		shouldFail    bool
		errorContains string
	}{
		{
			name:          "should_timeout_with_high_difficulty",
			challengeID:   testChallengeID,
			seed:          testSeed,
			difficulty:    testDifficultyHigh,
			timeout:       testTimeout,
			shouldFail:    true,
			errorContains: "timed out",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			log := logger.NewStdLogger()
			cfg := &config.ClientConfig{
				SolveTimeout: tt.timeout,
			}
			factory := NewSolverFactory(log, cfg)
			solver := factory.NewSolver()

			// Create a context with the test timeout
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			// Execute
			_, err := solver.SolveWithContext(ctx, tt.challengeID, tt.seed, tt.difficulty, DefaultScryptN, DefaultScryptR, DefaultScryptP, DefaultScryptKeyLen)

			// Assert
			if tt.shouldFail {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSolver_SolveWithDefaults(t *testing.T) {
	tests := []struct {
		name        string
		challengeID string
		seed        string
		shouldFail  bool
	}{
		{
			name:        "should_solve_with_default_parameters",
			challengeID: testChallengeID,
			seed:        testSeed,
			shouldFail:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			log := logger.NewStdLogger()
			cfg := &config.ClientConfig{
				SolveTimeout: 5 * time.Second, // Longer timeout for tests
			}
			factory := NewSolverFactory(log, cfg)
			solver := factory.NewSolver()

			// Execute
			nonce, err := solver.SolveWithDefaults(tt.challengeID, tt.seed)

			// Assert
			if tt.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, nonce)

				// Verify the solution meets the default difficulty requirement (1)
				hash, err := scryptKey([]byte(tt.challengeID+tt.seed+nonce), []byte(tt.challengeID), DefaultScryptN, DefaultScryptR, DefaultScryptP, DefaultScryptKeyLen)
				assert.NoError(t, err)
				assert.True(t, strings.HasPrefix(hash, "0"))
			}
		})
	}
}

// Helper function to calculate scrypt hash and return as hex string
func scryptKey(password, salt []byte, N, r, p, keyLen int) (string, error) {
	hash, err := scrypt.Key(password, salt, N, r, p, keyLen)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash), nil
}
