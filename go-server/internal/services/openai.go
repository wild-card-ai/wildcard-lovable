package services

import (
	"context"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIService handles interactions with OpenAI API
type OpenAIService struct {
	client *openai.Client
}

// NewOpenAIService creates a new OpenAI service
func NewOpenAIService(apiKey string) *OpenAIService {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIService{
		client: client,
	}
}

// InterpretMessage sends a message to OpenAI to determine if it requires Stripe integration
func (s *OpenAIService) InterpretMessage(ctx context.Context, message string) (bool, string, error) {
	resp, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.F(openai.ChatModelGPT4o),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a helpful assistant. If the user's message requires Stripe API integration (like payments, customers, products, etc.), respond with 'true' followed by a brief explanation. Otherwise, respond with 'false' and provide a helpful response to their query."),
			openai.UserMessage(message),
		}),
	})

	if err != nil {
		return false, "", fmt.Errorf("failed to interpret message: %w", err)
	}

	if len(resp.Choices) > 0 {
		content := resp.Choices[0].Message.Content
		isStripeRelated := len(content) >= 4 && content[:4] == "true"
		response := content[4:] // Get the explanation/response part
		return isStripeRelated, response, nil
	}

	return false, "", nil
} 
