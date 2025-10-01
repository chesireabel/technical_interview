package models
import "time"

type Order struct {
	ID int64 `json:"id" db:"id"`
	CustomerID  string `json:"customer_id" db:"customer_id"`
	Item       string    `json:"item" db:"item"`
	Amount     float64   `json:"amount" db:"amount"`
	OrderedAt  time.Time `json:"ordered_at" db:"ordered_at"`
    CreatedAt  time.Time `db:"created_at" json:"created_at"`
}