package models

import "time"

type Order struct {
	ID         int
	TableID    int
	SessionID  string
	TotalPrice float64
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Items      []OrderItem
}

type OrderItem struct {
	ID       int
	Name     string
	Quantity int
	Price    float64
}

type Category struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	CategoryItems []Product `json:"category_items,omitempty"`
}

type Product struct {
	ID         int     `json:"id"`
	CategoryID int     `json:"category_id"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}
