package services

import (
	"context"
	"testing"

	"github.com/chesireabel/Technical-Interview/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCustomer(t *testing.T) {
	mockRepo := new(MockCustomerRepo)
	service := NewCustomerService(mockRepo)

	// Success case
	customer := &models.Customer{
		Customer_name: "Omondi",
		Email:         "omonditimon@example.com",
		Password:      "12345",
		Phone:         "+254712345678",
	}
mockRepo.On("Create", mock.Anything, customer).Return(int64(1), nil)

	id, err := service.CreateCustomer(context.Background(), customer)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockRepo.AssertExpectations(t)

	// Failure: missing fields
	badCustomer := &models.Customer{}
	_, err = service.CreateCustomer(context.Background(), badCustomer)
	assert.Error(t, err)
	assert.Equal(t, "all fields  are required", err.Error())
}

func TestGetCustomer(t *testing.T) {
	mockRepo := new(MockCustomerRepo)
	service := NewCustomerService(mockRepo)

	// Success case
	expected := &models.Customer{
		ID:            1,
		Customer_name: "Jane",
		Email:         "jane@example.com",
	}
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

	customer, err := service.GetCustomer(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, customer)

	// Failure: invalid ID
	_, err = service.GetCustomer(context.Background(), 0)
	assert.Error(t, err)
	assert.Equal(t, "id is required", err.Error())
}
