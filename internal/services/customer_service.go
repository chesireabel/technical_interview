package services

import (
	"context"
	"errors"
	"time"

	"github.com/chesireabel/Technical-Interview/internal/models"
	"github.com/chesireabel/Technical-Interview/internal/repositories"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, customer *models.Customer) (int64, error)
	GetCustomer(ctx context.Context, id int64) (*models.Customer, error)
	GetAllCustomers(ctx context.Context) ([]models.Customer, error)
	UpdateCustomer(ctx context.Context, customer *models.Customer) error
	DeleteCustomer(ctx context.Context, id int64) error
}

type customerService struct {
	repo repositories.CustomerRepository
}

func NewCustomerService(repo repositories.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) CreateCustomer(ctx context.Context, customer *models.Customer) (int64, error) {
	if customer.Customer_name == "" || customer.Email == ""  || customer.Password == "" || customer.Phone == ""{
		return 0, errors.New("all fields  are required")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.Create(ctx, customer)
}

func (s *customerService) GetCustomer(ctx context.Context, id int64) (*models.Customer, error) {
	if id == 0 {
		return nil, errors.New("id is required")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.GetByID(ctx, id)
}

func (s *customerService) GetAllCustomers(ctx context.Context) ([]models.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.repo.GetAll(ctx)
}

func (s *customerService) UpdateCustomer(ctx context.Context, customer *models.Customer) error {
	if customer.ID == 0 {
		return errors.New("id is required for update")
	}

	if customer.Customer_name == "" || customer.Email == "" {
		return errors.New("name and email are required")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.Update(ctx, customer)
}

func (s *customerService) DeleteCustomer(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("id is required for delete")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.Delete(ctx, id)
}