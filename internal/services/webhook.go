package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookService struct {
	secret string
}

type WebhookPayload struct {
	Event     string                 `json:"event"`
	PaymentID string                 `json:"payment_id"`
	Status    string                 `json:"status"`
	Amount    string                 `json:"amount"`
	Currency  string                 `json:"currency"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp int64                  `json:"timestamp"`
}

func NewWebhookService(secret string) *WebhookService {
	return &WebhookService{
		secret: secret,
	}
}

func (s *WebhookService) SendWebhook(url string, payload WebhookPayload) error {
	if url == "" {
		return nil // No webhook URL provided
	}

	// Add timestamp
	payload.Timestamp = time.Now().Unix()

	// Serialize payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create signature
	signature := s.createSignature(payloadBytes)

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", signature)
	req.Header.Set("User-Agent", "CryptoPaymentGateway/1.0")

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

func (s *WebhookService) createSignature(payload []byte) string {
	h := hmac.New(sha256.New, []byte(s.secret))
	h.Write(payload)
	return "sha256=" + hex.EncodeToString(h.Sum(nil))
}

func (s *WebhookService) VerifySignature(payload []byte, signature string) bool {
	expectedSignature := s.createSignature(payload)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}