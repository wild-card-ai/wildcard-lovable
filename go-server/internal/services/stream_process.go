package services

import (
	"context"
	"encoding/json"
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

	// Add a slice to collect all messages and results
	var actionResults []string
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
		case wildcard.EventExec:
			// Step 4: Execute the function since we have an available action
			result, _ := p.wildcardClient.HandleExecEvent(resp.Data, resp.API)

			if !result.Success {
				handleError(updates, "Failed to execute function", fmt.Errorf("function execution failed"))
				currentMessage = fmt.Sprintf("Failed to execute function '%s'. Received Response: %v", resp.Data["name"], result.Error)
				continue
			}

			send(updates, EventProgress, map[string]interface{}{
				"message": "Function executed successfully",
				"result":  result.Data,
			})

			// Marshal the data into JSON and add prefix with function name and response
			dataBytes, err := json.Marshal(result.Data)
			if err != nil {
				// Fallback: just use fmt.Sprintf
				currentMessage = fmt.Sprintf("Successfully executed function '%s'. Received Response: %v", resp.Data["name"], result.Data)
			} else {
				currentMessage = fmt.Sprintf("Successfully executed function '%s'. Received Response: %s", resp.Data["name"], string(dataBytes))
			}

			// Store the result message
			actionResults = append(actionResults, currentMessage)

		case wildcard.EventStop:
			wildcardResp, err := p.wildcardClient.HandleResponse(resp)
			if handleError(updates, "Failed to handle Wildcard response", err) {
				return
			}
			data, ok := wildcardResp.Data.(map[string]interface{})
			if !ok {
				handleError(updates, "Invalid response data format", fmt.Errorf("expected map[string]interface{}, got %T", wildcardResp.Data))
				return
			}

			// Send progress update that we're generating a summary
			send(updates, EventProgress, map[string]interface{}{
				"message": "Generating summary of actions taken...",
			})

			// Collect all relevant information for OpenAI
			summaryContext := fmt.Sprintf("User request: %s\n", message)
			for i, result := range actionResults {
				summaryContext += fmt.Sprintf("Action %d: %s\n", i+1, result)
			}
			summaryContext += fmt.Sprintf("Final results: %v", data)

			// Get OpenAI to generate a user-friendly summary
			summary, err := p.openaiService.GenerateSummary(context.Background(), summaryContext)
			if err != nil {
				handleError(updates, "Failed to generate summary", err)
				return
			}

			// Send the final response with the OpenAI-generated summary
			send(updates, EventComplete, map[string]interface{}{
				"message": summary,
				"data":    data, // Include original data as well
			})
			return

		case wildcard.EventError:
			wildcardResp, err := p.wildcardClient.HandleResponse(resp)
			if err != nil {
				handleError(updates, "Failed to handle Wildcard error", err)
				return
			}
			handleError(updates, wildcardResp.Error, nil)
			return

		default:
			handleError(updates, "Unknown event", fmt.Errorf("unknown event type: %s", resp.Event))
			return
		}
	}
}
