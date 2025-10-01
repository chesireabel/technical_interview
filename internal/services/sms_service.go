// File: internal/services/sms_service.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chesireabel/Technical-Interview/internal/models"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}


type SMSService interface {
	SendOrderConfirmation(ctx context.Context, order *models.Order, customer *models.Customer) error
	SendOrderUpdate(ctx context.Context, order *models.Order, phoneNumber, status string) error
}

type smsService struct {
	username   string
	apiKey     string
	baseURL    string
	httpClient HTTPClient
}

type SMSResponse struct {
	SMSMessageData struct {
		Message    string `json:"Message"`
		Recipients []struct {
			StatusCode int    `json:"statusCode"`
			Number     string `json:"number"`
			Status     string `json:"status"`
			Cost       string `json:"cost"`
			MessageID  string `json:"messageId"`
		} `json:"Recipients"`
	} `json:"SMSMessageData"`
}

func NewSMSService() (SMSService, error) {
	username := os.Getenv("AT_USERNAME")
	apiKey := os.Getenv("AT_API_KEY")
	baseURL := os.Getenv("AT_BASE_URL")

	if username == "" || apiKey == "" {
		return nil, fmt.Errorf("AT_USERNAME and AT_API_KEY must be set in environment")
	}

	if baseURL == "" {
		baseURL = "https://api.africastalking.com/version1"
	}

	return &smsService{
		username: username,
		apiKey:   apiKey,
		baseURL:  baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (s *smsService) sendSMS(phoneNumber, message string) error {
	// Validate phone number format
	if !strings.HasPrefix(phoneNumber, "+") {
		return fmt.Errorf("phone number must start with country code (e.g., +254)")
	}

	// Prepare form data
	data := url.Values{}
	data.Set("username", s.username)
	data.Set("to", phoneNumber)
	data.Set("message", message)

	// Create request
	endpoint := fmt.Sprintf("%s/messaging", s.baseURL)
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", s.apiKey)

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("SMS API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response to check if message was sent
	var smsResp SMSResponse
	if err := json.Unmarshal(body, &smsResp); err != nil {
		return fmt.Errorf("failed to parse SMS response: %w", err)
	}

	// Check if any recipient failed
	if len(smsResp.SMSMessageData.Recipients) > 0 {
		recipient := smsResp.SMSMessageData.Recipients[0]
		if recipient.StatusCode != 100 && recipient.StatusCode != 101 {
			return fmt.Errorf("SMS delivery failed: %s (code: %d)", recipient.Status, recipient.StatusCode)
		}
	}

	return nil
}

func (s *smsService) SendOrderConfirmation(ctx context.Context, order *models.Order, customer *models.Customer) error {
	message := fmt.Sprintf(
		"Hello %s! Order #%d confirmed. Item: %s, Amount: KES %.2f. Thank you for your order!",
		customer.Customer_name,
		order.ID,
		order.Item,
		order.Amount,
	)

	return s.sendSMS(customer.Phone, message)
}

func (s *smsService) SendOrderUpdate(ctx context.Context, order *models.Order, phoneNumber, status string) error {
	message := fmt.Sprintf(
		"Order Update: Your order #%d (%s) is now %s. Amount: KES %.2f",
		order.ID,
		order.Item,
		status,
		order.Amount,
	)

	return s.sendSMS(phoneNumber, message)
}