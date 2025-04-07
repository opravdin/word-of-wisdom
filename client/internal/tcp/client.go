package tcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/opravdin/word-of-wisdom/client/internal/config"
	"github.com/opravdin/word-of-wisdom/client/internal/domain"
	"github.com/opravdin/word-of-wisdom/client/internal/logger"
)

// tcpClient implements the TCPClient interface
type tcpClient struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	logger logger.Logger
	config *config.ClientConfig
}

// tcpClientFactory implements the TCPClientFactory interface
type tcpClientFactory struct {
	logger logger.Logger
	config *config.ClientConfig
}

// NewTCPClientFactory creates a new TCP client factory
func NewTCPClientFactory(logger logger.Logger, config *config.ClientConfig) TCPClientFactory {
	return &tcpClientFactory{
		logger: logger,
		config: config,
	}
}

// NewClient creates a new TCP client with timeouts
func (f *tcpClientFactory) NewClient(ctx context.Context, serverAddr string) (TCPClient, error) {
	f.logger.Debug("Creating new TCP client", "server", serverAddr)

	// Create a dialer with timeout
	dialer := net.Dialer{
		Timeout: f.config.ConnectTimeout,
	}

	// Connect with timeout
	conn, err := dialer.DialContext(ctx, "tcp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %w", err)
	}

	// Set initial read and write deadlines
	if err := conn.SetReadDeadline(time.Now().Add(f.config.ReadTimeout)); err != nil {
		conn.Close()
		return nil, fmt.Errorf("error setting read deadline: %w", err)
	}

	if err := conn.SetWriteDeadline(time.Now().Add(f.config.WriteTimeout)); err != nil {
		conn.Close()
		return nil, fmt.Errorf("error setting write deadline: %w", err)
	}

	return &tcpClient{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		logger: f.logger.With("component", "tcp_client", "server", serverAddr),
		config: f.config,
	}, nil
}

// Close closes the client connection
func (c *tcpClient) Close() error {
	c.logger.Debug("Closing TCP connection")
	return c.conn.Close()
}

// SendMessage sends a message to the server
func (c *tcpClient) SendMessage(ctx context.Context, msg domain.Message) error {
	c.logger.Debug("Sending message", "type", msg.Type)

	// Use context for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Continue with the operation
	}

	// Set deadline based on context or default timeout
	deadline, ok := ctx.Deadline()
	if !ok {
		// No deadline in context, use default timeout
		deadline = time.Now().Add(c.config.WriteTimeout)
	}

	// Reset write deadline
	if err := c.conn.SetWriteDeadline(deadline); err != nil {
		return fmt.Errorf("error setting write deadline: %w", err)
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	// Add newline as message delimiter
	msgBytes = append(msgBytes, '\n')

	if _, err := c.writer.Write(msgBytes); err != nil {
		return fmt.Errorf("error writing to server: %w", err)
	}

	return c.writer.Flush()
}

// ReadMessage reads a message from the server
func (c *tcpClient) ReadMessage(ctx context.Context) (domain.Message, error) {
	c.logger.Debug("Reading message from server")

	// Use context for cancellation
	select {
	case <-ctx.Done():
		return domain.Message{}, ctx.Err()
	default:
		// Continue with the operation
	}

	// Set deadline based on context or default timeout
	deadline, ok := ctx.Deadline()
	if !ok {
		// No deadline in context, use default timeout
		deadline = time.Now().Add(c.config.ReadTimeout)
	}

	// Reset read deadline
	if err := c.conn.SetReadDeadline(deadline); err != nil {
		return domain.Message{}, fmt.Errorf("error setting read deadline: %w", err)
	}

	msgBytes, err := c.reader.ReadBytes('\n')
	if err != nil {
		return domain.Message{}, fmt.Errorf("error reading from server: %w", err)
	}

	var msg domain.Message
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return domain.Message{}, fmt.Errorf("error parsing message: %w", err)
	}

	c.logger.Debug("Received message", "type", msg.Type)
	return msg, nil
}

// ProcessChallenge processes a PoW challenge and sends the solution
func (c *tcpClient) ProcessChallenge(ctx context.Context, msg domain.Message, solveFn func(string, string, int, int, int, int, int) (string, error)) error {
	// Parse challenge
	var challengeData domain.PowChallengeData
	if err := json.Unmarshal(msg.Data, &challengeData); err != nil {
		return fmt.Errorf("error parsing challenge data: %w", err)
	}

	c.logger.Info("Processing PoW challenge",
		"challengeID", challengeData.ChallengeID,
		"difficultyLevel", challengeData.DifficultyLevel,
		"scryptN", challengeData.ScryptN,
		"scryptR", challengeData.ScryptR,
		"scryptP", challengeData.ScryptP,
		"keyLen", challengeData.KeyLen)

	// Extract PoW parameters from challenge
	difficultyLevel := challengeData.DifficultyLevel
	scryptN := challengeData.ScryptN
	scryptR := challengeData.ScryptR
	scryptP := challengeData.ScryptP
	keyLen := challengeData.KeyLen

	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Continue with the operation
	}

	// Solve challenge with parameters from server
	c.logger.Debug("Solving challenge")
	nonce, err := solveFn(challengeData.ChallengeID, challengeData.Seed, difficultyLevel, scryptN, scryptR, scryptP, keyLen)
	if err != nil {
		return fmt.Errorf("error solving challenge: %w", err)
	}

	c.logger.Debug("Challenge solved", "nonce", nonce)

	// Check context cancellation again
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Continue with the operation
	}

	// Send solution
	solutionData := domain.PowSolutionData{
		ChallengeID: challengeData.ChallengeID,
		Nonce:       nonce,
	}

	solutionDataBytes, err := json.Marshal(solutionData)
	if err != nil {
		return fmt.Errorf("error marshaling solution data: %w", err)
	}

	solutionMsg := domain.Message{
		Type: domain.TypePowSolution,
		Data: solutionDataBytes,
	}

	c.logger.Debug("Sending solution")
	return c.SendMessage(ctx, solutionMsg)
}

// ProcessChallengeWithDefaults processes a PoW challenge using default parameters
// This is for backward compatibility with older servers that don't send parameters
func (c *tcpClient) ProcessChallengeWithDefaults(ctx context.Context, msg domain.Message, solveFn func(string, string) (string, error)) error {
	// Parse challenge
	var challengeData domain.PowChallengeData
	if err := json.Unmarshal(msg.Data, &challengeData); err != nil {
		return fmt.Errorf("error parsing challenge data: %w", err)
	}

	c.logger.Info("Processing PoW challenge with default parameters",
		"challengeID", challengeData.ChallengeID)

	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Continue with the operation
	}

	// Solve challenge with default parameters
	c.logger.Debug("Solving challenge with defaults")
	nonce, err := solveFn(challengeData.ChallengeID, challengeData.Seed)
	if err != nil {
		return fmt.Errorf("error solving challenge: %w", err)
	}

	c.logger.Debug("Challenge solved", "nonce", nonce)

	// Check context cancellation again
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Continue with the operation
	}

	// Send solution
	solutionData := domain.PowSolutionData{
		ChallengeID: challengeData.ChallengeID,
		Nonce:       nonce,
	}

	solutionDataBytes, err := json.Marshal(solutionData)
	if err != nil {
		return fmt.Errorf("error marshaling solution data: %w", err)
	}

	solutionMsg := domain.Message{
		Type: domain.TypePowSolution,
		Data: solutionDataBytes,
	}

	c.logger.Debug("Sending solution")
	return c.SendMessage(ctx, solutionMsg)
}

// GetQuote reads and parses a quote response
func (c *tcpClient) GetQuote(ctx context.Context) (*domain.Quote, error) {
	c.logger.Debug("Reading quote response")

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Continue with the operation
	}

	// Read quote response
	msg, err := c.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading quote response: %w", err)
	}

	if msg.Type == domain.TypeError {
		var errorData domain.ErrorData
		if err := json.Unmarshal(msg.Data, &errorData); err != nil {
			return nil, fmt.Errorf("error parsing error data: %w", err)
		}
		return nil, fmt.Errorf("server error: %s - %s", errorData.Code, errorData.Message)
	}

	if msg.Type != domain.TypeQuoteResponse {
		return nil, fmt.Errorf("unexpected message type: %s", msg.Type)
	}

	// Parse quote
	var quoteData domain.QuoteResponseData
	if err := json.Unmarshal(msg.Data, &quoteData); err != nil {
		return nil, fmt.Errorf("error parsing quote data: %w", err)
	}

	c.logger.Info("Received quote", "author", quoteData.Author)

	return &domain.Quote{
		Text:   quoteData.Text,
		Author: quoteData.Author,
	}, nil
}
