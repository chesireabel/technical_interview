package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/chesireabel/Technical-Interview/internal/models"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *models.Customer) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Customer, error)
	GetAll(ctx context.Context) ([]models.Customer, error)
	Update(ctx context.Context, customer *models.Customer) error
	Delete(ctx context.Context, id int64) error
}

type customerRepository struct {
	db *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *models.Customer) (int64, error) {
	query := `
		INSERT INTO customers (customer_name, email, password, phone, code, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id
	`

	var id int64
	err := r.db.QueryRow(ctx, query,
		customer.Customer_name,
		customer.Email,
		customer.Password,
		customer.Phone,
		customer.Code,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create customer: %w", err)
	}
	return id, nil
}

func (r *customerRepository) GetByID(ctx context.Context, id int64) (*models.Customer, error) {
	var c models.Customer
	query := `
		SELECT id, customer_name, email, password, phone, code, created_at 
		FROM customers 
		WHERE id = $1
	`
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.Customer_name,
		&c.Email,
		&c.Password,
		&c.Phone,
		&c.Code,
		&c.CreatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}
	return &c, nil
}

func (r *customerRepository) GetAll(ctx context.Context) ([]models.Customer, error) {
	var customers []models.Customer
	query := `
		SELECT id, customer_name, email, password, phone, code, created_at 
		FROM customers 
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Customer
		err := rows.Scan(
			&c.ID,
			&c.Customer_name,
			&c.Email,
			&c.Password,
			&c.Phone,
			&c.Code,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating customers: %w", err)
	}

	return customers, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *models.Customer) error {
	query := `
		UPDATE customers
		SET customer_name = $1, email = $2, password = $3, phone = $4, code = $5
		WHERE id = $6
	`
	
	cmdTag, err := r.db.Exec(ctx, query,
		customer.Customer_name,
		customer.Email,
		customer.Password,
		customer.Phone,
		customer.Code,
		customer.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("customer with id %d not found", customer.ID)
	}

	return nil
}

func (r *customerRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM customers WHERE id = $1"
	
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("customer with id %d not found", id)
	}

	return nil
}