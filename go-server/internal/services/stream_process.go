package services

import (
	"context"
	"fmt"

	"github.com/wildcard-lovable/go-server/internal/models"
	"github.com/wildcard-lovable/go-server/pkg/wildcard"
)

// Event types for stream updates
const (
	EventStart    = "start"    // Initial event when processing starts
	EventProgress = "progress" // Progress updates during processing
	EventComplete = "complete" // Final success event
	EventError    = "error"    // Error event
)

// Helper functions
func send(updates chan<- models.StreamUpdate, eventType string, data map[string]interface{}) {
	updates <- models.StreamUpdate{
		Type: eventType,
		Data: data,
	}
}

func handleError(updates chan<- models.StreamUpdate, msg string, err error) bool {
	if err != nil {
		send(updates, EventError, map[string]interface{}{
			"message": msg,
			"error":   err.Error(),
		})
		return true
	}
	return false
}

// StreamProcessMessage - Processes a user message, executes integrations actions if needed
func (p *Processor) StreamProcessMessage(userID, message string, updates chan<- models.StreamUpdate) {
	defer close(updates)

	// Start processing
	send(updates, EventStart, map[string]interface{}{
		"message": "Starting message processing",
	})

	// Step 1: Process with OpenAI to determine if the given action is related to an integration
	send(updates, EventProgress, map[string]interface{}{
		"message": "Analyzing message with OpenAI",
	})

	isStripeRelated, llmResponse, err := p.openaiService.InterpretMessage(context.Background(), message)
	if err != nil {
		handleError(updates, "Failed to process with OpenAI", err)
		return
	}

	if !isStripeRelated {
		send(updates, EventComplete, map[string]interface{}{
			"message": llmResponse,
		})
		return
	}

	// Step 2: Create Wildcard session since we know the action is related to Stripe
	send(updates, EventProgress, map[string]interface{}{
		"message": "Creating Wildcard session",
	})

	sessionID, err := p.wildcardClient.CreateSession(userID)
	if handleError(updates, "Failed to create session", err) {
		return
	}

	// Step 3: Process with Wildcard to get the Stripe action to execute
	currentMessage := message
	for {
		send(updates, EventProgress, map[string]interface{}{
			"message": "Processing with Wildcard",
		})

		resp, err := p.wildcardClient.ProcessMessage(userID, sessionID, currentMessage)
		if handleError(updates, "Failed to process with Wildcard", err) {
			return
		}

		switch resp.Event {
		case "EXEC":
			// Step 4: Execute the Stripe function since we have an available action
			result, err := p.wildcardClient.(*wildcard.StripeClient).handleStripeExec(resp.Data)
			if handleError(updates, "Failed to execute Stripe function", err) {
				return
			}

			send(updates, EventProgress, map[string]interface{}{
				"message": "Stripe function executed successfully",
				"result":  result,
			})
			currentMessage = fmt.Sprintf("%v", result.Data)

		case "STOP":
			wildcardResp, err := p.wildcardClient.HandleResponse(resp)
			if handleError(updates, "Failed to handle Wildcard response", err) {
				return
			}
			send(updates, EventComplete, wildcardResp.Data)
			return

		case "ERROR":
			wildcardResp, err := p.wildcardClient.HandleResponse(resp)
			if err != nil {
				handleError(updates, "Failed to handle Wildcard error", err)
				return
			}
			handleError(updates, "Wildcard error", fmt.Errorf("%v", wildcardResp.Error))
			return

		default:
			handleError(updates, "Unknown event", fmt.Errorf("unknown event type: %s", resp.Event))
			return
		}
	}
}
