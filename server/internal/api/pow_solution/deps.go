package powsolution

import (
	"context"

	"github.com/opravdin/word-of-wisdom/internal/domain"
)

//go:generate mockgen -destination=./mocks/mock_contracts.go -package=mocks github.com/opravdin/word-of-wisdom/internal/api/pow_solution PowService,QuoteUsecase

// powService defines the interface for the Proof of Work service
type powService interface {
	ValidateChallenge(ctx context.Context, clientIP, challengeID, solution string) error
}

// quoteUsecase defines the interface for the quote usecase
type quoteUsecase interface {
	Quote(ctx context.Context) domain.Quote
}
