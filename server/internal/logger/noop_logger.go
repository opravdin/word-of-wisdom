package logger

import (
	"context"
)

// NoopLogger is a logger that does nothing
type NoopLogger struct{}

// NewNoopLogger creates a new NoopLogger
func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

// Debug logs a debug message
func (l *NoopLogger) Debug(msg string, args ...any) {}

// Info logs an info message
func (l *NoopLogger) Info(msg string, args ...any) {}

// Warn logs a warning message
func (l *NoopLogger) Warn(msg string, args ...any) {}

// Error logs an error message
func (l *NoopLogger) Error(msg string, args ...any) {}

// With returns a logger with the given attributes
func (l *NoopLogger) With(args ...any) Logger {
	return l
}

// WithContext returns a logger with context
func (l *NoopLogger) WithContext(ctx context.Context) Logger {
	return l
}
