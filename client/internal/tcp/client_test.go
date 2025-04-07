package tcp

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/opravdin/word-of-wisdom/client/internal/config"
	"github.com/opravdin/word-of-wisdom/client/internal/domain"
	"github.com/opravdin/word-of-wisdom/client/internal/logger"
	"github.com/stretchr/testify/assert"
)

// Define constants for test values
const (
	longTimeout       = 10 * time.Second
	shortTimeout      = 10 * time.Millisecond
	veryShortTimeout  = 1 * time.Millisecond
	waitForExpiration = 20 * time.Millisecond
	testChallengeID   = "test-id"
	testSeed          = "test-seed"
	testDifficulty    = 1
	testScryptN       = 16384
	testScryptR       = 8
	testScryptP       = 1
	testKeyLen        = 32
	testQuoteText     = "Test quote"
	testQuoteAuthor   = "Test author"
)

// mockServer is a simple TCP server for testing
type mockServer struct {
	listener net.Listener
	done     chan struct{}
	stopped  bool
}

// newMockServer creates a new mock TCP server
func newMockServer(t *testing.T) (*mockServer, string) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	server := &mockServer{
		listener: listener,
		done:     make(chan struct{}),
		stopped:  false,
	}

	// Start the server in a goroutine
	go server.serve(t)

	return server, listener.Addr().String()
}

// serve handles incoming connections
func (s *mockServer) serve(t *testing.T) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// Check if the server is shutting down
			if s.stopped {
				return
			}
			t.Logf("Error accepting connection: %v", err)
			continue
		}

		// Handle the connection in a goroutine
		go s.handleConnection(t, conn)
	}
}

// handleConnection handles a single connection
func (s *mockServer) handleConnection(t *testing.T, conn net.Conn) {
	defer conn.Close()

	// Read the request
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Logf("Error reading from connection: %v", err)
		return
	}

	// Parse the request
	var msg domain.Message
	if err := json.Unmarshal(buf[:n-1], &msg); err != nil {
		t.Logf("Error parsing message: %v", err)
		return
	}

	// Send a response based on the request type
	switch msg.Type {
	case domain.TypeQuoteRequest:
		// Send a challenge
		challenge := domain.Message{
			Type: domain.TypePowChallenge,
			Data: json.RawMessage(
				`{"challenge_id":"` + testChallengeID + `","seed":"` + testSeed + `",` +
					`"difficulty_level":` + strconv.Itoa(testDifficulty) + `,` +
					`"scrypt_n":` + strconv.Itoa(testScryptN) + `,` +
					`"scrypt_r":` + strconv.Itoa(testScryptR) + `,` +
					`"scrypt_p":` + strconv.Itoa(testScryptP) + `,` +
					`"key_len":` + strconv.Itoa(testKeyLen) + `}`),
		}
		challengeBytes, _ := json.Marshal(challenge)
		challengeBytes = append(challengeBytes, '\n')
		if _, err := conn.Write(challengeBytes); err != nil {
			t.Logf("Error writing challenge: %v", err)
			return
		}
	case domain.TypePowSolution:
		// Send a quote
		quote := domain.Message{
			Type: domain.TypeQuoteResponse,
			Data: json.RawMessage(`{"text":"` + testQuoteText + `","author":"` + testQuoteAuthor + `"}`),
		}
		quoteBytes, _ := json.Marshal(quote)
		quoteBytes = append(quoteBytes, '\n')
		if _, err := conn.Write(quoteBytes); err != nil {
			t.Logf("Error writing quote: %v", err)
			return
		}
	}
}

// stop stops the mock server
func (s *mockServer) stop() {
	if !s.stopped {
		s.stopped = true
		close(s.done)
		s.listener.Close()
	}
}

// TestTCPClient tests the TCP client functionality
func TestTCPClient(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() (context.Context, context.CancelFunc)
		waitFunc    func()
		shouldFail  bool
		expectedErr error
	}{
		{
			name: "should_respect_context_timeout",
			setupFunc: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), shortTimeout)
			},
			waitFunc: func() {
				time.Sleep(waitForExpiration)
			},
			shouldFail:  true,
			expectedErr: context.DeadlineExceeded,
		},
		{
			name: "should_respect_context_cancellation",
			setupFunc: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				// Cancel immediately
				cancel()
				return ctx, cancel
			},
			waitFunc: func() {
				// No need to wait
			},
			shouldFail:  true,
			expectedErr: context.Canceled,
		},
		{
			name: "should_propagate_context_deadline",
			setupFunc: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), veryShortTimeout)
			},
			waitFunc: func() {
				time.Sleep(waitForExpiration)
			},
			shouldFail:  true,
			expectedErr: context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			server, addr := newMockServer(t)
			defer server.stop()

			log := logger.NewStdLogger()

			cfg := &config.ClientConfig{
				ConnectTimeout: longTimeout,
				ReadTimeout:    longTimeout,
				WriteTimeout:   longTimeout,
			}

			factory := NewTCPClientFactory(log, cfg)

			// Create a client with a background context
			client, err := factory.NewClient(context.Background(), addr)
			assert.NoError(t, err)
			defer client.Close()

			// Create the test context
			ctx, cancel := tt.setupFunc()
			defer cancel()

			// Wait if needed
			tt.waitFunc()

			// Execute
			err = client.SendMessage(ctx, domain.Message{Type: domain.TypeQuoteRequest})

			// Assert
			if tt.shouldFail {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestTCPClient_ReadMessage tests the TCPClient's ReadMessage method
func TestTCPClient_ReadMessage(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() (context.Context, context.CancelFunc)
		shouldFail  bool
		expectedErr error
	}{
		{
			name: "should_read_message_successfully",
			setupFunc: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
			shouldFail: false,
		},
		{
			name: "should_return_error_when_context_canceled",
			setupFunc: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Immediately cancel
				return ctx, cancel
			},
			shouldFail:  true,
			expectedErr: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			server, addr := newMockServer(t)
			defer server.stop()

			log := logger.NewStdLogger()

			cfg := &config.ClientConfig{
				ConnectTimeout: longTimeout,
				ReadTimeout:    longTimeout,
				WriteTimeout:   longTimeout,
			}

			factory := NewTCPClientFactory(log, cfg)

			// Create a client with a background context
			client, err := factory.NewClient(context.Background(), addr)
			assert.NoError(t, err)
			defer client.Close()

			// Create the test context
			ctx, cancel := tt.setupFunc()
			defer cancel()

			// First send a message to get a response
			if !tt.shouldFail {
				err = client.SendMessage(context.Background(), domain.Message{Type: domain.TypeQuoteRequest})
				assert.NoError(t, err)
			}

			// Execute
			msg, err := client.ReadMessage(ctx)

			// Assert
			if tt.shouldFail {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, domain.TypePowChallenge, msg.Type)
			}
		})
	}
}

// TestTCPClient_ProcessChallenge tests the TCPClient's ProcessChallenge method
func TestTCPClient_ProcessChallenge(t *testing.T) {
	// Create a mock server
	server, addr := newMockServer(t)
	defer server.stop()

	// Create a logger
	log := logger.NewStdLogger()

	// Create a client config
	cfg := &config.ClientConfig{
		ConnectTimeout: longTimeout,
		ReadTimeout:    longTimeout,
		WriteTimeout:   longTimeout,
	}

	// Create a client factory
	factory := NewTCPClientFactory(log, cfg)

	// Create a client
	client, err := factory.NewClient(context.Background(), addr)
	assert.NoError(t, err)
	defer client.Close()

	// Send a message to get a challenge
	err = client.SendMessage(context.Background(), domain.Message{Type: domain.TypeQuoteRequest})
	assert.NoError(t, err)

	// Read the challenge
	msg, err := client.ReadMessage(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, domain.TypePowChallenge, msg.Type)

	// Mock solve function
	solveFn := func(challengeID, seed string, difficultyLevel, scryptN, scryptR, scryptP, keyLen int) (string, error) {
		return "test-nonce", nil
	}

	// Process the challenge
	err = client.ProcessChallenge(context.Background(), msg, solveFn)
	assert.NoError(t, err)
}
