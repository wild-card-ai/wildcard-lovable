package wildcard

// Response represents a response from the Wildcard API
type Response struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

// APIResponse represents a standardized API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SessionResponse represents the response from creating a new session
type SessionResponse struct {
	SessionID string `json:"sessionId"`
}

// Function represents a function to be executed by an integration
type Function struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}
