package services

import (
	"go-qr-order/internal/models"
	"go-qr-order/internal/repository"
)

// CategoryService coordinates category (menu group) operations.
type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) ListCategories() []models.Category {
	return s.repo.GetAllCategories()
}

func (s *CategoryService) GetCategory(id int) (models.Category, bool) {
	return s.repo.GetCategoryByID(id)
}

func (s *CategoryService) CreateCategory(category models.Category) models.Category {
	return s.repo.SaveCategory(category)
}
