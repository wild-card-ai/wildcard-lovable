package services

import (
	"context"
	"fmt"

	"github.com/wildcard-lovable/go-server/pkg/wildcard"
	"github.com/wildcard-lovable/go-server/pkg/wildcard/integrations/stripe"
)

// Processor handles the processing of user messages
type Processor struct {
	wildcardClient *wildcard.Client
	openaiService  *OpenAIService
}

// NewProcessor creates a new processor instance
func NewProcessor(wildcardBaseURL string, stripeExecutor *stripe.Executor, openaiService *OpenAIService) *Processor {
	client := wildcard.NewClient(wildcardBaseURL)
	client.RegisterExecutor(wildcard.APINameStripe, stripeExecutor)

	return &Processor{
		wildcardClient: client,
		openaiService:  openaiService,
	}
}

// ProcessMessage handles the complete flow of processing a user message
func (p *Processor) ProcessMessage(userID, message string) (*wildcard.APIResponse, error) {
	ctx := context.Background()

	// First, interpret the message using OpenAI to determine if it's Stripe-related
	isStripeRelated, llmResponse, err := p.openaiService.InterpretMessage(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to interpret message: %w", err)
	}

	if !isStripeRelated {
		return &wildcard.APIResponse{
			Success: true,
			Data:    llmResponse,
		}, nil
	}

	// If it is Stripe-related, use Wildcard to process it
	return p.wildcardClient.ProcessAPIMessage(userID, message)
}
