package repository

import (
	"go-qr-order/internal/models"
)

type InMemoryOrderRepo struct {
	orders map[int]models.Order
}

func NewInMemoryOrderRepo() *InMemoryOrderRepo {
	return &InMemoryOrderRepo{
		orders: make(map[int]models.Order),
	}
}

func (r *InMemoryOrderRepo) Save(order models.Order) {
	r.orders[order.ID] = order
}

func (r *InMemoryOrderRepo) GetById(id int) (models.Order, bool) {
	order, exist := r.orders[id]
	return order, exist
}
