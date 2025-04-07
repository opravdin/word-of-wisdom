package pow

import (
	"context"
	"errors"
	"time"
)

// Repository errors

var (
	ErrTaskAlreadyExists = errors.New("task already exists")
	ErrTaskNotFound      = errors.New("task not found")
	ErrRedisOperation    = errors.New("redis operation failed")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

type Task struct {
	ID              string
	Seed            string
	DifficultyLevel int
	Result          string // Keep for backward compatibility
}

type Repository interface {
	GetAndIncrementRequestCount(ctx context.Context, ip string) (int64, error)
	CreateTask(ctx context.Context, task Task, ttl time.Duration) error
	GetTask(ctx context.Context, taskID string) (*Task, error)
	DeleteTask(ctx context.Context, taskID string) error
	IncrementUnsolvedCount(ctx context.Context, ip string) (int64, error)
	DecrementUnsolvedCount(ctx context.Context, ip string) error
	DecrementUnsolvedCountBy(ctx context.Context, ip string, count int) error
	GetUnsolvedCount(ctx context.Context, ip string) (int64, error)
}
