package logger

import "context"

//go:generate mockgen -source=${GOFILE} -package=mocks -destination=./mocks/deps.go

// Logger defines the interface for logging operations
type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	With(keyvals ...interface{}) Logger
	WithContext(ctx context.Context) Logger
}
