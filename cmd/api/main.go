package main

import (
	"fmt"
	"go-qr-order/internal/handlers"
	"go-qr-order/internal/repository"
	"go-qr-order/internal/services"
	"net/http"
)

func main() {
	repo := repository.NewInMemoryOrderRepo()
	service := services.NewOrderService(repo)
	handler := handlers.NewOrderHandler(service)

	http.HandleFunc("POST /orders", handler.CreateOrderHandler)
	http.HandleFunc("GET /orders/{id}", handler.GetOrderHandler)
	http.HandleFunc("GET /orders", handler.GetAllOrdersHandler)
	http.HandleFunc("PUT /orders/{id}", handler.UpdateOrderById)
	http.HandleFunc("DELETE /orders/{id}", handler.DeleteOrderById)
	fmt.Println("Server starting..")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Print("Serves not started\n", err)
	}
}
