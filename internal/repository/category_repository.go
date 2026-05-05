package repository

import "go-qr-order/internal/models"

// CategoryRepository persists menu categories and their products.
type CategoryRepository interface {
	SaveCategory(category models.Category) models.Category
	GetCategoryByID(id int) (models.Category, bool)
	GetAllCategories() []models.Category
}
