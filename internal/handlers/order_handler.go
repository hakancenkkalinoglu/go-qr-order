package handlers

import (
	"encoding/json"
	"go-qr-order/internal/models"
	"go-qr-order/internal/services"
	"net/http"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderService(s *services.OrderService) *OrderHandler {
	return &OrderHandler{
		service: s,
	}
}

func (h *OrderHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var reqOrder models.Order

	err := json.NewDecoder(r.Body).Decode(&reqOrder)
	if err != nil {
		http.Error(w, "Invalid Json data", http.StatusBadRequest)
		return
	}

	createdOrder := h.service.CreateOrder(reqOrder)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(createdOrder)

}
