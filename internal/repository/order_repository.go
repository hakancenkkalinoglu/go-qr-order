package repository

import "go-qr-order/internal/models"

// OrderRepository persists orders and their line items.
type OrderRepository interface {
	Save(order models.Order) models.Order
	GetById(id int) (models.Order, bool)
	GetAll() []models.Order
	UpdateOrderById(id int, status string) models.Order
	DeleteOrderById(id int) bool
}
