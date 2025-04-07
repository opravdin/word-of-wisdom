package logger

import (
	"context"
	"fmt"
	"log"
	"os"
)

// stdLogger is a simple implementation of the Logger interface using the standard log package
type stdLogger struct {
	logger *log.Logger
	prefix string
}

// NewStdLogger creates a new standard logger
func NewStdLogger() Logger {
	return &stdLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Debug logs a debug message
func (l *stdLogger) Debug(msg string, keyvals ...interface{}) {
	l.log("DEBUG", msg, keyvals...)
}

// Info logs an info message
func (l *stdLogger) Info(msg string, keyvals ...interface{}) {
	l.log("INFO", msg, keyvals...)
}

// Warn logs a warning message
func (l *stdLogger) Warn(msg string, keyvals ...interface{}) {
	l.log("WARN", msg, keyvals...)
}

// Error logs an error message
func (l *stdLogger) Error(msg string, keyvals ...interface{}) {
	l.log("ERROR", msg, keyvals...)
}

// With returns a logger with the given key-value pairs
func (l *stdLogger) With(keyvals ...interface{}) Logger {
	newLogger := &stdLogger{
		logger: l.logger,
		prefix: l.prefix,
	}

	// Add key-value pairs to prefix
	if len(keyvals) > 0 {
		for i := 0; i < len(keyvals); i += 2 {
			key := keyvals[i]
			var value interface{} = "MISSING"
			if i+1 < len(keyvals) {
				value = keyvals[i+1]
			}
			newLogger.prefix += " " + key.(string) + "=" + stringify(value)
		}
	}

	return newLogger
}

// WithContext returns a logger with context
func (l *stdLogger) WithContext(ctx context.Context) Logger {
	// In a real implementation, we might extract trace IDs or other context values
	return l
}

// log formats and logs a message with the given level and key-value pairs
func (l *stdLogger) log(level, msg string, keyvals ...interface{}) {
	// Format key-value pairs
	kvString := ""
	for i := 0; i < len(keyvals); i += 2 {
		key := keyvals[i]
		var value interface{} = "MISSING"
		if i+1 < len(keyvals) {
			value = keyvals[i+1]
		}
		kvString += " " + key.(string) + "=" + stringify(value)
	}

	// Log the message with level, prefix, and key-value pairs
	l.logger.Printf("[%s]%s %s%s", level, l.prefix, msg, kvString)
}

// stringify converts a value to a string
func stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}

	switch v := value.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return fmt.Sprintf("%v", v)
	}
}
