package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	powsolution "github.com/opravdin/word-of-wisdom/internal/api/pow_solution"
	quoterequest "github.com/opravdin/word-of-wisdom/internal/api/quote_request"
	"github.com/opravdin/word-of-wisdom/internal/configuration"
	"github.com/opravdin/word-of-wisdom/internal/configuration/env"
	"github.com/opravdin/word-of-wisdom/internal/infrastructure"
	"github.com/opravdin/word-of-wisdom/internal/logger"
	"github.com/opravdin/word-of-wisdom/internal/protocol"
	"github.com/opravdin/word-of-wisdom/internal/tcp"
)

func main() {
	// Initialize logger
	appLogger := logger.NewSlogLogger(
		logger.WithLevel(slog.LevelInfo),
		logger.WithFormat(logger.FormatJSON),
	)
	appLogger.Info("Starting Word of Wisdom server")

	// Initialize environment configuration
	appConfig := env.LoadFromEnv()
	appLogger.Info("Configuration loaded", "port", appConfig.Server.Port)

	// Initialize infrastructure
	storage, err := infrastructure.NewStorageConfiguration(appConfig.Redis, appLogger)
	if err != nil {
		appLogger.Error("Failed to initialize storage", "error", err)
		os.Exit(1)
	}
	defer storage.Close()

	// Initialize configuration
	config := configuration.NewDefaultConfiguration(*storage, appConfig, appLogger)

	// Get powService from configuration
	powService := config.PowService

	// Initialize TCP server
	port := appConfig.Server.Port

	// Create TCP listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		appLogger.Error("Failed to create TCP listener", "error", err)
		os.Exit(1)
	}

	// Create TCP server
	server := tcp.NewServer(listener, appLogger)

	// Create handlers
	quoteRequestHandler := quoterequest.NewHandler(powService, appLogger)
	powSolutionHandler := powsolution.NewHandler(powService, config.GetQuote, appLogger)

	// Register handlers
	server.RegisterHandler(protocol.TypeQuoteRequest, quoteRequestHandler)
	server.RegisterHandler(protocol.TypePowSolution, powSolutionHandler)

	// Start server in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		appLogger.Info("Starting TCP server", "port", port)
		if err := server.Start(ctx); err != nil {
			appLogger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Graceful shutdown
	appLogger.Info("Shutting down server...")
	cancel()
	if err := listener.Close(); err != nil {
		appLogger.Error("Error during server shutdown", "error", err)
	}
}
