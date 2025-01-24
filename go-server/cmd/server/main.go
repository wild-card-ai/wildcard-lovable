package main

import (
	"log"
	"net/http"

	"github.com/wildcard-lovable/go-server/internal/config"
	"github.com/wildcard-lovable/go-server/internal/handlers"
	"github.com/wildcard-lovable/go-server/internal/middleware"
	"github.com/wildcard-lovable/go-server/internal/services"
	"github.com/wildcard-lovable/go-server/pkg/wildcard/integrations/stripe"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Initialize services
	stripeStore := services.NewStripeKeyStore()
	stripeExecutor := stripe.NewExecutor(stripeStore)
	openaiService := services.NewOpenAIService(cfg.OpenAIAPIKey)
	processor := services.NewProcessor(cfg.WildcardBackendURL, stripeExecutor, openaiService)

	// Initialize handler
	messageHandler := handlers.NewMessageHandler(processor, stripeStore)

	// Set up routes with CORS middleware
	http.HandleFunc("/process", middleware.CorsMiddleware(messageHandler.ProcessMessage))
	http.HandleFunc("/process-stream", middleware.CorsMiddleware(messageHandler.StreamProcess))
	http.HandleFunc("/register-stripe", middleware.CorsMiddleware(messageHandler.HandleStripeRegistration))

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe("0.0.0.0:"+cfg.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
