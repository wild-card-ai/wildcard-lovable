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
export WILDCARD_BACKEND_URL=http://localhost:8000 # Wildcard backend URL (optional)
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

### POST /process

Process a user message and execute any necessary Stripe operations.

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
    "success": boolean,
    "data": object,
    "error": string
}
```

## Development

To add new Stripe functions:

1. Add the function to the `FunctionMap` in `pkg/stripe/executor.go`
2. Implement the corresponding method in the `Executor` struct
