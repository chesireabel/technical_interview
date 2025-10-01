package services

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/chesireabel/Technical-Interview/internal/models"
	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestSendOrderConfirmation(t *testing.T) {
	mockResp := `{
		"SMSMessageData": {
			"Message": "Sent",
			"Recipients": [{
				"statusCode": 100,
				"number": "+254700000000",
				"status": "Success",
				"cost": "KES 0.00",
				"messageId": "12345"
			}]
		}
	}`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResp)),
			}, nil
		},
	}

	svc := &smsService{
		username:   "testuser",
		apiKey:     "testkey",
		baseURL:    "https://mockapi.test",
		httpClient: mockClient,
	}

	order := &models.Order{ID: 1, Item: "Book", Amount: 500}
	customer := &models.Customer{Customer_name: "Jane", Phone: "+254700000000"}

	err := svc.SendOrderConfirmation(context.Background(), order, customer)
	assert.NoError(t, err)
}

func TestSendOrderUpdate(t *testing.T) {
	mockResp := `{
		"SMSMessageData": {
			"Message": "Sent",
			"Recipients": [{
				"statusCode": 100,
				"number": "+254700000000",
				"status": "Success",
				"cost": "KES 0.00",
				"messageId": "12345"
			}]
		}
	}`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResp)),
			}, nil
		},
	}

	svc := &smsService{
		username:   "testuser",
		apiKey:     "testkey",
		baseURL:    "https://mockapi.test",
		httpClient: mockClient,
	}

	order := &models.Order{ID: 1, Item: "Book", Amount: 500}
	phone := "+254700000000"
	status := "Delivered"

	err := svc.SendOrderUpdate(context.Background(), order, phone, status)
	assert.NoError(t, err)
}
