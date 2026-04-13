# Goma - Go HTTP Benchmarking Tool

Goma is a lightweight, concurrent HTTP benchmarking and load-testing tool written in Go. It allows users to send multiple HTTP requests to a specified endpoint with configurable concurrency and verbosity.

## Project Overview

- **Purpose:** Simple CLI tool for HTTP request performance testing.
- **Language:** Go (1.23+).
- **Core Dependencies:** Standard library (`net/http`, `sync`, `flag`, etc.).

## Architecture

The project follows a modular worker-pool pattern:
- **Runner:** Manages the worker pool and `http.Client` lifecycle.
- **Config:** Encapsulates benchmark parameters and request construction logic.
- **Workers:** Concurrent goroutines that execute requests and log results.
- **Synchronization:** Uses `sync.WaitGroup` and channels for task distribution.

## Core Features

- **Concurrency:** Configurable number of workers (`-c`).
- **Iterations:** Configurable total number of requests (`-n`).
- **Authentication:** Supports Bearer tokens and Basic Auth.
- **Customization:** Supports custom HTTP methods, headers, host overrides, and body data.
- **Security:** Optional TLS verification skipping (`-insecure`).
- **Verbosity:** Multiple logging levels (0-3).

## Development Workflow

### Building
```bash
go build -o goma main.go
```

### Testing
```bash
go test -v ./...
```

### Running
```bash
./goma -url https://example.com -c 10 -n 100
```

### CI/CD
GitHub Actions workflow runs builds, tests, and benchmarks on Ubuntu, Windows, and macOS. Triggered on `push` and `pull_request` to `master`.

## Key Files
- `main.go`: Application entry point and core logic.
- `main_test.go`: Unit tests for request construction and configuration.
- `go.mod`: Module definition.
- `README.md`: Public documentation and usage guide.

## Future Improvements
- **Result Aggregation:** Track success/failure rates and latency percentiles.
- **Reporting:** Provide a final summary report upon benchmark completion.
