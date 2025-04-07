package quoterequest

import (
	"context"

	"github.com/opravdin/word-of-wisdom/internal/protocol"
)

//go:generate mockgen -destination=./mocks/mock_contracts.go -package=mocks github.com/opravdin/word-of-wisdom/internal/api/quote_request PowService

// powService defines the interface for the Proof of Work service
type powService interface {
	CreateChallenge(ctx context.Context, clientIP string) (*protocol.Challenge, error)
}
