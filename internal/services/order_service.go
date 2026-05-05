package services

import (
	"go-qr-order/internal/models"
	"go-qr-order/internal/repository"
	"time"
)

type OrderService struct {
	repo repository.OrderRepository
}

// Constructor
func NewOrderService(r repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: r,
	}
}

func (s *OrderService) CreateOrder(order models.Order) models.Order {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	if order.Status == "" {
		order.Status = "PENDING"
	}

	order = s.repo.Save(order)

	return order
}

func (s *OrderService) GetOrder(id int) (models.Order, bool) {
	return s.repo.GetById(id)
}

func (s *OrderService) GetAllOrders() []models.Order {
	return s.repo.GetAll()
}

func (s *OrderService) UpdateOrderById(id int, status string) (models.Order, bool) {
	_, isExist := s.GetOrder(id)
	if !isExist {
		return models.Order{}, false
	}

	updatedOrder := s.repo.UpdateOrderById(id, status)
	return updatedOrder, true
}

func (s *OrderService) DeleteOrderById(id int) bool {
	_, isExist := s.GetOrder(id)
	if !isExist {
		return false
	}
	s.repo.DeleteOrderById(id)
	return true
}
