package models

import "time"

type Customer struct {
	ID  int64  `json:"id" db:"id"`
	Customer_name string `json:"customer_name" db:"customer_name"`
	Email string `json:"email" db:"email"` 
	Password string `json:"password,omitempty"`
	Phone string `json:"phone" db:"phone"`
	Code string `json:"code" db:"code"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}