package config

import (
	"log"
	"os"
)

type Config struct {
	Port              string
	WildcardBackendURL string
	OpenAIAPIKey      string
	StripeAPIKey      string
}

func NewConfig() *Config {
	return &Config{
		Port:              getEnvOrDefault("PORT", "8080"),
		WildcardBackendURL: getEnvOrDefault("WILDCARD_BACKEND_URL", "http://localhost:8000"),
		OpenAIAPIKey:      getEnv("OPENAI_API_KEY"),
		StripeAPIKey:      getEnv("STRIPE_API_KEY"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return value
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 