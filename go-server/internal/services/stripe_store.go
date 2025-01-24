package services

import (
	"fmt"
	"sync"
)

// StripeKeyStore manages Stripe API keys for users
type StripeKeyStore struct {
	keys map[string]string // userID -> stripeAPIKey
	mu   sync.RWMutex
}

// NewStripeKeyStore creates a new StripeKeyStore
func NewStripeKeyStore() *StripeKeyStore {
	return &StripeKeyStore{
		keys: make(map[string]string),
	}
}

// RegisterKey registers a Stripe API key for a user
func (s *StripeKeyStore) RegisterKey(userID, apiKey string) error {
	if userID == "" || apiKey == "" {
		return fmt.Errorf("userID and apiKey cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.keys[userID] = apiKey
	return nil
}

// GetStripeKey retrieves a user's Stripe API key
func (s *StripeKeyStore) GetStripeKey(userID string) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("userID cannot be empty")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	key, exists := s.keys[userID]
	if !exists {
		return "", fmt.Errorf("no Stripe API key found for user %s", userID)
	}
	return key, nil
}

// RemoveKey removes a user's Stripe API key
func (s *StripeKeyStore) RemoveKey(userID string) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.keys, userID)
	return nil
}
