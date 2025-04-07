package domain

import "encoding/json"

// Message types
const (
	TypeQuoteRequest  = "quote_request"
	TypePowChallenge  = "pow_challenge"
	TypePowSolution   = "pow_solution"
	TypeQuoteResponse = "quote_response"
	TypeError         = "error"
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

// Quote represents a wisdom quote
type Quote struct {
	Text   string
	Author string
}
