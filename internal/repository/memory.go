package repository

import (
	"go-qr-order/internal/models"
)

type InMemoryOrderRepo struct {
	orders map[int]models.Order
	lastID int
}

func NewInMemoryOrderRepo() *InMemoryOrderRepo {
	return &InMemoryOrderRepo{
		orders: make(map[int]models.Order),
	}
}

func (r *InMemoryOrderRepo) Save(order models.Order) models.Order {
	r.lastID++

	order.ID = r.lastID
	r.orders[order.ID] = order

	return order
}

func (r *InMemoryOrderRepo) GetById(id int) (models.Order, bool) {
	order, exist := r.orders[id]
	return order, exist
}
