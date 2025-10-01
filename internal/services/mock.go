package services

import (
	"context"

	"github.com/chesireabel/Technical-Interview/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockCustomerRepo struct {
	mock.Mock
}

func (m *MockCustomerRepo) Create(ctx context.Context, customer *models.Customer) (int64, error) {
	args := m.Called(ctx, customer)
return args.Get(0).(int64), args.Error(1)
}

func (m *MockCustomerRepo) GetByID(ctx context.Context, id int64) (*models.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCustomerRepo) GetAll(ctx context.Context) ([]models.Customer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Customer), args.Error(1)
}

func (m *MockCustomerRepo) Update(ctx context.Context, customer *models.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomerRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(ctx context.Context, order *models.Order) (int64, error) {
	args := m.Called(ctx, order)
return args.Get(0).(int64), args.Error(1)
}

func (m *MockOrderRepo) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockOrderRepo) GetByCustomerID(ctx context.Context, customerID int64) ([]models.Order, error) {
	args := m.Called(ctx, customerID)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockOrderRepo) GetAll(ctx context.Context) ([]models.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockOrderRepo) Update(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockSMSService struct {
	mock.Mock
}

func (m *MockSMSService) SendOrderConfirmation(ctx context.Context, order *models.Order, customer *models.Customer) error {
	args := m.Called(ctx, order, customer)
	return args.Error(0)
}

func (m *MockSMSService) SendOrderUpdate(ctx context.Context, order *models.Order, phoneNumber, status string) error {
	args := m.Called(ctx, order, phoneNumber, status)
	return args.Error(0)
}
