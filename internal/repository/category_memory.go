package repository

import "go-qr-order/internal/models"

type InMemoryCategoryRepo struct {
	categories map[int]models.Category
	lastID     int
}

func NewInCategoryRepo() *InMemoryCategoryRepo {
	return &InMemoryCategoryRepo{
		categories: make(map[int]models.Category),
	}
}

func (c *InMemoryCategoryRepo) SaveCategory(category models.Category) models.Category {
	c.lastID++

	category.ID = c.lastID
	c.categories[category.ID] = category

	return category
}
