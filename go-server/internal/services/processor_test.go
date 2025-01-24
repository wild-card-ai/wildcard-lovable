package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wildcard-lovable/go-server/pkg/wildcard"
)

func TestHandleWildcardResponse(t *testing.T) {
	tests := []struct {
		name     string
		response *wildcard.Response
		want     *wildcard.APIResponse
		wantErr  bool
	}{
		{
			name: "EXEC event",
			response: &wildcard.Response{
				Event: wildcard.EventExec,
				Data: map[string]interface{}{
					"name": "listCustomers",
					"arguments": map[string]interface{}{
						"limit": 10,
					},
				},
			},
			want: &wildcard.APIResponse{
				Success: true,
				Data: map[string]interface{}{
					"name": "listCustomers",
					"arguments": map[string]interface{}{
						"limit": 10,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "STOP event",
			response: &wildcard.Response{
				Event: wildcard.EventStop,
				Data: map[string]interface{}{
					"message": "Done",
				},
			},
			want: &wildcard.APIResponse{
				Success: true,
				Data: map[string]interface{}{
					"message": "Done",
				},
			},
			wantErr: false,
		},
		{
			name: "ERROR event",
			response: &wildcard.Response{
				Event: wildcard.EventError,
				Data: map[string]interface{}{
					"message": "Something went wrong",
				},
			},
			want: &wildcard.APIResponse{
				Success: false,
				Error:   "Wildcard error: map[message:Something went wrong]",
			},
			wantErr: false,
		},
		{
			name: "Unknown event",
			response: &wildcard.Response{
				Event: "UNKNOWN",
				Data:  map[string]interface{}{},
			},
			want:    nil,
			wantErr: true,
		},
	}

	client := wildcard.NewClient("http://localhost:8080")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.HandleResponse(tt.response)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}
