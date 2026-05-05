package main

import (
	"fmt"
	"go-qr-order/internal/database"
	"go-qr-order/internal/handlers"
	"go-qr-order/internal/middleware"
	"go-qr-order/internal/repository"
	"go-qr-order/internal/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbPath := strings.TrimSpace(os.Getenv("SQLITE_PATH"))
	if dbPath == "" {
		dbPath = "data/app.db"
	}

	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	orderRepo := repository.NewSQLiteOrderRepo(db)
	categoryRepo := repository.NewSQLiteCategoryRepo(db)
	productRepo := repository.NewSQLiteProductRepo(db)

	orderService := services.NewOrderService(orderRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo)

	orderHandler := handlers.NewOrderHandler(orderService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("GET /categories", middleware.RequestLogger(middleware.RequireAuth(categoryHandler.ListCategories)))
	http.HandleFunc("GET /categories/{id}", middleware.RequestLogger(middleware.RequireAuth(categoryHandler.GetCategory)))
	http.HandleFunc("POST /categories", middleware.RequestLogger(middleware.RequireAuth(categoryHandler.CreateCategory)))
	http.HandleFunc("GET /categories/{id}/products", middleware.RequestLogger(middleware.RequireAuth(productHandler.ListProductsByCategory)))
	http.HandleFunc("POST /products", middleware.RequestLogger(middleware.RequireAuth(productHandler.CreateProduct)))

	http.HandleFunc("POST /orders", middleware.RequestLogger(middleware.RequireAuth(orderHandler.CreateOrderHandler)))
	http.HandleFunc("GET /orders/{id}", middleware.RequestLogger(middleware.RequireAuth(orderHandler.GetOrderHandler)))
	http.HandleFunc("GET /orders", middleware.RequestLogger(middleware.RequireAuth(orderHandler.GetAllOrdersHandler)))
	http.HandleFunc("PUT /orders/{id}", middleware.RequestLogger(middleware.RequireAuth(orderHandler.UpdateOrderById)))
	http.HandleFunc("DELETE /orders/{id}", middleware.RequestLogger(middleware.RequireAuth(orderHandler.DeleteOrderById)))

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	fmt.Println("Server starting on", addr)

	err = http.ListenAndServe(addr, nil)

	if err != nil {
		fmt.Print("Serves not started\n", err)
	}
}
