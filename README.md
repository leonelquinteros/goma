# goma 🦍

[![Go Report Card](https://goreportcard.com/badge/github.com/leonelquinteros/goma)](https://goreportcard.com/report/github.com/leonelquinteros/goma)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Goma** is a lightweight, concurrent HTTP benchmarking and load-testing tool written in Go. It's designed for quick performance checks with a minimal footprint.

## Features

- 🚀 **Concurrent execution:** Run benchmarks using a configurable number of worker goroutines.
- 🔐 **Authentication support:** Easily test endpoints requiring Bearer tokens or Basic Auth.
- 🛠 **Customizable requests:** Support for custom HTTP methods, headers, host overrides, and body data.
- 🛡 **Security:** Skip TLS verification for testing local or development environments.
- 📊 **Verbosity levels:** From silent execution to full request/response logging.

## Installation

### From Source
Ensure you have [Go](https://golang.org/dl/) installed (version 1.14 or later).

```bash
go install github.com/leonelquinteros/goma@latest
```

Alternatively, clone and build locally:

```bash
git clone https://github.com/leonelquinteros/goma.git
cd goma
go build -o goma
```

## Usage

```bash
$ goma -h
Usage of goma:
  -bearer string
    	Authorization bearer token
  -c int
    	Concurrent workers (default 1)
  -data string
    	Raw body data as string
  -headers string
    	List of headers to send in the in the following format: Header1:Value1;Header2:Value2;HeaderN:ValueN
  -host string
    	Value for the Host header to be sent in the request
  -insecure
    	Allow invalid SSL/TLS certificates
  -method string
    	HTTP request method (default "GET")
  -n int
    	Amount of iterations (default 1)
  -pass string
    	Basic Auth password
  -url string
    	Endpoint URL to request (default "https://example.com")
  -user string
    	Basic Auth username
  -v int
    	Verbosity level [0,1,2,3] (default 1)
```

### Examples

#### Basic benchmark
Run 100 requests with 10 concurrent workers:
```bash
goma -url https://api.example.com/v1/health -c 10 -n 100
```

#### POST request with data and custom headers
```bash
goma -url https://api.example.com/v1/data \
     -method POST \
     -data '{"key": "value"}' \
     -headers "Content-Type:application/json;X-Custom-Header:MyValue" \
     -c 5 -n 20
```

#### Testing with Bearer Authentication
```bash
goma -url https://api.example.com/v1/secure \
     -bearer "your-auth-token-here" \
     -n 10
```

## Development

### Running Tests
We maintain high quality with unit tests for our core logic. You can run them with:

```bash
go test -v ./...
```

### CI/CD
The project uses GitHub Actions for continuous integration, automatically building and testing changes across multiple operating systems.

## License
Distributed under the MIT License. See `LICENSE` for more information.
