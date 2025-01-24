package models

// MessageRequest represents the incoming user message
type MessageRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// SessionResponse represents the response from creating a new session
type SessionResponse struct {
	SessionID string `json:"session_id"`
}

// WildcardResponse represents the response from the Wildcard backend
type WildcardResponse struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

// OpenAIFunction represents a function call from OpenAI
type OpenAIFunction struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string     `json:"error,omitempty"`
} 