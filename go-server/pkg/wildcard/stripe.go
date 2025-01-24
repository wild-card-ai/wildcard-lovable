package wildcard

import (
	"fmt"

	"github.com/wildcard-lovable/go-server/pkg/wildcard/integrations/stripe"
)

// StripeClient extends the base Wildcard client with Stripe-specific functionality
type StripeClient struct {
	*Client
	stripeExecutor *stripe.Executor
}

// NewStripeClient creates a new Wildcard client with Stripe integration
func NewStripeClient(baseURL string, stripeExecutor *stripe.Executor) *StripeClient {
	return &StripeClient{
		Client:         NewClient(baseURL),
		stripeExecutor: stripeExecutor,
	}
}

// ProcessStripeMessage handles the complete flow of processing a Stripe-related message
func (c *StripeClient) ProcessStripeMessage(userID, message string) (*APIResponse, error) {
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

		// For EXEC events, execute the Stripe function and continue the conversation
		if resp.Event == "EXEC" {
			result, err := c.handleStripeExec(resp.Data)
			if err != nil {
				return nil, err
			}
			currentMessage = fmt.Sprintf("%v", result.Data)
			continue
		}

		// For all other events, handle the response and return
		return c.HandleResponse(resp)
	}
}

// handleStripeExec processes the EXEC event and executes the Stripe function
func (c *StripeClient) handleStripeExec(data map[string]interface{}) (*APIResponse, error) {
	function, err := c.HandleExecEvent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to handle exec event: %w", err)
	}

	result, err := c.stripeExecutor.ExecuteFunction(function.Name, function.Arguments)
	if err != nil {
		return &APIResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &APIResponse{
		Success: true,
		Data:    result,
	}, nil
}
