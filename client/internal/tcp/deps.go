package tcp

import (
	"context"

	"github.com/opravdin/word-of-wisdom/client/internal/domain"
)

//go:generate mockgen -source=${GOFILE} -package=mocks -destination=./mocks/deps.go

// TCPClient defines the interface for TCP client operations
type TCPClient interface {
	// SendMessage sends a message to the server
	SendMessage(ctx context.Context, msg domain.Message) error

	// ReadMessage reads a message from the server
	ReadMessage(ctx context.Context) (domain.Message, error)

	// ProcessChallenge processes a PoW challenge and sends the solution
	ProcessChallenge(ctx context.Context, msg domain.Message, solveFn func(string, string, int, int, int, int, int) (string, error)) error

	// ProcessChallengeWithDefaults processes a PoW challenge using default parameters
	ProcessChallengeWithDefaults(ctx context.Context, msg domain.Message, solveFn func(string, string) (string, error)) error

	// GetQuote reads and parses a quote response
	GetQuote(ctx context.Context) (*domain.Quote, error)

	// Close closes the client connection
	Close() error
}

// TCPClientFactory defines the interface for creating TCP clients
type TCPClientFactory interface {
	// NewClient creates a new TCP client
	NewClient(ctx context.Context, serverAddr string) (TCPClient, error)
}
