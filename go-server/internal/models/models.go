package models

// MessageRequest represents the incoming user message
type MessageRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// StreamUpdate represents a single update in the SSE stream
type StreamUpdate struct {
	Type string                 `json:"type"` // "start", "progress", "complete", "error"
	Data map[string]interface{} `json:"data"`
}

// Event types for stream updates
const (
	EventStart    = "start"    // Initial event when processing starts
	EventProgress = "progress" // Progress updates during processing
	EventComplete = "complete" // Final success event
	EventError    = "error"    // Error event
)
