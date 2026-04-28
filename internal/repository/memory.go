package repository

import (
	"go-qr-order/internal/models"
	"time"
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

func (r *InMemoryOrderRepo) GetAll() []models.Order {
	var allOrdersSlice []models.Order
	for _, value := range r.orders {
		allOrdersSlice = append(allOrdersSlice, value)
	}
	return allOrdersSlice
}

func (r *InMemoryOrderRepo) UpdateOrderById(id int, status string) models.Order {
	order := r.orders[id]
	order.Status = status
	order.UpdatedAt = time.Now()

	//avoid pass by value and ovveride the slice

	r.orders[id] = order

	return order
}

func (r *InMemoryOrderRepo) DeleteOrderById(id int) bool {
	delete(r.orders, id)
	return true
}
