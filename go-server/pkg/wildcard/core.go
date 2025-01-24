package wildcard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Executor is the interface that all integration executors must implement
type Executor interface {
	ExecuteFunction(name string, arguments map[string]interface{}) (interface{}, error)
}

// Client handles core Wildcard operations
type Client struct {
	baseURL   string
	executors map[string]Executor
}

// NewClient creates a new Wildcard client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:   baseURL,
		executors: make(map[string]Executor),
	}
}

// RegisterExecutor registers an executor for a specific API
func (c *Client) RegisterExecutor(apiName string, executor Executor) {
	c.executors[apiName] = executor
}

// CreateSession creates a new session for the user
func (c *Client) CreateSession(userID string) (string, error) {
	url := fmt.Sprintf("%s/session/%s", c.baseURL, userID)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer resp.Body.Close()

	var sessionResp SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionResp); err != nil {
		return "", fmt.Errorf("failed to decode session response: %w", err)
	}

	fmt.Println("Session ID:", sessionResp.SessionID)
	return sessionResp.SessionID, nil
}

// ProcessMessage sends a message to Wildcard for processing
func (c *Client) ProcessMessage(userID, sessionID, message string) (*Response, error) {
	url := fmt.Sprintf("%s/process/%s/%s", c.baseURL, userID, sessionID)

	// Create the JSON payload
	payload := map[string]interface{}{
		"message": message,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to process message: %w", err)
	}
	defer resp.Body.Close()

	var wildcardResp Response
	if err := json.NewDecoder(resp.Body).Decode(&wildcardResp); err != nil {
		return nil, fmt.Errorf("failed to decode wildcard response: %w", err)
	}

	return &wildcardResp, nil
}

// HandleResponse processes the response from Wildcard
func (c *Client) HandleResponse(resp *Response) (*APIResponse, error) {
	fmt.Println("Handling response:", resp)
	switch resp.Event {
	case EventStop:
		return &APIResponse{
			Success: true,
			Data:    resp.Data,
		}, nil
	case EventError:
		return &APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Wildcard error: %v", resp.Data),
		}, nil
	case EventExec:
		return &APIResponse{
			Success: true,
			Data:    resp.Data,
		}, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", resp.Event)
	}
}

// HandleExecEvent processes the EXEC event data into a Function and executes it
func (c *Client) HandleExecEvent(data map[string]interface{}, apiName string) (*APIResponse, error) {
	// Debug logging
	fmt.Printf("HandleExecEvent received data: %+v\n", data)
	fmt.Printf("HandleExecEvent received apiName: %s\n", apiName)

	// Safely extract and validate required fields
	name, ok := data["name"].(string)
	if !ok {
		fmt.Printf("Error: name field is missing or invalid. Data type: %T\n", data["name"])
		return &APIResponse{
			Success: false,
			Error:   fmt.Sprintf("We tried to execute a function, but the function name was missing or invalid"),
		}, nil
	}

	arguments, ok := data["arguments"].(map[string]interface{})
	if !ok {
		fmt.Printf("Error: arguments field is missing or invalid. Data type: %T\n", data["arguments"])
		return &APIResponse{
			Success: false,
			Error:   fmt.Sprintf("We tried to execute function '%s', but the arguments were missing or invalid", name),
		}, nil
	}

	function := Function{
		Name:      name,
		API:       apiName,
		Arguments: arguments,
	}

	// Verify API name matches
	if function.API != apiName {
		return &APIResponse{
			Success: false,
			Error:   fmt.Sprintf("We tried to execute function '%s', but there was an API mismatch (expected %s, got %s)", name, apiName, function.API),
		}, nil
	}

	// Get the executor for this API
	executor, ok := c.executors[apiName]
	if !ok {
		return &APIResponse{
			Success: false,
			Error:   fmt.Sprintf("We tried to execute function '%s', but no executor was found for API '%s'", name, apiName),
		}, nil
	}

	// Execute the function
	result, err := executor.ExecuteFunction(function.Name, function.Arguments)
	if err != nil {
		return &APIResponse{
			Success: false,
			Error:   fmt.Sprintf("We tried to execute function '%s', but received error: %v", name, err),
		}, nil
	}

	return &APIResponse{
		Success: true,
		Data:    result,
	}, nil
}

// ProcessAPIMessage handles the complete flow of processing an API-specific message
func (c *Client) ProcessAPIMessage(userID, message string) (*APIResponse, error) {
	// Create a session
	sessionID, err := c.CreateSession(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Process messages with Wildcard until we get a final response
	currentMessage := message
	for {
		resp, err := c.ProcessMessage(userID, sessionID, currentMessage)
		if err != nil {
			return nil, fmt.Errorf("failed to process message: %w", err)
		}

		// For EXEC events, execute the function and continue the conversation
		if resp.Event == EventExec {
			result, _ := c.HandleExecEvent(resp.Data, resp.API)
			if !result.Success {
				// Send the error message back to continue the conversation
				currentMessage = result.Error
				continue
			}

			// Format as a descriptive string message
			currentMessage = fmt.Sprintf("The %s operation was successful. The result was: %v", resp.Data["name"], result.Data)
			continue
		}

		// For all other events, handle the response and return
		return c.HandleResponse(resp)
	}
}
