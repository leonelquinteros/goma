# Goma - Go HTTP Benchmarking Tool

Goma is a lightweight, concurrent HTTP benchmarking and load-testing tool written in Go. It allows users to send multiple HTTP requests to a specified endpoint with configurable concurrency and verbosity.

## Project Overview

- **Purpose:** Simple CLI tool for HTTP request performance testing.
- **Language:** Go (1.23+ in `go.mod`).
- **Core Dependencies:** Standard library (`net/http`, `sync`, `flag`, etc.).

## Architecture

The project follows a simple worker-pool pattern:
- **Main Thread:** Parses flags, initializes the shared channel, and spawns worker goroutines.
- **Workers:** Listen on a shared channel for request IDs and perform HTTP requests.
- **Synchronization:** Uses `sync.WaitGroup` to wait for all workers to complete and a channel to distribute tasks.

## Core Features

- **Concurrency:** Configurable number of workers (`-c`).
- **Iterations:** Configurable total number of requests (`-n`).
- **Authentication:** Supports Bearer tokens and Basic Auth.
- **Customization:** Supports custom HTTP methods, headers, host overrides, and body data.
- **Security:** Option to skip TLS verification (`-insecure`).
- **Verbosity:** Multiple logging levels (0-3).

## Development Workflow

### Building
```bash
go build -o goma main.go
```

### Running
```bash
./goma -url https://example.com -c 10 -n 100
```

### CI/CD
The project uses GitHub Actions for continuous integration, running builds, tests, and benchmarks on Ubuntu and Windows.
Location: `.github/workflows/ci.yml`

## Project Status & Recommendations

- **Testing:** Unit tests have been added in `main_test.go` to verify request creation, header parsing, and authentication configuration.
- **Structure:** The code has been modularized into `Config` and `Runner` structs, separating CLI parsing from execution logic.
- **Configuration:** A custom `http.Client` is now used, avoiding global modifications to `http.DefaultTransport`.
- **Error Handling:** Errors are logged, and the benchmark continues. Future work could include tracking success/failure rates and providing a final summary.
- **Dependencies:** Unused dependencies (like `delve`) have been removed from `go.mod`.

## Key Files
- `main.go`: Entry point and core logic.
- `go.mod`: Module definition and dependencies.
- `README.md`: Usage instructions and examples.
