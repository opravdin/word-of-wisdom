package pow

import (
	"context"
	"time"

	powrepo "github.com/opravdin/word-of-wisdom/internal/repository/pow"
)

//go:generate mockgen -source=${GOFILE} -package=mocks -destination=./mocks/deps.go

// Repository defines the interface for PoW repository operations
type Repository interface {
	GetAndIncrementRequestCount(ctx context.Context, ip string) (int64, error)
	CreateTask(ctx context.Context, task powrepo.Task, ttl time.Duration) error
	GetTask(ctx context.Context, taskID string) (*powrepo.Task, error)
	DeleteTask(ctx context.Context, taskID string) error
	IncrementUnsolvedCount(ctx context.Context, ip string) (int64, error)
	DecrementUnsolvedCount(ctx context.Context, ip string) error
	DecrementUnsolvedCountBy(ctx context.Context, ip string, count int) error
	GetUnsolvedCount(ctx context.Context, ip string) (int64, error)
}

// RandomProvider defines the interface for random number generation
type RandomProvider interface {
	Float32() float32
}

// PoWUtilsInterface defines the interface for PoW utility functions
type PoWUtilsInterface interface {
	CalculateDifficultyLevel(requestCount int64) int
	GenerateRandomSeed() (string, error)
	VerifySolution(challengeID, seed, nonce string, difficultyLevel int) bool
}

// Ensure PoWUtils implements PoWUtilsInterface
var _ PoWUtilsInterface = (*PoWUtils)(nil)
