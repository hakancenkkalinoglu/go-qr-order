package handlers

import (
	"encoding/json"
	"go-qr-order/internal/models"
	"go-qr-order/internal/services"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(s *services.OrderService) *OrderHandler {
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

func (h *OrderHandler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid Path Parameter", http.StatusBadRequest)
		return
	}
	order, isExist := h.service.GetOrder(id)
	if !isExist {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
