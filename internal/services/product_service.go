package services

import (
	"go-qr-order/internal/models"
	"go-qr-order/internal/repository"
)

// ProductService coordinates catalog product operations.
type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product models.Product) (models.Product, bool) {
	if product.CategoryID <= 0 {
		return models.Product{}, false
	}
	return s.repo.Save(product), true
}

func (s *ProductService) ListByCategoryID(categoryID int) []models.Product {
	return s.repo.ListProductsByCategoryID(categoryID)
}
