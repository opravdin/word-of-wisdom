# Code Generation Rules

This document outlines coding standards for the project. All AI code generation agents must comply with these rules when generating code.

## Project Structure

- `cmd/`: Application entry points
- `internal/`: Private application code
  - Organize by feature/domain, not by technical layer
  - Use clear package boundaries with well-defined interfaces

## Package Organization

- Define interfaces in `deps.go` files
- Generate mocks in `mocks/` subdirectory
- Use dependency injection pattern
- Keep implementation details private (unexported)

## Code Style

- Use `CamelCase` for exported identifiers
- Use `camelCase` for unexported identifiers
- Keep functions focused on a single responsibility
- Document exported functions, types, and constants
- Group related constants in `const` blocks

## Dependency Management

- Define interfaces in `deps.go` files
- Generate mocks with: `//go:generate mockgen -source=${GOFILE} -package=mocks -destination=./mocks/deps.go`
- Pass dependencies via constructor functions
- Avoid global state and singletons

## Error Handling

- Use `fmt.Errorf("context: %w", err)` to wrap errors
- Log errors with context: `logger.Error("Failed to create task", "id", id, "error", err)`
- Return wrapped errors: `return nil, fmt.Errorf("failed to create task: %w", err)`

## Logging

- Use structured logging with key-value pairs
- Use appropriate log levels (Debug, Info, Warn, Error)
- Include relevant context in log entries
- Use `WithContext` and `With` methods to add context

## Configuration

- Load configuration from environment variables
- Group related configuration in structs
- Pass configuration to components that need it
- Validate configuration at startup

## Sample Code

```go
// Service implements the Proof of Work service
type Service struct {
    repo                  Repository
    config                *env.PowConfig
    utils                 PoWUtilsInterface
    logger                logger.Logger
}

// NewService creates a new Proof of Work service
func NewService(repo Repository, config *env.PowConfig, log logger.Logger) *Service {
    return &Service{
        repo:                  repo,
        config:                config,
        utils:                 NewPoWUtils(config, log),
        logger:                log,
    }
}

// CreateChallenge creates a new PoW challenge for a client
func (s *Service) CreateChallenge(ctx context.Context, clientIP string) (*protocol.Challenge, error) {
    s.logger.Debug("Creating challenge for client", "ip", clientIP)
    
    // Implementation...
    if err := s.repo.CreateTask(ctx, task, s.config.ChallengeTTL); err != nil {
        s.logger.Error("Failed to create task", "id", challengeID, "error", err)
        return nil, fmt.Errorf("failed to create challenge task: %w", err)
    }
    
    return challenge, nil
}