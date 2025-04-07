package logger

import (
	"context"
)

//go:generate mockgen -source=${GOFILE} -package=mocks -destination=./mocks/repository.go

// Logger defines the interface for logging
type Logger interface {
	// Debug logs a debug message
	Debug(msg string, args ...any)
	// Info logs an info message
	Info(msg string, args ...any)
	// Warn logs a warning message
	Warn(msg string, args ...any)
	// Error logs an error message
	Error(msg string, args ...any)
	// With returns a logger with the given attributes
	With(args ...any) Logger
	// WithContext returns a logger with context
	WithContext(ctx context.Context) Logger
}
