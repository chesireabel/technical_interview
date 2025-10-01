// File: internal/services/order_service.go
package services

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/chesireabel/Technical-Interview/internal/models"
	"github.com/chesireabel/Technical-Interview/internal/repositories"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *models.Order) (int64, error)
	GetOrder(ctx context.Context, id int64) (*models.Order, error)
	GetOrdersByCustomer(ctx context.Context, customerID int64) ([]models.Order, error)
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
	DeleteOrder(ctx context.Context, id int64) error
}

type orderService struct {
	repo            repositories.OrderRepository
	customerRepo    repositories.CustomerRepository
	smsService      SMSService
}

func NewOrderService(repo repositories.OrderRepository, customerRepo repositories.CustomerRepository, smsService SMSService) OrderService {
	return &orderService{
		repo:         repo,
		customerRepo: customerRepo,
		smsService:   smsService,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, order *models.Order) (int64, error) {
	// Validation
	if order.CustomerID == "" {
		return 0, errors.New("customer_id is required")
	}
	if order.Item == "" {
		return 0, errors.New("item is required")
	}
	if order.Amount <= 0 {
		return 0, errors.New("amount must be greater than 0")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get customer details to fetch phone number
	customerIDInt, err := strconv.ParseInt(order.CustomerID, 10, 64)
	if err != nil {
		return 0, errors.New("invalid customer_id format")
	}

	customer, err := s.customerRepo.GetByID(ctx, customerIDInt)
	if err != nil {
		return 0, errors.New("customer not found")
	}

	// Create order in database
	orderID, err := s.repo.Create(ctx, order)
	if err != nil {
		return 0, err
	}

	// Set the order ID for SMS notification
	order.ID = orderID

	// Send SMS notification asynchronously to avoid blocking
	if s.smsService != nil && customer.Phone != "" {
			// Create a new context for the async operation with timeout
			smsCtx, smsCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer smsCancel()

			if err := s.smsService.SendOrderConfirmation(smsCtx, order, customer); err != nil {
				// Log error but don't fail the order creation
				log.Printf("Failed to send SMS for order %d to %s: %v", orderID, customer.Phone, err)
			} else {
				log.Printf("SMS sent successfully for order %d to %s (%s)", orderID, customer.Customer_name, customer.Phone)
			}
		
	} else {
		log.Printf("SMS service not available or customer phone missing for order %d", orderID)
	}

	return orderID, nil
}

func (s *orderService) GetOrder(ctx context.Context, id int64) (*models.Order, error) {
	if id == 0 {
		return nil, errors.New("id is required")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.GetByID(ctx, id)
}

func (s *orderService) GetOrdersByCustomer(ctx context.Context, customerID int64) ([]models.Order, error) {
	if customerID == 0 {
		return nil, errors.New("customer_id is required")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.repo.GetByCustomerID(ctx, customerID)
}

func (s *orderService) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.repo.GetAll(ctx)
}

func (s *orderService) UpdateOrder(ctx context.Context, order *models.Order) error {
	if order.ID == 0 {
		return errors.New("id is required for update")
	}
	if order.CustomerID == "" {
		return errors.New("customer_id is required")
	}
	if order.Item == "" {
		return errors.New("item is required")
	}
	if order.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.Update(ctx, order)
}

func (s *orderService) DeleteOrder(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("id is required for delete")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.Delete(ctx, id)
}