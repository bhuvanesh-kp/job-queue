# Phase 1 — Single Process (No Distribution)

Forget Redis, Kafka, Kubernetes.

## Implement everything inside one Go program

    Client
    │
    HTTP API
    │
    In-memory Queue
    │
    Worker Goroutine
    │
    Job Executor

### Features

    POST /jobs
    Queue jobs
    Worker consumes jobs
    Execute
    Store status

## Example

### POST /jobs

    {
        "type":"sleep",
        "duration":10
    }

    Response

    {
    "job_id":"abc123"
    }

### Worker

    Queue

    ↓

    Receive Job

    ↓

    Execute

    ↓

    Update Status

    ↓

    Done

At this point you'll already understand the lifecycle.

## Project Structure

    project/
    ├── cmd/
    │   ├── api/              # HTTP API entry point
    │   │   └── main.go       # Starts HTTP server + worker pool
    │   └── worker/           # Optional: Separate worker binary
    │       └── main.go       # Starts worker pool only (for scaling)
    ├── internal/
    │   ├── api/              # HTTP layer
    │   │   ├── handler.go    # Receives client requests, enqueues jobs
    │   │   └── server.go     # HTTP server setup
    │   ├── queue/            # In-memory queue implementation
    │   │   ├── queue.go      # Channel-based job queue
    │   │   └── job.go        # Job struct definition
    │   ├── worker/           # Worker goroutine logic
    │   │   ├── pool.go       # Worker pool manager (start/stop/shutdown)
    │   │   └── worker.go     # Individual worker goroutine logic
    │   ├── executor/         # Job execution logic
    │   │   └── executor.go   # Actual business logic for processing jobs
    │   └── config/           # Configuration (queue size, worker count)
    ├── pkg/                  # Optional: Shared utilities (logger, metrics)
    ├── go.mod
    └── Makefile   
