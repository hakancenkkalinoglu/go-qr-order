package repository

import "go-qr-order/internal/models"

// ProductRepository persists products (always tied to a category).
type ProductRepository interface {
	Save(product models.Product) models.Product
	GetProductByID(id int) (models.Product, bool)
	ListProductsByCategoryID(categoryID int) []models.Product
}
