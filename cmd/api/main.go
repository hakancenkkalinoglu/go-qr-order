package main

import (
	"fmt"
	"go-qr-order/internal/handlers"
	"go-qr-order/internal/middleware"
	"go-qr-order/internal/repository"
	"go-qr-order/internal/services"
	"net/http"
)

func main() {
	repo := repository.NewInMemoryOrderRepo()
	service := services.NewOrderService(repo)
	handler := handlers.NewOrderHandler(service)

	http.HandleFunc("POST /orders", middleware.RequestLogger(middleware.RequireAuth(handler.CreateOrderHandler)))
	http.HandleFunc("GET /orders/{id}", middleware.RequestLogger(middleware.RequireAuth(handler.GetOrderHandler)))
	http.HandleFunc("GET /orders", middleware.RequestLogger(middleware.RequireAuth(handler.GetAllOrdersHandler)))
	http.HandleFunc("PUT /orders/{id}", middleware.RequestLogger(middleware.RequireAuth(handler.UpdateOrderById)))
	http.HandleFunc("DELETE /orders/{id}", middleware.RequestLogger(middleware.RequireAuth(handler.DeleteOrderById)))

	fmt.Println("Server starting..")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Print("Serves not started\n", err)
	}
}
