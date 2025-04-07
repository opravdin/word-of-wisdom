package pow

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/opravdin/word-of-wisdom/internal/configuration/env"
	"github.com/opravdin/word-of-wisdom/internal/logger"
	"github.com/opravdin/word-of-wisdom/internal/protocol"
	"github.com/opravdin/word-of-wisdom/internal/random"
	powrepo "github.com/opravdin/word-of-wisdom/internal/repository/pow"
)

// Service implements the Proof of Work service
type Service struct {
	repo                  Repository
	config                *env.PowConfig
	utils                 PoWUtilsInterface
	maxUnsolvedChallenges int
	logger                logger.Logger
	random                RandomProvider
}

// NewService creates a new Proof of Work service
func NewService(repo Repository, config *env.PowConfig, log logger.Logger) *Service {
	return &Service{
		repo:                  repo,
		config:                config,
		utils:                 NewPoWUtils(config, log),
		maxUnsolvedChallenges: config.MaxUnsolvedChallenges,
		logger:                log,
		random:                random.NewProvider(),
	}
}

// CreateChallenge creates a new PoW challenge for a client
func (s *Service) CreateChallenge(ctx context.Context, clientIP string) (*protocol.Challenge, error) {
	s.logger.Debug("Creating challenge for client", "ip", clientIP)

	// Check and increment unsolved count
	unsolvedCount, err := s.repo.IncrementUnsolvedCount(ctx, clientIP)
	if err != nil {
		s.logger.Error("Failed to increment unsolved count", "ip", clientIP, "error", err)
		return nil, fmt.Errorf("failed to increment unsolved count: %w", err)
	}

	// Check if we exceeded max count
	if unsolvedCount > int64(s.maxUnsolvedChallenges) {
		s.logger.Warn("Too many unsolved challenges", "ip", clientIP, "count", unsolvedCount, "max", s.maxUnsolvedChallenges)
		return nil, fmt.Errorf("too many unsolved challenges")
	}

	// Generate challenge ID
	challengeID := uuid.New().String()

	// Get request count to determine difficulty
	requestCount, err := s.repo.GetAndIncrementRequestCount(ctx, clientIP)
	if err != nil {
		s.logger.Error("Failed to get request count", "ip", clientIP, "error", err)
		return nil, fmt.Errorf("failed to get request count: %w", err)
	}

	// Calculate difficulty level based on request count
	difficultyLevel := s.utils.CalculateDifficultyLevel(requestCount)
	s.logger.Debug("Calculated difficulty level", "ip", clientIP, "level", difficultyLevel, "requestCount", requestCount)

	// Generate random seed
	seed, err := s.utils.GenerateRandomSeed()
	if err != nil {
		s.logger.Error("Failed to generate random seed", "error", err)
		return nil, fmt.Errorf("failed to generate random seed: %w", err)
	}

	// Create task with challenge ID and seed
	task := powrepo.Task{
		ID:              challengeID,
		Seed:            seed,
		DifficultyLevel: difficultyLevel,
	}

	if err := s.repo.CreateTask(ctx, task, s.config.ChallengeTTL); err != nil {
		s.logger.Error("Failed to create task", "id", challengeID, "error", err)
		return nil, fmt.Errorf("failed to create challenge task: %w", err)
	}

	challenge := &protocol.Challenge{
		ID:              challengeID,
		Seed:            seed,
		DifficultyLevel: difficultyLevel,
		ScryptN:         s.config.ScryptN,
		ScryptR:         s.config.ScryptR,
		ScryptP:         s.config.ScryptP,
		KeyLen:          s.config.KeyLen,
	}

	s.logger.Debug("Challenge created", "id", challengeID, "difficulty", difficultyLevel)
	return challenge, nil
}

// ValidateChallenge validates a PoW challenge solution
func (s *Service) ValidateChallenge(ctx context.Context, clientIP, challengeID, nonce string) error {
	s.logger.Debug("Validating challenge", "ip", clientIP, "id", challengeID)

	// Validate UUID format
	if _, err := uuid.Parse(challengeID); err != nil {
		s.logger.Debug("Invalid challenge ID format", "id", challengeID)
		return fmt.Errorf("invalid challenge ID format")
	}

	// Get task from repository
	task, err := s.repo.GetTask(ctx, challengeID)
	if err != nil {
		s.logger.Debug("Invalid challenge", "id", challengeID, "error", err)
		return fmt.Errorf("invalid challenge: %w", err)
	}

	// Verify the solution meets the difficulty requirement
	if !s.utils.VerifySolution(challengeID, task.Seed, nonce, task.DifficultyLevel) {
		s.logger.Debug("Invalid challenge solution", "id", challengeID, "nonce", nonce)
		return fmt.Errorf("invalid challenge solution")
	}

	// Delete the task from storage to prevent replay attacks
	if err := s.repo.DeleteTask(ctx, challengeID); err != nil {
		s.logger.Error("Failed to delete task after validation", "id", challengeID, "error", err)
		// Continue execution even if deletion fails, as the validation was successful
	}

	// With 5% chance, decrement unsolved count by 2 instead of 1 (forgiving a missed job)
	if s.random.Float32() < 0.05 {
		s.logger.Debug("Lucky client! Decrementing unsolved count by 2", "ip", clientIP)
		if err := s.repo.DecrementUnsolvedCountBy(ctx, clientIP, 2); err != nil {
			s.logger.Error("Failed to decrement unsolved count by 2", "ip", clientIP, "error", err)
			return fmt.Errorf("failed to update challenge count: %w", err)
		}
	} else {
		// Regular case: decrement unsolved count by 1
		if err := s.repo.DecrementUnsolvedCount(ctx, clientIP); err != nil {
			s.logger.Error("Failed to decrement unsolved count", "ip", clientIP, "error", err)
			return fmt.Errorf("failed to update challenge count: %w", err)
		}
	}

	s.logger.Debug("Challenge validated successfully", "ip", clientIP, "id", challengeID)

	return nil
}
