package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/wildcard-lovable/go-server/internal/models"
	"github.com/wildcard-lovable/go-server/pkg/stripe"
)

// Processor handles the processing of user messages
type Processor struct {
	wildcardBaseURL string
	stripeExecutor  *stripe.Executor
	openaiService   *OpenAIService
}

// NewProcessor creates a new processor instance
func NewProcessor(wildcardBaseURL string, stripeExecutor *stripe.Executor, openaiService *OpenAIService) *Processor {
	return &Processor{
		wildcardBaseURL: wildcardBaseURL,
		stripeExecutor:  stripeExecutor,
		openaiService:   openaiService,
	}
}

// ProcessMessage handles the complete flow of processing a user message
func (p *Processor) ProcessMessage(userID, message string) (*models.APIResponse, error) {
	ctx := context.Background()

	// First, interpret the message using OpenAI to determine if it's Stripe-related
	isStripeRelated, llmResponse, err := p.openaiService.InterpretMessage(ctx, message)
	// Print the result for debugging purposes
	fmt.Printf("Message interpretation result - isStripeRelated: %v, response: %v\n", isStripeRelated, llmResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to interpret message: %w", err)
	}

	if !isStripeRelated {
		return &models.APIResponse{
			Success: true,
			Data:    llmResponse,
		}, nil
	}

	// If it is Stripe-related, use Wildcard to process it
	sessionID, err := p.createSession(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Start the conversation loop with Wildcard
	// Continue the conversation loop with Wildcard until we get a STOP or ERROR event
	// We may want to include a "context" parameter as an input
	currentMessage := message
	for {
		resp, err := p.processWithWildcard(userID, sessionID, currentMessage)
		if err != nil {
			return nil, fmt.Errorf("failed to process with wildcard: %w", err)
		}

		result, err := p.handleWildcardResponse(resp)
		if err != nil {
			return nil, err
		}

		// If we got a STOP or ERROR event, return the result
		if resp.Event == "STOP" || resp.Event == "ERROR" {
			return result, nil
		}

		// For EXEC events, send the result back to Wildcard
		currentMessage = fmt.Sprintf("%v", result.Data)
	}
}

// createSession creates a new session for the user
func (p *Processor) createSession(userID string) (string, error) {
	url := fmt.Sprintf("%s/session/%s", p.wildcardBaseURL, userID)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var sessionResp models.SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionResp); err != nil {
		return "", err
	}

	return sessionResp.SessionID, nil
}

// processWithWildcard sends the message to the Wildcard backend for processing
func (p *Processor) processWithWildcard(userID, sessionID, message string) (*models.WildcardResponse, error) {
	url := fmt.Sprintf("%s/process/%s/%s", p.wildcardBaseURL, userID, sessionID)
	
	payload := map[string]string{"message": message}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wildcardResp models.WildcardResponse
	if err := json.NewDecoder(resp.Body).Decode(&wildcardResp); err != nil {
		return nil, err
	}

	return &wildcardResp, nil
}

// handleWildcardResponse processes the response from the Wildcard backend
func (p *Processor) handleWildcardResponse(resp *models.WildcardResponse) (*models.APIResponse, error) {
	switch resp.Event {
	case "EXEC":
		return p.handleExecEvent(resp.Data)
	case "STOP":
		return &models.APIResponse{
			Success: true,
			Data:    resp.Data,
		}, nil
	case "ERROR":
		return &models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Wildcard error: %v", resp.Data),
		}, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", resp.Event)
	}
}

// handleExecEvent processes the EXEC event and executes the appropriate Stripe function
func (p *Processor) handleExecEvent(data map[string]interface{}) (*models.APIResponse, error) {
	var function models.WildcardFunction
	functionData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(functionData, &function); err != nil {
		return nil, err
	}

	result, err := p.stripeExecutor.ExecuteFunction(function.Name, function.Arguments)
	if err != nil {
		return &models.APIResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &models.APIResponse{
		Success: true,
		Data:    result,
	}, nil
}