package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TurnstileService struct {
	secretKey string
	enabled   bool
}

type TurnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
}

func NewTurnstileService(secretKey string, enabled bool) *TurnstileService {
	return &TurnstileService{
		secretKey: secretKey,
		enabled:   enabled,
	}
}

func (s *TurnstileService) VerifyToken(token string) (bool, error) {
	if !s.enabled {
		return true, nil
	}

	url := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	payload := map[string]string{
		"secret":   s.secretKey,
		"response": token,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		// TODO: log error
		return false, fmt.Errorf("failed to marshal turnstile request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// TODO: log error
		return false, fmt.Errorf("failed to send turnstile request: %w", err)
	}
	defer resp.Body.Close()

	var result TurnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode turnstile response: %w", err)
	}

	return result.Success, nil
} 
