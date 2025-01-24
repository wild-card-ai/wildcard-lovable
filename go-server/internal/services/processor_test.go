package services

import (
	"encoding/json"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/wildcard-lovable/go-server/internal/models"
	"github.com/wildcard-lovable/go-server/pkg/stripe"
	"github.com/wildcard-lovable/go-server/internal/config"
)

func TestHandleWildcardResponse(t *testing.T) {
	// Get config to access Stripe API key
	cfg := config.NewConfig()

	// Create stripe executor with actual API key
	stripeExecutor := stripe.NewExecutor(cfg.StripeAPIKey)

	// Create the processor with minimal config
	processor := &Processor{
		stripeExecutor: stripeExecutor,
	}

	// Raw JSON response we expect from Wildcard
	rawJSONResponse := `{
		"event": "EXEC",
		"data": {
			"name": "stripe_post_products",
			"arguments": {
				"name": "Test Product",
				"description": "A test product",
				"active": true
			}
		}
	}`

	var parsedResponse models.WildcardResponse
	err := json.Unmarshal([]byte(rawJSONResponse), &parsedResponse)
	assert.NoError(t, err, "Failed to parse JSON response")

	tests := []struct {
		name     string
		response *models.WildcardResponse
		wantErr  bool
	}{
		{
			name:     "Create Product from JSON response",
			response: &parsedResponse,
			wantErr:  false,
		},
		{
			name: "Create Product direct struct",
			response: &models.WildcardResponse{
				Event: "EXEC",
				Data: map[string]interface{}{
					"name": "stripe_post_products",
					"arguments": map[string]interface{}{
						"name": "Test Product",
						"description": "A test product",
						"active": true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "STOP event",
			response: &models.WildcardResponse{
				Event: "STOP",
				Data: map[string]interface{}{
					"message": "Operation completed",
				},
			},
			wantErr: false,
		},
		{
			name: "ERROR event",
			response: &models.WildcardResponse{
				Event: "ERROR",
				Data: map[string]interface{}{
					"error": "Something went wrong",
				},
			},
			wantErr: false,
		},
		{
			name: "Unknown event",
			response: &models.WildcardResponse{
				Event: "UNKNOWN",
				Data: map[string]interface{}{
					"message": "Unknown event",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.handleWildcardResponse(tt.response)
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			assert.NoError(t, err)
			assert.NotNil(t, result)

			switch tt.response.Event {
			case "EXEC":
				assert.True(t, result.Success)
				// For EXEC events, verify the Data field contains the Stripe response
				assert.NotNil(t, result.Data)
			case "STOP":
				assert.True(t, result.Success)
				assert.Equal(t, "Operation completed", result.Data.(map[string]interface{})["message"])
			case "ERROR":
				assert.False(t, result.Success)
				assert.Contains(t, result.Error, "Something went wrong")
			}
		})
	}
} 