package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opravdin/word-of-wisdom/client/internal/config"
	"github.com/opravdin/word-of-wisdom/client/internal/domain"
	"github.com/opravdin/word-of-wisdom/client/internal/http"
	"github.com/opravdin/word-of-wisdom/client/internal/logger"
	"github.com/opravdin/word-of-wisdom/client/internal/pow"
	"github.com/opravdin/word-of-wisdom/client/internal/tcp"
)

func main() {
	// Define command-line flags
	serverAddr := flag.String("server", "localhost:8080", "Word of Wisdom TCP server address")
	httpAddr := flag.String("http", "localhost:3000", "HTTP server address")
	cliMode := flag.Bool("cli", false, "Run in CLI mode instead of starting HTTP server")
	flag.Parse()

	// Print banner
	fmt.Println("Word of Wisdom Client")
	fmt.Println("=====================")

	// Create logger
	log := logger.NewStdLogger()

	// Load configuration
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Error("Invalid configuration", "error", err)
		os.Exit(1)
	}

	// Override config with command-line flags
	if *serverAddr != "localhost:8080" {
		cfg.ServerAddress = *serverAddr
	}
	if *httpAddr != "localhost:3000" {
		cfg.HTTPAddress = *httpAddr
	}

	// Create dependencies
	tcpFactory := tcp.NewTCPClientFactory(log, cfg)
	solverFactory := pow.NewSolverFactory(log, cfg)
	solver := solverFactory.NewSolver()

	if *cliMode {
		// Run in CLI mode (legacy behavior)
		runCLIMode(cfg.ServerAddress, tcpFactory, solver, log)
	} else {
		// Run HTTP server with web UI
		fmt.Printf("Starting HTTP server on %s\n", cfg.HTTPAddress)
		fmt.Printf("Connecting to Word of Wisdom server at %s\n", cfg.ServerAddress)

		// Create and start HTTP server
		quoteService := http.NewQuoteService(cfg.ServerAddress, tcpFactory, solver, log)
		server := http.NewServer(quoteService, log)

		// Handle graceful shutdown
		go func() {
			// Create a channel to listen for OS signals
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			// Wait for a signal
			sig := <-sigChan
			log.Info("Received signal, shutting down", "signal", sig)

			// Create a context with timeout for shutdown
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Shutdown the server
			if err := server.Stop(ctx); err != nil {
				log.Error("Error shutting down server", "error", err)
			}
		}()

		// Start the server (this will block until the server is stopped)
		if err := server.Start(cfg.HTTPAddress); err != nil {
			log.Error("Error starting HTTP server", "error", err)
			os.Exit(1)
		}
	}
}

// runCLIMode runs the client in CLI mode (legacy behavior)
func runCLIMode(serverAddr string, clientFactory tcp.TCPClientFactory, solver pow.Solver, log logger.Logger) {
	log.Info("Running in CLI mode", "server", serverAddr)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create TCP client
	client, err := clientFactory.NewClient(ctx, serverAddr)
	if err != nil {
		log.Error("Error connecting to server", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	// Request quote
	log.Info("Requesting quote...")

	// Send quote request
	requestMsg := domain.Message{
		Type: domain.TypeQuoteRequest,
	}

	if err := client.SendMessage(ctx, requestMsg); err != nil {
		log.Error("Error sending quote request", "error", err)
		os.Exit(1)
	}

	// Read challenge
	msg, err := client.ReadMessage(ctx)
	if err != nil {
		log.Error("Error reading challenge", "error", err)
		os.Exit(1)
	}

	// Check for error response
	if msg.Type == domain.TypeError {
		var errorData domain.ErrorData
		if err := json.Unmarshal(msg.Data, &errorData); err != nil {
			log.Error("Error parsing error data", "error", err)
			os.Exit(1)
		}
		log.Error("Server error", "code", errorData.Code, "message", errorData.Message)
		os.Exit(1)
	}

	// Verify message type
	if msg.Type != domain.TypePowChallenge {
		log.Error("Unexpected message type", "type", msg.Type)
		os.Exit(1)
	}

	// Process challenge
	log.Info("Solving PoW challenge...")

	// Parse challenge data to check if it contains PoW parameters
	var challengeData domain.PowChallengeData
	if err := json.Unmarshal(msg.Data, &challengeData); err != nil {
		log.Error("Error parsing challenge data", "error", err)
		os.Exit(1)
	}

	// Check if the server sent PoW parameters
	if challengeData.ScryptN > 0 && challengeData.ScryptR > 0 && challengeData.ScryptP > 0 && challengeData.KeyLen > 0 {
		log.Info("Using server-provided PoW parameters",
			"N", challengeData.ScryptN,
			"R", challengeData.ScryptR,
			"P", challengeData.ScryptP,
			"KeyLen", challengeData.KeyLen)

		if err := client.ProcessChallenge(ctx, msg, solver.Solve); err != nil {
			log.Error("Error processing challenge", "error", err)
			os.Exit(1)
		}
	} else {
		// Fallback to default parameters for backward compatibility
		log.Info("Server did not provide PoW parameters, using defaults")
		if err := client.ProcessChallengeWithDefaults(ctx, msg, solver.SolveWithDefaults); err != nil {
			log.Error("Error processing challenge", "error", err)
			os.Exit(1)
		}
	}

	// Get quote
	quote, err := client.GetQuote(ctx)
	if err != nil {
		log.Error("Error getting quote", "error", err)
		os.Exit(1)
	}

	// Display quote
	fmt.Printf("\nQuote: %s\n", quote.Text)
	fmt.Printf("Author: %s\n", quote.Author)
}
