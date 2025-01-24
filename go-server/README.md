# Wildcard Go Server

This is a Go implementation of a server that processes user messages through the Wildcard backend and executes Stripe API operations based on the responses.

## Prerequisites

- Go 1.19 or later
- Stripe API key
- OpenAI API key
- Access to Wildcard backend

## Project Structure

```
go-server/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── api/                  # API related code
│   ├── config/              # Configuration management
│   ├── handlers/            # HTTP request handlers
│   ├── models/              # Data models
│   └── services/            # Business logic
├── pkg/
│   └── stripe/              # Stripe integration
└── README.md
```

## Configuration

The server requires the following environment variables:

```bash
export PORT=8080                                  # Server port (optional, defaults to 8080)
export WILDCARD_BACKEND_URL=http://localhost:8000 # Wildcard backend URL (if hosted)
export OPENAI_API_KEY=your_openai_api_key        # OpenAI API key
export STRIPE_API_KEY=your_stripe_api_key        # Stripe API key
```

## Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```

## Running the Server

```bash
go run cmd/server/main.go
```

## API Endpoints

### Process Message (Regular)
```
POST /process
```
Processes a message and returns a single response.

Request body:
```json
{
    "user_id": "string",
    "message": "string"
}
```

Response:
```json
{
    "success": true,
    "data": {},
    "error": "string"
}
```

### Process Message (Streaming)
```
POST /process-stream
```
Processes a message and streams updates using Server-Sent Events (SSE).

Request body:
```json
{
    "user_id": "string",
    "message": "string"
}
```

Example curl request:
```bash
curl -N -X POST http://localhost:8082/process-stream \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "message": "Create a new product called Premium Plan for $10 per month"
  }'
```
Note: The `-N` flag is required for curl to disable buffering and show the stream events in real-time.

Stream Events Format:
```json
{
    "type": "start|progress|complete|error",
    "data": {
        "message": "string",
        "result": {},
        "error": "string"
    }
}
```

Event Types:
- `start`: Initial event when processing starts
- `progress`: Progress updates during processing
- `complete`: Final success event
- `error`: Error event

## Development

To add new Stripe functions:

1. Add the function to the `FunctionMap` in `pkg/stripe/executor.go`
2. Implement the corresponding method in the `Executor` struct
