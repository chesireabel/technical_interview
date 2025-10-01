package services

import (
	"context"
	"testing"

	"github.com/chesireabel/Technical-Interview/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrder(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockCustomerRepo := new(MockCustomerRepo)
	mockSMS := new(MockSMSService)
	

	service := NewOrderService(mockOrderRepo, mockCustomerRepo, mockSMS)

	customer := &models.Customer{
		ID:            1,
		Customer_name: "Jane",
		Email:         "jane@example.com",
		Phone:         "+254712345678",
	}
	order := &models.Order{
		CustomerID: "1",
		Item:       "Laptop",
		Amount:     1000,
	}

	// Setup mock responses
	mockCustomerRepo.On("GetByID", mock.Anything, int64(1)).Return(customer, nil)
    mockOrderRepo.On("Create", mock.Anything, order).Return(int64(1), nil)
mockSMS.On(
    "SendOrderConfirmation",
    mock.Anything, // context.Context
    mock.AnythingOfType("*models.Order"),
    mock.AnythingOfType("*models.Customer"),
).Return(nil)
	orderID, err := service.CreateOrder(context.Background(), order)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), orderID)

	mockCustomerRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
	mockSMS.AssertExpectations(t)
}

func TestCreateOrder_InvalidCustomerID(t *testing.T) {
	service := NewOrderService(nil, nil, nil)

	order := &models.Order{
		CustomerID: "",
		Item:       "Laptop",
		Amount:     1000,
	}

	_, err := service.CreateOrder(context.Background(), order)

	assert.Error(t, err)
	assert.Equal(t, "customer_id is required", err.Error())
}

func TestGetOrder(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	service := NewOrderService(mockOrderRepo, nil, nil)

	expected := &models.Order{ID: 1, Item: "Laptop", Amount: 1000}
	mockOrderRepo.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

	order, err := service.GetOrder(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, order)

	_, err = service.GetOrder(context.Background(), 0)
	assert.Error(t, err)
	assert.Equal(t, "id is required", err.Error())
}
