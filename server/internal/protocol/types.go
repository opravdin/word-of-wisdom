package protocol

import (
	"context"
	"encoding/json"
	"net"
)

// Message types
const (
	TypeQuoteRequest  = "quote_request"
	TypePowChallenge  = "pow_challenge"
	TypePowSolution   = "pow_solution"
	TypeQuoteResponse = "quote_response"
	TypeError         = "error"
)

// Error codes
const (
	ErrInvalidRequest    = "invalid_request"
	ErrInvalidChallenge  = "invalid_challenge"
	ErrInvalidSolution   = "invalid_solution"
	ErrRateLimitExceeded = "rate_limit_exceeded"
	ErrInternalError     = "internal_error"
)

// Message represents a TCP protocol message
type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// PowChallengeData represents the data for a PoW challenge
type PowChallengeData struct {
	ChallengeID     string `json:"challenge_id"`
	Seed            string `json:"seed"`
	DifficultyLevel int    `json:"difficulty_level"`
	ScryptN         int    `json:"scrypt_n"`
	ScryptR         int    `json:"scrypt_r"`
	ScryptP         int    `json:"scrypt_p"`
	KeyLen          int    `json:"key_len"`
	Task            string `json:"task"`
}

// PowSolutionData represents the data for a PoW solution
type PowSolutionData struct {
	ChallengeID string `json:"challenge_id"`
	Nonce       string `json:"nonce"`
}

// QuoteResponseData represents the data for a quote response
type QuoteResponseData struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

// ErrorData represents the data for an error message
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// MessageHandler defines the interface for handling TCP messages
type MessageHandler interface {
	// HandleMessage handles a message and returns a response
	HandleMessage(ctx context.Context, conn net.Conn, clientIP string, msg Message) error
}

// SendMessage sends a message to the client
func SendMessage(conn net.Conn, msgType string, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := Message{
		Type: msgType,
		Data: dataBytes,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Add newline as message delimiter
	msgBytes = append(msgBytes, '\n')

	_, err = conn.Write(msgBytes)
	return err
}

// SendError sends an error message to the client
func SendError(conn net.Conn, code, message string) error {
	errorData := ErrorData{
		Code:    code,
		Message: message,
	}
	return SendMessage(conn, TypeError, errorData)
}
