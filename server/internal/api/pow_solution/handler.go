package powsolution

import (
	"context"
	"encoding/json"
	"net"

	"github.com/opravdin/word-of-wisdom/internal/logger"
	"github.com/opravdin/word-of-wisdom/internal/protocol"
)

// Handler handles PoW solution messages
type Handler struct {
	powService   powService
	quoteUsecase quoteUsecase
	logger       logger.Logger
}

// NewHandler creates a new PoW solution handler
func NewHandler(powService powService, quoteUsecase quoteUsecase, log logger.Logger) *Handler {
	return &Handler{
		powService:   powService,
		quoteUsecase: quoteUsecase,
		logger:       log,
	}
}

// HandleMessage handles a PoW solution message
func (h *Handler) HandleMessage(ctx context.Context, conn net.Conn, clientIP string, msg protocol.Message) error {
	// Parse solution data
	var solutionData protocol.PowSolutionData
	if err := json.Unmarshal(msg.Data, &solutionData); err != nil {
		h.logger.Error("Error parsing solution data", "ip", clientIP, "error", err)
		return protocol.SendError(conn, protocol.ErrInvalidRequest, "Invalid solution format")
	}

	// Validate solution
	h.logger.Debug("Validating solution", "ip", clientIP, "id", solutionData.ChallengeID)
	err := h.powService.ValidateChallenge(ctx, clientIP, solutionData.ChallengeID, solutionData.Nonce)
	if err != nil {
		h.logger.Error("Error validating solution", "ip", clientIP, "id", solutionData.ChallengeID, "error", err)
		return protocol.SendError(conn, protocol.ErrInvalidSolution, err.Error())
	}

	// Get quote
	h.logger.Debug("Getting quote for client", "ip", clientIP)
	quote := h.quoteUsecase.Quote(ctx)

	// Send quote to client
	quoteData := protocol.QuoteResponseData{
		Text:   quote.Text,
		Author: quote.Author,
	}

	h.logger.Debug("Sending quote to client", "ip", clientIP, "author", quote.Author)
	return protocol.SendMessage(conn, protocol.TypeQuoteResponse, quoteData)
}
