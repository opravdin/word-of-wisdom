package tcp

import (
	"context"
	"net"

	"github.com/opravdin/word-of-wisdom/internal/protocol"
)

//go:generate mockgen -destination=./mocks/mock_handler.go -package=mocks github.com/opravdin/word-of-wisdom/internal/tcp MessageHandler

// MessageHandler defines the interface for handling TCP messages
type MessageHandler interface {
	// HandleMessage handles a message and returns a response
	HandleMessage(ctx context.Context, conn net.Conn, clientIP string, msg protocol.Message) error
}

// HandlerRegistry maps message types to handlers
type HandlerRegistry struct {
	handlers map[string]MessageHandler
}

// NewHandlerRegistry creates a new handler registry
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]MessageHandler),
	}
}

// RegisterHandler registers a handler for a message type
func (r *HandlerRegistry) RegisterHandler(msgType string, handler MessageHandler) {
	r.handlers[msgType] = handler
}

// GetHandler returns the handler for a message type
func (r *HandlerRegistry) GetHandler(msgType string) (MessageHandler, bool) {
	handler, ok := r.handlers[msgType]
	return handler, ok
}
