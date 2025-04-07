# Word of Wisdom Client

This is a client for the Word of Wisdom service, which provides quotes of wisdom after solving a Proof of Work (PoW) challenge.

## Architecture

The client has been refactored to follow clean architecture principles and best practices:

- **Dependency Injection**: All components receive their dependencies through constructors
- **Interface-based Design**: Components interact through interfaces, not concrete implementations
- **Separation of Concerns**: Code is organized by feature/domain, not by technical layer
- **Error Handling**: Errors are wrapped with context for better debugging
- **Structured Logging**: Consistent logging with key-value pairs

## Project Structure

```
client/
├── cmd/                    # Application entry points
│   └── client/             # Main client application
│       └── main.go         # Entry point
├── internal/               # Private application code
│   ├── config/             # Configuration
│   │   └── env.go          # Environment-based configuration
│   ├── domain/             # Domain models
│   │   └── models.go       # Core domain models
│   ├── http/               # HTTP server
│   │   ├── deps.go         # HTTP interfaces
│   │   ├── server.go       # HTTP server implementation
│   │   └── static/         # Static web UI files
│   ├── logger/             # Logging
│   │   ├── deps.go         # Logger interface
│   │   └── logger.go       # Logger implementation
│   ├── pow/                # Proof of Work
│   │   ├── deps.go         # PoW interfaces
│   │   ├── solver.go       # PoW solver implementation
│   │   └── solver_test.go  # PoW solver tests
│   └── tcp/                # TCP client
│       ├── deps.go         # TCP client interfaces
│       └── client.go       # TCP client implementation
```

## Features

- **CLI Mode**: Run as a command-line application to fetch a single quote
- **HTTP Server**: Run as a web server with a UI for interacting with the service
- **Load Testing**: Built-in load testing capabilities
- **Graceful Shutdown**: Handles termination signals properly
- **Configurable**: Settings can be configured via environment variables or command-line flags

## Usage

### Command-line Mode

```bash
# Run in CLI mode
go run cmd/client/main.go --cli --server localhost:8080
```

### HTTP Server Mode

```bash
# Run HTTP server
go run cmd/client/main.go --server localhost:8080 --http localhost:3000
```

### Environment Variables

The following environment variables can be used to configure the client:

- `SERVER_ADDRESS`: Word of Wisdom server address (default: "localhost:8080")
- `HTTP_ADDRESS`: HTTP server address (default: "localhost:3000")
- `CONNECT_TIMEOUT`: TCP connection timeout (default: 10s)
- `READ_TIMEOUT`: TCP read timeout (default: 30s)
- `WRITE_TIMEOUT`: TCP write timeout (default: 30s)
- `SOLVE_TIMEOUT`: Maximum time to spend solving a PoW challenge (default: 30s)

## Development

### Running Tests

```bash
go test ./...
```

### Generating Mocks

```bash
go generate ./...