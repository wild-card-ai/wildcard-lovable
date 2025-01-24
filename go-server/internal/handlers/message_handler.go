package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/wildcard-lovable/go-server/internal/models"
	"github.com/wildcard-lovable/go-server/internal/services"
)

// MessageHandler handles HTTP requests for message processing
type MessageHandler struct {
	processor *services.Processor
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(processor *services.Processor) *MessageHandler {
	return &MessageHandler{
		processor: processor,
	}
}

// ProcessMessage handles the HTTP request for processing a message
func (h *MessageHandler) ProcessMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.processor.ProcessMessage(req.UserID, req.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 