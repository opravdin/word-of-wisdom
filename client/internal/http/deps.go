package http

import (
	"context"

	"github.com/opravdin/word-of-wisdom/client/internal/domain"
)

//go:generate mockgen -source=${GOFILE} -package=mocks -destination=./mocks/deps.go

// Server defines the interface for the HTTP server
type Server interface {
	// Start starts the HTTP server
	Start(httpAddr string) error

	// Stop stops the HTTP server
	Stop(ctx context.Context) error
}

// QuoteService defines the interface for quote operations
type QuoteService interface {
	// GetQuote gets a quote from the server
	GetQuote(ctx context.Context) (*domain.Quote, *ChallengeInfo, error)

	// GetChallenge gets a challenge from the server without solving it
	GetChallenge(ctx context.Context) (*ChallengeInfo, error)

	// GetStats gets the current statistics
	GetStats() *Stats

	// StartLoadTest starts a load test
	StartLoadTest() error

	// StopLoadTest stops a load test
	StopLoadTest() error
}

// ChallengeInfo represents information about a PoW challenge
type ChallengeInfo struct {
	ChallengeID     string  `json:"challengeId"`
	Task            string  `json:"task"`
	DifficultyLevel int     `json:"difficultyLevel"`
	ScryptN         int     `json:"scryptN"`
	ScryptR         int     `json:"scryptR"`
	ScryptP         int     `json:"scryptP"`
	KeyLen          int     `json:"keyLen"`
	EstComplexity   float64 `json:"estimatedComplexity"`
}

// Stats represents statistics for the client
type Stats struct {
	RequestCount           int     `json:"requestCount"`
	SuccessCount           int     `json:"successCount"`
	FailureCount           int     `json:"failureCount"`
	LastDifficulty         int     `json:"lastDifficulty"`         // ScryptN parameter
	LastDifficultyLevel    int     `json:"lastDifficultyLevel"`    // Number of leading zeros required
	LastScryptR            int     `json:"lastScryptR"`            // ScryptR parameter
	LastScryptP            int     `json:"lastScryptP"`            // ScryptP parameter
	LastKeyLen             int     `json:"lastKeyLen"`             // Key length parameter
	EstimatedComplexity    float64 `json:"estimatedComplexity"`    // Estimated computational complexity
	AverageSolveTime       float64 `json:"averageSolveTime"`       // Average time to solve in seconds
	TotalSolveTime         float64 `json:"totalSolveTime"`         // Total time spent solving
	MinSolveTime           float64 `json:"minSolveTime"`           // Minimum solve time observed
	MaxSolveTime           float64 `json:"maxSolveTime"`           // Maximum solve time observed
	LastSolveTime          float64 `json:"lastSolveTime"`          // Last solve time
	LoadTestActive         bool    `json:"loadTestActive"`         // Whether load test is active
	LoadTestRequests       int     `json:"loadTestRequests"`       // Number of requests in load test
	LoadTestRequestsPerSec float64 `json:"loadTestRequestsPerSec"` // Requests per second during load test
}

// QuoteResponse represents a response with a quote
type QuoteResponse struct {
	Success   bool           `json:"success"`
	Quote     *domain.Quote  `json:"quote,omitempty"`
	Error     string         `json:"error,omitempty"`
	Stats     *Stats         `json:"stats,omitempty"`
	Challenge *ChallengeInfo `json:"challenge,omitempty"`
}
