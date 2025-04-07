package http

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/opravdin/word-of-wisdom/client/internal/domain"
	"github.com/opravdin/word-of-wisdom/client/internal/logger"
	"github.com/opravdin/word-of-wisdom/client/internal/pow"
	"github.com/opravdin/word-of-wisdom/client/internal/tcp"
)

// serverImpl implements the Server interface
type serverImpl struct {
	quoteService QuoteService
	logger       logger.Logger
	server       *http.Server
}

// statsImpl implements the Stats struct with mutex for thread safety
type statsImpl struct {
	mu                     sync.Mutex
	RequestCount           int       `json:"requestCount"`
	SuccessCount           int       `json:"successCount"`
	FailureCount           int       `json:"failureCount"`
	LastDifficulty         int       `json:"lastDifficulty"`         // ScryptN parameter
	LastDifficultyLevel    int       `json:"lastDifficultyLevel"`    // Number of leading zeros required
	LastScryptR            int       `json:"lastScryptR"`            // ScryptR parameter
	LastScryptP            int       `json:"lastScryptP"`            // ScryptP parameter
	LastKeyLen             int       `json:"lastKeyLen"`             // Key length parameter
	EstimatedComplexity    float64   `json:"estimatedComplexity"`    // Estimated computational complexity
	AverageSolveTime       float64   `json:"averageSolveTime"`       // Average time to solve in seconds
	TotalSolveTime         float64   `json:"totalSolveTime"`         // Total time spent solving
	MinSolveTime           float64   `json:"minSolveTime"`           // Minimum solve time observed
	MaxSolveTime           float64   `json:"maxSolveTime"`           // Maximum solve time observed
	LastSolveTime          float64   `json:"lastSolveTime"`          // Last solve time
	LastRequestTime        time.Time `json:"lastRequestTime"`        // Time of last request
	LoadTestActive         bool      `json:"loadTestActive"`         // Whether load test is active
	LoadTestStartTime      time.Time `json:"loadTestStartTime"`      // When load test started
	LoadTestRequests       int       `json:"loadTestRequests"`       // Number of requests in load test
	LoadTestRequestsPerSec float64   `json:"loadTestRequestsPerSec"` // Requests per second during load test
}

// quoteServiceImpl implements the QuoteService interface
type quoteServiceImpl struct {
	serverAddr    string
	clientFactory tcp.TCPClientFactory
	solver        pow.Solver
	stats         *statsImpl
	logger        logger.Logger
}

// NewServer creates a new HTTP server
func NewServer(quoteService QuoteService, logger logger.Logger) Server {
	return &serverImpl{
		quoteService: quoteService,
		logger:       logger.With("component", "http_server"),
	}
}

// NewQuoteService creates a new quote service
func NewQuoteService(serverAddr string, clientFactory tcp.TCPClientFactory, solver pow.Solver, logger logger.Logger) QuoteService {
	return &quoteServiceImpl{
		serverAddr:    serverAddr,
		clientFactory: clientFactory,
		solver:        solver,
		stats: &statsImpl{
			LastRequestTime: time.Now(),
		},
		logger: logger.With("component", "quote_service"),
	}
}

// Start starts the HTTP server
func (s *serverImpl) Start(httpAddr string) error {
	s.logger.Info("Starting HTTP server", "address", httpAddr)

	// Serve static files
	fs := http.FileServer(http.Dir("./internal/http/static"))
	http.Handle("/", fs)

	// API endpoints
	http.HandleFunc("/api/quote", s.handleGetQuote)
	http.HandleFunc("/api/challenge", s.handleGetChallenge)
	http.HandleFunc("/api/stats", s.handleGetStats)
	http.HandleFunc("/api/load/start", s.handleStartLoad)
	http.HandleFunc("/api/load/stop", s.handleStopLoad)

	s.server = &http.Server{Addr: httpAddr}
	return s.server.ListenAndServe()
}

// Stop stops the HTTP server
func (s *serverImpl) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")
	return s.server.Shutdown(ctx)
}

// handleGetQuote handles the request to get a quote
func (s *serverImpl) handleGetQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.logger.Debug("Handling quote request")
	quote, challenge, err := s.quoteService.GetQuote(r.Context())
	if err != nil {
		s.sendErrorResponse(w, "Error getting quote: "+err.Error())
		return
	}

	// Send response
	response := QuoteResponse{
		Success:   true,
		Quote:     quote,
		Stats:     s.quoteService.GetStats(),
		Challenge: challenge,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("Error encoding response", "error", err)
	}
}

// handleGetChallenge handles the request to get a challenge without solving it
func (s *serverImpl) handleGetChallenge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.logger.Debug("Handling challenge request")
	challenge, err := s.quoteService.GetChallenge(r.Context())
	if err != nil {
		s.sendErrorResponse(w, "Error getting challenge: "+err.Error())
		return
	}

	// Send response
	response := QuoteResponse{
		Success:   true,
		Challenge: challenge,
		Stats:     s.quoteService.GetStats(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("Error encoding response", "error", err)
	}
}

// handleGetStats handles the request to get stats
func (s *serverImpl) handleGetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.logger.Debug("Handling stats request")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.quoteService.GetStats()); err != nil {
		s.logger.Error("Error encoding stats", "error", err)
	}
}

// handleStartLoad handles the request to start load testing
func (s *serverImpl) handleStartLoad(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.logger.Info("Starting load test")
	err := s.quoteService.StartLoadTest()
	if err != nil {
		s.sendErrorResponse(w, "Error starting load test: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
		s.logger.Error("Error encoding response", "error", err)
	}
}

// handleStopLoad handles the request to stop load testing
func (s *serverImpl) handleStopLoad(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.logger.Info("Stopping load test")
	err := s.quoteService.StopLoadTest()
	if err != nil {
		s.sendErrorResponse(w, "Error stopping load test: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
		s.logger.Error("Error encoding response", "error", err)
	}
}

// sendErrorResponse sends an error response
func (s *serverImpl) sendErrorResponse(w http.ResponseWriter, errMsg string) {
	s.logger.Error("Error response", "error", errMsg)

	response := QuoteResponse{
		Success: false,
		Error:   errMsg,
		Stats:   s.quoteService.GetStats(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("Error encoding error response", "error", err)
	}
}

// GetQuote gets a quote from the server
func (s *quoteServiceImpl) GetQuote(ctx context.Context) (*domain.Quote, *ChallengeInfo, error) {
	startTime := time.Now()
	s.stats.mu.Lock()
	s.stats.RequestCount++
	s.stats.LastRequestTime = startTime
	s.stats.mu.Unlock()

	s.logger.Debug("Getting quote from server", "serverAddr", s.serverAddr)

	// Create TCP client
	client, err := s.clientFactory.NewClient(ctx, s.serverAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to server: %w", err)
	}
	defer client.Close()

	// Send quote request
	requestMsg := domain.Message{
		Type: domain.TypeQuoteRequest,
	}

	if err := client.SendMessage(ctx, requestMsg); err != nil {
		return nil, nil, fmt.Errorf("error sending quote request: %w", err)
	}

	// Read challenge
	msg, err := client.ReadMessage(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading challenge: %w", err)
	}

	// Check for error response
	if msg.Type == domain.TypeError {
		var errorData domain.ErrorData
		if err := json.Unmarshal(msg.Data, &errorData); err != nil {
			return nil, nil, fmt.Errorf("error parsing error data: %w", err)
		}
		return nil, nil, fmt.Errorf("server error: %s - %s", errorData.Code, errorData.Message)
	}

	// Verify message type
	if msg.Type != domain.TypePowChallenge {
		return nil, nil, fmt.Errorf("unexpected message type: %s", msg.Type)
	}

	// Parse challenge data
	var challengeData domain.PowChallengeData
	if err := json.Unmarshal(msg.Data, &challengeData); err != nil {
		return nil, nil, fmt.Errorf("error parsing challenge data: %w", err)
	}

	// Create challenge info
	challenge := &ChallengeInfo{
		ChallengeID:     challengeData.ChallengeID,
		Task:            challengeData.Task,
		DifficultyLevel: challengeData.DifficultyLevel,
		ScryptN:         challengeData.ScryptN,
		ScryptR:         challengeData.ScryptR,
		ScryptP:         challengeData.ScryptP,
		KeyLen:          challengeData.KeyLen,
	}

	// Solve challenge
	var solveErr error
	solveStartTime := time.Now()

	if challengeData.ScryptN > 0 && challengeData.ScryptR > 0 && challengeData.ScryptP > 0 && challengeData.KeyLen > 0 {
		solveErr = client.ProcessChallenge(ctx, msg, s.solver.Solve)

		// Calculate estimated complexity based on Scrypt parameters
		// Higher values mean more computational work
		estimatedComplexity := float64(challengeData.ScryptN) * float64(challengeData.ScryptR) * float64(challengeData.ScryptP) *
			float64(challengeData.KeyLen) * math.Pow(16, float64(challengeData.DifficultyLevel))

		// Update stats with difficulty parameters
		s.stats.mu.Lock()
		s.stats.LastDifficulty = challengeData.ScryptN
		s.stats.LastDifficultyLevel = challengeData.DifficultyLevel
		s.stats.LastScryptR = challengeData.ScryptR
		s.stats.LastScryptP = challengeData.ScryptP
		s.stats.LastKeyLen = challengeData.KeyLen
		s.stats.EstimatedComplexity = estimatedComplexity
		challenge.EstComplexity = estimatedComplexity
		s.stats.mu.Unlock()
	} else {
		solveErr = client.ProcessChallengeWithDefaults(ctx, msg, s.solver.SolveWithDefaults)

		// Update stats with default difficulty
		s.stats.mu.Lock()
		s.stats.LastDifficulty = pow.DefaultScryptN
		s.stats.LastDifficultyLevel = 1 // Default difficulty level is 1
		s.stats.LastScryptR = pow.DefaultScryptR
		s.stats.LastScryptP = pow.DefaultScryptP
		s.stats.LastKeyLen = pow.DefaultScryptKeyLen
		// Calculate estimated complexity for default parameters
		estimatedComplexity := float64(pow.DefaultScryptN) * float64(pow.DefaultScryptR) *
			float64(pow.DefaultScryptP) * float64(pow.DefaultScryptKeyLen) * 16.0 // 16^1 for difficulty level 1
		s.stats.EstimatedComplexity = estimatedComplexity
		challenge.EstComplexity = estimatedComplexity
		s.stats.mu.Unlock()
	}

	if solveErr != nil {
		return nil, nil, fmt.Errorf("error solving challenge: %w", solveErr)
	}

	// Get quote
	quote, err := client.GetQuote(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting quote: %w", err)
	}

	// Update stats with solve time
	solveTime := time.Since(solveStartTime).Seconds()

	s.stats.mu.Lock()
	s.stats.SuccessCount++
	s.stats.TotalSolveTime += solveTime
	s.stats.LastSolveTime = solveTime

	// Update min/max solve times
	if s.stats.MinSolveTime == 0 || solveTime < s.stats.MinSolveTime {
		s.stats.MinSolveTime = solveTime
	}
	if solveTime > s.stats.MaxSolveTime {
		s.stats.MaxSolveTime = solveTime
	}

	// Calculate average
	s.stats.AverageSolveTime = s.stats.TotalSolveTime / float64(s.stats.SuccessCount)

	// Calculate requests per second during load test if active
	if s.stats.LoadTestActive && time.Since(s.stats.LoadTestStartTime).Seconds() > 0 {
		s.stats.LoadTestRequestsPerSec = float64(s.stats.LoadTestRequests) / time.Since(s.stats.LoadTestStartTime).Seconds()
	}

	s.stats.mu.Unlock()

	return quote, challenge, nil
}

// GetChallenge gets a challenge from the server without solving it
func (s *quoteServiceImpl) GetChallenge(ctx context.Context) (*ChallengeInfo, error) {
	s.logger.Debug("Getting challenge from server", "serverAddr", s.serverAddr)

	// Create TCP client
	client, err := s.clientFactory.NewClient(ctx, s.serverAddr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %w", err)
	}
	defer client.Close()

	// Send quote request
	requestMsg := domain.Message{
		Type: domain.TypeQuoteRequest,
	}

	if err := client.SendMessage(ctx, requestMsg); err != nil {
		return nil, fmt.Errorf("error sending quote request: %w", err)
	}

	// Read challenge
	msg, err := client.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading challenge: %w", err)
	}

	// Check for error response
	if msg.Type == domain.TypeError {
		var errorData domain.ErrorData
		if err := json.Unmarshal(msg.Data, &errorData); err != nil {
			return nil, fmt.Errorf("error parsing error data: %w", err)
		}
		return nil, fmt.Errorf("server error: %s - %s", errorData.Code, errorData.Message)
	}

	// Verify message type
	if msg.Type != domain.TypePowChallenge {
		return nil, fmt.Errorf("unexpected message type: %s", msg.Type)
	}

	// Parse challenge data
	var challengeData domain.PowChallengeData
	if err := json.Unmarshal(msg.Data, &challengeData); err != nil {
		return nil, fmt.Errorf("error parsing challenge data: %w", err)
	}

	// Create challenge info
	challenge := &ChallengeInfo{
		ChallengeID:     challengeData.ChallengeID,
		Task:            challengeData.Task,
		DifficultyLevel: challengeData.DifficultyLevel,
		ScryptN:         challengeData.ScryptN,
		ScryptR:         challengeData.ScryptR,
		ScryptP:         challengeData.ScryptP,
		KeyLen:          challengeData.KeyLen,
	}

	// If parameters are not provided, use defaults
	if challenge.ScryptN <= 0 {
		challenge.ScryptN = pow.DefaultScryptN
	}
	if challenge.ScryptR <= 0 {
		challenge.ScryptR = pow.DefaultScryptR
	}
	if challenge.ScryptP <= 0 {
		challenge.ScryptP = pow.DefaultScryptP
	}
	if challenge.KeyLen <= 0 {
		challenge.KeyLen = pow.DefaultScryptKeyLen
	}
	if challenge.DifficultyLevel <= 0 {
		challenge.DifficultyLevel = 1 // Default difficulty level is 1
	}

	// Calculate estimated complexity
	challenge.EstComplexity = float64(challenge.ScryptN) * float64(challenge.ScryptR) *
		float64(challenge.ScryptP) * float64(challenge.KeyLen) *
		math.Pow(16, float64(challenge.DifficultyLevel))

	return challenge, nil
}

// GetStats gets the current statistics
func (s *quoteServiceImpl) GetStats() *Stats {
	// Convert internal stats to the interface Stats type
	s.stats.mu.Lock()
	defer s.stats.mu.Unlock()

	return &Stats{
		RequestCount:           s.stats.RequestCount,
		SuccessCount:           s.stats.SuccessCount,
		FailureCount:           s.stats.FailureCount,
		LastDifficulty:         s.stats.LastDifficulty,
		LastDifficultyLevel:    s.stats.LastDifficultyLevel,
		LastScryptR:            s.stats.LastScryptR,
		LastScryptP:            s.stats.LastScryptP,
		LastKeyLen:             s.stats.LastKeyLen,
		EstimatedComplexity:    s.stats.EstimatedComplexity,
		AverageSolveTime:       s.stats.AverageSolveTime,
		TotalSolveTime:         s.stats.TotalSolveTime,
		MinSolveTime:           s.stats.MinSolveTime,
		MaxSolveTime:           s.stats.MaxSolveTime,
		LastSolveTime:          s.stats.LastSolveTime,
		LoadTestActive:         s.stats.LoadTestActive,
		LoadTestRequests:       s.stats.LoadTestRequests,
		LoadTestRequestsPerSec: s.stats.LoadTestRequestsPerSec,
	}
}

// StartLoadTest starts a load test
func (s *quoteServiceImpl) StartLoadTest() error {
	s.stats.mu.Lock()
	if s.stats.LoadTestActive {
		s.stats.mu.Unlock()
		return fmt.Errorf("load test already active")
	}

	s.stats.LoadTestActive = true
	s.stats.LoadTestStartTime = time.Now()
	s.stats.LoadTestRequests = 0
	s.stats.mu.Unlock()

	// Start load test in a goroutine
	go s.runLoadTest()
	return nil
}

// StopLoadTest stops a load test
func (s *quoteServiceImpl) StopLoadTest() error {
	s.stats.mu.Lock()
	s.stats.LoadTestActive = false
	s.stats.mu.Unlock()
	return nil
}

// runLoadTest runs a load test by continuously requesting quotes
func (s *quoteServiceImpl) runLoadTest() {
	s.logger.Info("Starting load test")

	for {
		s.stats.mu.Lock()
		active := s.stats.LoadTestActive
		s.stats.mu.Unlock()

		if !active {
			s.logger.Info("Load test stopped")
			return
		}

		// Create context for the request
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		// Get quote
		_, _, err := s.GetQuote(ctx)
		if err != nil {
			s.logger.Error("Load test error", "error", err)
		}

		cancel()

		// Update load test stats
		s.stats.mu.Lock()
		s.stats.LoadTestRequests++
		if time.Since(s.stats.LoadTestStartTime).Seconds() > 0 {
			s.stats.LoadTestRequestsPerSec = float64(s.stats.LoadTestRequests) / time.Since(s.stats.LoadTestStartTime).Seconds()
		}
		s.stats.mu.Unlock()
	}
}
