package quoterequest

import (
	"context"
	"net"

	"github.com/opravdin/word-of-wisdom/internal/logger"
	"github.com/opravdin/word-of-wisdom/internal/protocol"
)

// Handler handles quote request messages
type Handler struct {
	powService powService
	logger     logger.Logger
}

// NewHandler creates a new quote request handler
func NewHandler(powService powService, log logger.Logger) *Handler {
	return &Handler{
		powService: powService,
		logger:     log,
	}
}

// HandleMessage handles a quote request message
func (h *Handler) HandleMessage(ctx context.Context, conn net.Conn, clientIP string, msg protocol.Message) error {
	// Create a new challenge
	h.logger.Debug("Creating challenge for client", "ip", clientIP)
	challenge, err := h.powService.CreateChallenge(ctx, clientIP)
	if err != nil {
		h.logger.Error("Error creating challenge", "ip", clientIP, "error", err)
		return protocol.SendError(conn, protocol.ErrRateLimitExceeded, err.Error())
	}

	// Send challenge to client
	challengeData := protocol.PowChallengeData{
		ChallengeID:     challenge.ID,
		Seed:            challenge.Seed,
		DifficultyLevel: challenge.DifficultyLevel,
		ScryptN:         challenge.ScryptN,
		ScryptR:         challenge.ScryptR,
		ScryptP:         challenge.ScryptP,
		KeyLen:          challenge.KeyLen,
		Task:            challenge.ID, // Use challenge ID as the task
	}

	h.logger.Debug("Sending challenge to client", "ip", clientIP, "id", challenge.ID)
	return protocol.SendMessage(conn, protocol.TypePowChallenge, challengeData)
}
