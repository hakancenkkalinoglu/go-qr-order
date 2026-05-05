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
	ID            int
	name          string
	CategoryItems []Product
}

type Product struct {
	ID       int
	Name     string
	Quantity int
	Price    float64
}
