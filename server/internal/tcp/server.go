package tcp

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/opravdin/word-of-wisdom/internal/logger"
	"github.com/opravdin/word-of-wisdom/internal/protocol"
)

// Default timeout values
const (
	// ReadTimeout is the maximum duration for reading the entire request
	ReadTimeout = 30 * time.Second
	// WriteTimeout is the maximum duration for writing the entire response
	WriteTimeout = 30 * time.Second
	// IdleTimeout is the maximum amount of time to wait for the next request
	IdleTimeout = 60 * time.Second
)

// Server represents a TCP server
type Server struct {
	listener          net.Listener
	connectionHandler *ConnectionHandler
	connections       map[string]net.Conn
	mu                sync.Mutex
	logger            logger.Logger
}

// NewServer creates a new TCP server
func NewServer(
	listener net.Listener,
	log logger.Logger,
) *Server {
	// Create handler registry
	registry := NewHandlerRegistry()

	// Create connection handler
	connectionHandler := NewConnectionHandler(registry, log)

	return &Server{
		listener:          listener,
		connectionHandler: connectionHandler,
		connections:       make(map[string]net.Conn),
		logger:            log,
	}
}

// RegisterHandler registers a handler for a message type
func (s *Server) RegisterHandler(msgType string, handler protocol.MessageHandler) {
	s.connectionHandler.Registry.RegisterHandler(msgType, handler)
}

// Start starts the TCP server
func (s *Server) Start(ctx context.Context) error {
	defer s.listener.Close()

	// Create a goroutine to handle context cancellation
	go func() {
		<-ctx.Done()
		s.listener.Close()
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				s.logger.Error("Error accepting connection", "error", err)
				continue
			}
		}

		// Add connection to map
		clientAddr := conn.RemoteAddr().String()
		s.mu.Lock()
		s.connections[clientAddr] = conn
		s.mu.Unlock()

		// Handle connection in a goroutine
		go func(conn net.Conn) {
			defer func() {
				// Remove connection from map
				clientAddr := conn.RemoteAddr().String()
				s.mu.Lock()
				delete(s.connections, clientAddr)
				s.mu.Unlock()
			}()

			// Set initial read deadline
			if err := conn.SetReadDeadline(time.Now().Add(ReadTimeout)); err != nil {
				s.logger.Error("Failed to set read deadline", "addr", clientAddr, "error", err)
				conn.Close()
				return
			}

			// Set initial write deadline
			if err := conn.SetWriteDeadline(time.Now().Add(WriteTimeout)); err != nil {
				s.logger.Error("Failed to set write deadline", "addr", clientAddr, "error", err)
				conn.Close()
				return
			}

			s.connectionHandler.HandleConnection(ctx, conn)
		}(conn)
	}
}

// getIPFromAddr extracts the IP address from a net.Addr
func getIPFromAddr(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}
