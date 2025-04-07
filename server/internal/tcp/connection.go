package tcp

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/opravdin/word-of-wisdom/internal/logger"
	"github.com/opravdin/word-of-wisdom/internal/protocol"
)

// ConnectionHandler handles client connections
type ConnectionHandler struct {
	Registry *HandlerRegistry
	logger   logger.Logger
}

// NewConnectionHandler creates a new connection handler
func NewConnectionHandler(registry *HandlerRegistry, log logger.Logger) *ConnectionHandler {
	return &ConnectionHandler{
		Registry: registry,
		logger:   log,
	}
}

// HandleConnection handles a client connection
func (h *ConnectionHandler) HandleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	clientAddr := conn.RemoteAddr().String()
	clientIP := getIPFromAddr(clientAddr)

	h.logger.Info("New connection", "addr", clientAddr)

	for {
		select {
		case <-ctx.Done():
			h.logger.Debug("Connection handler context done", "addr", clientAddr)
			return
		default:
			// Reset read deadline for this operation
			if err := conn.SetReadDeadline(time.Now().Add(ReadTimeout)); err != nil {
				h.logger.Error("Failed to set read deadline", "addr", clientAddr, "error", err)
				return
			}

			// Read message
			msgBytes, err := reader.ReadBytes('\n')
			if err != nil {
				h.logger.Error("Error reading from client", "addr", clientAddr, "error", err)
				return
			}

			// Reset write deadline for the response
			if err := conn.SetWriteDeadline(time.Now().Add(WriteTimeout)); err != nil {
				h.logger.Error("Failed to set write deadline", "addr", clientAddr, "error", err)
				return
			}

			// Parse message
			var msg protocol.Message
			if err := json.Unmarshal(msgBytes, &msg); err != nil {
				if err := protocol.SendError(conn, protocol.ErrInvalidRequest, "Invalid JSON format"); err != nil {
					h.logger.Error("Failed to send error response", "addr", clientAddr, "error", err)
					return
				}
				continue
			}

			// Get handler for message type
			handler, ok := h.Registry.GetHandler(msg.Type)
			if !ok {
				if err := protocol.SendError(conn, protocol.ErrInvalidRequest, "Unknown message type"); err != nil {
					h.logger.Error("Failed to send error response", "addr", clientAddr, "error", err)
					return
				}
				continue
			}

			// Handle message
			if err := handler.HandleMessage(ctx, conn, clientIP, msg); err != nil {
				h.logger.Error("Error handling message", "addr", clientAddr, "type", msg.Type, "error", err)
				// Error already sent to client by handler
			}
		}
	}
}
