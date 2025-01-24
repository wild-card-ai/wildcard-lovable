package wildcard

// Response represents a response from the Wildcard API
type Response struct {
	Event string                 `json:"event"`
	API   string                 `json:"api"`
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
	SessionID string `json:"session_id"`
}

// Function represents a function to be executed by an integration
type Function struct {
	Name      string                 `json:"name"`
	API       string                 `json:"api"`
	Arguments map[string]interface{} `json:"arguments"`
}

// Event types for Wildcard responses
const (
	EventExec  = "EXEC"  // Execute a function
	EventStop  = "STOP"  // Stop processing
	EventError = "ERROR" // Error occurred
)

// API names for different integrations
const (
	APINameStripe = "stripe" // Stripe API integration
)
