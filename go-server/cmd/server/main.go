package main

import (
	"log"
	"net/http"
	"github.com/wildcard-lovable/go-server/internal/config"
	"github.com/wildcard-lovable/go-server/internal/handlers"
	"github.com/wildcard-lovable/go-server/internal/middleware"
	"github.com/wildcard-lovable/go-server/internal/services"
	"github.com/wildcard-lovable/go-server/pkg/stripe"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Initialize Stripe executor
	stripeExecutor := stripe.NewExecutor(cfg.StripeAPIKey)

	// Initialize OpenAI service
	openaiService := services.NewOpenAIService(cfg.OpenAIAPIKey)

	// Initialize processor
	processor := services.NewProcessor(cfg.WildcardBackendURL, stripeExecutor, openaiService)

	// Initialize handler
	messageHandler := handlers.NewMessageHandler(processor)

	// Set up routes with CORS middleware
	http.HandleFunc("/process", middleware.CorsMiddleware(messageHandler.ProcessMessage))
	http.HandleFunc("/process-stream", middleware.CorsMiddleware(messageHandler.StreamProcess))

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
} 