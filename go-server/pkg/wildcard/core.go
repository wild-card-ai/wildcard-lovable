package wildcard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client handles core Wildcard operations
type Client struct {
	baseURL string
}

// NewClient creates a new Wildcard client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
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

	return sessionResp.SessionID, nil
}

// ProcessMessage sends a message to Wildcard for processing
func (c *Client) ProcessMessage(userID, sessionID, message string) (*Response, error) {
	url := fmt.Sprintf("%s/process/%s/%s", c.baseURL, userID, sessionID)

	payload := map[string]string{"message": message}
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
	switch resp.Event {
	case "STOP":
		return &APIResponse{
			Success: true,
			Data:    resp.Data,
		}, nil
	case "ERROR":
		return &APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Wildcard error: %v", resp.Data),
		}, nil
	case "EXEC":
		return &APIResponse{
			Success: true,
			Data:    resp.Data,
		}, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", resp.Event)
	}
}

// HandleExecEvent processes the EXEC event data into a Function
func (c *Client) HandleExecEvent(data map[string]interface{}) (*Function, error) {
	var function Function
	functionData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal function data: %w", err)
	}

	if err := json.Unmarshal(functionData, &function); err != nil {
		return nil, fmt.Errorf("failed to unmarshal function data: %w", err)
	}

	return &function, nil
}
