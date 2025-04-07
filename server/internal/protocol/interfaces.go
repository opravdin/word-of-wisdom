package protocol

import (
	"context"

	"github.com/opravdin/word-of-wisdom/internal/domain"
)

// Challenge represents a PoW challenge
type Challenge struct {
	ID              string
	Seed            string
	DifficultyLevel int
	Result          string
	ScryptN         int
	ScryptR         int
	ScryptP         int
	KeyLen          int
}

// PowService defines the interface for the Proof of Work service
type PowService interface {
	CreateChallenge(ctx context.Context, clientIP string) (*Challenge, error)
	ValidateChallenge(ctx context.Context, clientIP, challengeID, solution string) error
}

// PowMiddleware is an alias for PowService for backward compatibility
type PowMiddleware = PowService

// QuoteUsecase defines the interface for the quote usecase
type QuoteUsecase interface {
	Quote(ctx context.Context) domain.Quote
}
