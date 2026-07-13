# Arachne

A concurrent web crawler in Go with graceful shutdown, statistics collection, and configurable depth limits.

## Features

- **Concurrent Crawling**: Worker pool architecture for efficient parallel crawling.
- **Configurable Depth**: Control crawl depth to limit recursion.
- **Graceful Shutdown**: Handles interrupt signals (SIGINT) cleanly.
- **Statistics Collection**: Tracks URLs crawled, fetch success/failure rates, links found, and total bytes transferred.
- **Thread-Safe Visited Set**: Prevents duplicate URL processing using mutex-protected set.
- **Context-Based Cancellation**: Uses Go's context package for coordinated shutdown.
- **URL Normalization**: Parses and normalizes URLs to avoid duplicates.
- **Configurable Timeout**: HTTP request timeout for network resilience.

## Tech Stack

- **Language**: Go 1.26.4
- **Libraries**: net/http, sync, context, golang.org/x/net/html

## Project Structure

```text
arachne/
├── cmd/
│   └── crawler/
│       └── main.go              # CLI entry point with flag parsing
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration struct
│   ├── crawler/
│   │   ├── crawler.go           # Main crawler orchestration logic
│   │   ├── worker.go            # Worker implementation for task processing
│   │   └── stats.go             # Statistics collection and reporting
│   ├── fetcher/
│   │   └── fetcher.go           # HTTP client with timeout
│   ├── frontier/
│   │   ├── queue.go             # Thread-safe task queue
│   │   ├── task.go              # Task struct with URL and depth
│   │   └── visited.go           # Thread-safe visited URL set
│   └── parser/
│       └── parser.go            # HTML parsing and URL extraction
└── go.mod                       # Go module descriptor
```

## Setup

1. Clone the repository:
```bash
git clone https://github.com/thedeepak12/arachne.git
cd arachne
```

2. Tidy the Go modules:
```bash
go mod tidy
```

3. Run the crawler:
```bash
go run ./cmd/crawler/main.go -url https://example.com/ -depth 2 -workers 10
```

## Usage

**Command-line flags:**
- `-url`: Seed URL to start crawling (required)
- `-depth`: Maximum crawl depth (default: 2)
- `-workers`: Number of concurrent workers (default: 10)
- `-timeout`: HTTP request timeout in seconds (default: 10)

**Examples:**

```bash
# Basic crawl with default settings
go run ./cmd/crawler/main.go -url https://example.com/

# Crawl with custom depth and workers
go run ./cmd/crawler/main.go -url https://example.com/ -depth 3 -workers 50

# Crawl with custom timeout
go run ./cmd/crawler/main.go -url https://example.com/ -timeout 15
```

**Interrupt Handling:**
Press `Ctrl+C` to gracefully shut down the crawler. Statistics will be printed upon termination.

## Architecture

1. **Frontier**: Manages the task queue and visited URL set
2. **Fetcher**: Handles HTTP requests with configurable timeout
3. **Parser**: Extracts links from HTML and normalizes URLs
4. **Worker Pool**: Concurrent workers process tasks from the queue
5. **Distributor**: Goroutine that distributes tasks to workers
6. **Statistics**: Thread-safe counters for crawl metrics

## License

Distributed under the MIT License. See [LICENSE](https://github.com/thedeepak12/arachne/blob/main/LICENSE) for more information.
