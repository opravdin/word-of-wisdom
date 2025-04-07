package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// SlogLogger is an implementation of Logger using slog
type SlogLogger struct {
	logger *slog.Logger
}

// NewSlogLogger creates a new SlogLogger with the given options
func NewSlogLogger(opts ...SlogOption) *SlogLogger {
	config := defaultSlogConfig()

	for _, opt := range opts {
		opt(config)
	}

	var handler slog.Handler

	switch config.format {
	case FormatJSON:
		handler = slog.NewJSONHandler(config.output, &slog.HandlerOptions{
			Level: config.level,
		})
	default:
		handler = slog.NewTextHandler(config.output, &slog.HandlerOptions{
			Level: config.level,
		})
	}

	return &SlogLogger{
		logger: slog.New(handler),
	}
}

// Debug logs a debug message
func (l *SlogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Info logs an info message
func (l *SlogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Warn logs a warning message
func (l *SlogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Error logs an error message
func (l *SlogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// With returns a logger with the given attributes
func (l *SlogLogger) With(args ...any) Logger {
	return &SlogLogger{
		logger: l.logger.With(args...),
	}
}

// WithContext returns a logger with context
func (l *SlogLogger) WithContext(ctx context.Context) Logger {
	// Extract any context values that should be added to the logger
	// For now, we just return the same logger
	return l
}

// Format represents the output format for the logger
type Format string

const (
	// FormatText outputs logs in a human-readable text format
	FormatText Format = "text"
	// FormatJSON outputs logs in JSON format
	FormatJSON Format = "json"
)

// slogConfig holds configuration for the SlogLogger
type slogConfig struct {
	level  slog.Level
	format Format
	output io.Writer
}

// defaultSlogConfig returns the default configuration for SlogLogger
func defaultSlogConfig() *slogConfig {
	return &slogConfig{
		level:  slog.LevelInfo,
		format: FormatText,
		output: os.Stdout,
	}
}

// SlogOption is a function that configures a SlogLogger
type SlogOption func(*slogConfig)

// WithLevel sets the minimum log level
func WithLevel(level slog.Level) SlogOption {
	return func(c *slogConfig) {
		c.level = level
	}
}

// WithFormat sets the output format
func WithFormat(format Format) SlogOption {
	return func(c *slogConfig) {
		c.format = format
	}
}

// WithOutput sets the output writer
func WithOutput(output io.Writer) SlogOption {
	return func(c *slogConfig) {
		c.output = output
	}
}
