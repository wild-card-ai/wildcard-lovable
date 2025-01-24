package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wildcard-lovable/go-server/internal/models"
	"github.com/wildcard-lovable/go-server/internal/services"
)

// MessageHandler handles HTTP requests for message processing
type MessageHandler struct {
	processor   *services.Processor
	stripeStore *services.StripeKeyStore
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(processor *services.Processor, stripeStore *services.StripeKeyStore) *MessageHandler {
	return &MessageHandler{
		processor:   processor,
		stripeStore: stripeStore,
	}
}

// ProcessMessage handles the regular HTTP POST request
func (h *MessageHandler) ProcessMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.processor.ProcessMessage(req.UserID, req.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// StreamProcess handles SSE streaming of message processing
func (h *MessageHandler) StreamProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel for updates
	updates := make(chan models.StreamUpdate)

	// Parse request
	var req models.MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendSSEError(w, "Failed to decode request", err)
		return
	}

	// Start processing in a goroutine
	go h.processor.StreamProcessMessage(req.UserID, req.Message, updates)

	// Stream updates to client
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Stream updates until done or client disconnects
	for update := range updates {
		data, err := json.Marshal(update)
		if err != nil {
			sendSSEError(w, "Failed to marshal update", err)
			return
		}

		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}
}

func sendSSEError(w http.ResponseWriter, msg string, err error) {
	errUpdate := models.StreamUpdate{
		Type: "error",
		Data: map[string]interface{}{
			"message": msg,
			"error":   err.Error(),
		},
	}
	data, _ := json.Marshal(errUpdate)
	fmt.Fprintf(w, "data: %s\n\n", data)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

type StripeRegistrationRequest struct {
	UserID string `json:"userId"`
	APIKey string `json:"apiKey"`
}

func (h *MessageHandler) HandleStripeRegistration(w http.ResponseWriter, r *http.Request) {
	var req StripeRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.stripeStore.RegisterKey(req.UserID, req.APIKey); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
