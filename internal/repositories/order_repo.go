package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/chesireabel/Technical-Interview/internal/models"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Order, error)
	GetByCustomerID(ctx context.Context, customerID int64) ([]models.Order, error)
	GetAll(ctx context.Context) ([]models.Order, error)
	Update(ctx context.Context, order *models.Order) error
	Delete(ctx context.Context, id int64) error
}

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) (int64, error) {
	query := `
		INSERT INTO orders (customer_id, item, amount, ordered_at, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id
	`

	var id int64
	err := r.db.QueryRow(
		ctx,
		query,
		order.CustomerID,
		order.Item,
		order.Amount,
		order.OrderedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}
	return id, nil
}

func (r *orderRepository) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	var o models.Order
	query := `
		SELECT id, customer_id, item, amount, ordered_at, created_at 
		FROM orders 
		WHERE id = $1
	`
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&o.ID,
		&o.CustomerID,
		&o.Item,
		&o.Amount,
		&o.OrderedAt,
		&o.CreatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	return &o, nil
}

func (r *orderRepository) GetByCustomerID(ctx context.Context, customerID int64) ([]models.Order, error) {
	var orders []models.Order
	query := `
		SELECT id, customer_id, item, amount, ordered_at, created_at 
		FROM orders 
		WHERE customer_id = $1
		ORDER BY ordered_at DESC
	`
	
	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders for customer: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Order
		err := rows.Scan(
			&o.ID,
			&o.CustomerID,
			&o.Item,
			&o.Amount,
			&o.OrderedAt,
			&o.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}

func (r *orderRepository) GetAll(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	query := `
		SELECT id, customer_id, item, amount, ordered_at, created_at 
		FROM orders 
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Order
		err := rows.Scan(
			&o.ID,
			&o.CustomerID,
			&o.Item,
			&o.Amount,
			&o.OrderedAt,
			&o.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}

func (r *orderRepository) Update(ctx context.Context, order *models.Order) error {
	query := `
		UPDATE orders
		SET customer_id = $1, item = $2, amount = $3, ordered_at = $4
		WHERE id = $5
	`
	
	cmdTag, err := r.db.Exec(
		ctx,
		query,
		order.CustomerID,
		order.Item,
		order.Amount,
		order.OrderedAt,
		order.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("order with id %d not found", order.ID)
	}

	return nil
}

func (r *orderRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM orders WHERE id = $1"
	
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("order with id %d not found", id)
	}

	return nil
}